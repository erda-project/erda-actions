package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/erda-project/erda/apistructs"

	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/strutil"
)

const sourceLineRegionLine int = 5

// see: https://docs.sonarqube.org/latest/user-guide/issues/

type IssueType string

var (
	IssueTypeBug           IssueType = "BUG"
	IssueTypeVulnerability IssueType = "VULNERABILITY"
	IssueTypeCodeSmell     IssueType = "CODE_SMELL"
)

func (i IssueType) String() string {
	return string(i)
}

type IssueSeverity string

var (
	IssueSeverityBlocker  IssueSeverity = "BLOCKER"
	IssueSeverityCritical IssueSeverity = "CRITICAL"
	IssueSeverityMajor    IssueSeverity = "MAJOR"
	IssueSeverityMinor    IssueSeverity = "MINOR"
	IssueSeverityInfo     IssueSeverity = "INFO"
)

func (i IssueSeverity) String() string {
	return string(i)
}

type IssueSearchResponse struct {
	Issues []Issue `json:"issues"`
}

type Issue struct {
	Key           string               `json:"key"` // issueKey
	Author        string               `json:"author"`
	Project       string               `json:"project"`
	LineNum       int                  `json:"line"`
	Message       string               `json:"message"`
	Rule          string               `json:"rule"`
	Severity      IssueSeverity        `json:"severity"`
	TextRange     apistructs.TextRange `json:"textRange"`
	Type          IssueType            `json:"type"`
	CreationDate  IssueTime            `json:"creationDate"`
	Component     string               `json:"component"`
	ComponentPath string               `json:"componentPath"`
	Flows         IssueFlows           `json:"flows"`

	// Snippet 返回的 codes 带有前端 CSS 样式，不能直接使用
	// 但是要使用 snippet 里的
	Snippet IssueSnippet `json:"snippet"`

	// 使用 snippet 里返回的行号，获取代码
	SourceLines []SourceLine `json:"sourceLines"`
}

// SourceLine raw 只返回 line 和 code 两个字段
type SourceLine struct {
	Line int    `json:"line"`
	Code string `json:"code"`

	Duplicated bool `json:"duplicated"` // 当前行是否重复
	LineHits   *int `json:"lineHits"`   // 当前行是否命中（是否重复或者是否未被测试覆盖），不为空则为命中
}

type IssueFlowLocation struct {
	Locations []struct {
		TextRange apistructs.TextRange `json:"textRange"`
	} `json:"locations"`
}

type IssueFlows []IssueFlowLocation

type IssueTime struct {
	time.Time
}

const IssueTimeLayout = "2006-01-02T15:04:05-0700"

func (t *IssueTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)
	if strings.HasPrefix(s, `"`) {
		s = s[1:]
	}
	if strings.HasSuffix(s, `"`) {
		s = s[:len(s)-1]
	}
	t.Time, err = time.Parse(IssueTimeLayout, s)
	return err
}

// querySonarIssues 查询 sonar issues
// web api: /api/issues/search
func (sonar *Sonar) querySonarIssues(projectKey string, issueTypes ...IssueType) ([]Issue, error) {
	if len(issueTypes) == 0 {
		return nil, nil
	}
	var issueTypeStrs []string
	for _, t := range issueTypes {
		issueTypeStrs = append(issueTypeStrs, t.String())
	}
	var body bytes.Buffer
	resp, err := httpclient.New(httpclient.WithCompleteRedirect()).
		BasicAuth(sonar.Auth.Login, sonar.Auth.Password).Get(sonar.Auth.HostURL).
		Path("/api/issues/search").
		Param("componentKeys", projectKey).
		Param("types", strutil.Join(issueTypeStrs, ",", true)).
		Param("s", "FILE_LINE").
		Param("ps", "500"). // max value: 500
		Do().
		Body(&body)
	if err != nil {
		return nil, err
	}
	if !resp.IsOK() {
		return nil, fmt.Errorf("statusCode: %d, responseBody: %s", resp.StatusCode(), body.String())
	}
	var issueSearchResp IssueSearchResponse
	if err := json.Unmarshal(body.Bytes(), &issueSearchResp); err != nil {
		return nil, err
	}
	// polish issues
	for i := range issueSearchResp.Issues {
		issue := issueSearchResp.Issues[i]
		// file path
		if idx := strings.Index(issue.Component, ":"); idx > -1 {
			issueSearchResp.Issues[i].ComponentPath = issue.Component[idx+1:]
		}
		// calc source line numbers
		var (
			lineNumFrom int
			lineNumTo   int
		)
		// get proper line number range from snippet
		issueSnippet, _ := sonar.querySonarIssueSnippets(issue.Key)
		if issueSnippet != nil && len(issueSnippet.Sources) > 0 {
			lineNumFrom = issueSnippet.Sources[0].Line
			lineNumTo = issueSnippet.Sources[len(issueSnippet.Sources)-1].Line
		}
		// 若起始行为 0，则默认取 issueLine 前几行
		if lineNumFrom == 0 {
			lineNumFrom = issue.LineNum - sourceLineRegionLine
			// 不能小于 1
			if lineNumFrom <= 0 {
				lineNumFrom = 1
			}
		}
		// 若结束行为 0，则默认取 issueLine 后几行
		if lineNumTo == 0 {
			lineNumTo = issue.LineNum + sourceLineRegionLine
		}
		// 结束行需要比起始行大
		if lineNumFrom >= lineNumTo {
			lineNumTo = lineNumFrom + sourceLineRegionLine
		}
		// 从 sonar source lines 接口获取指定代码
		sourceLines, err := sonar.querySonarSourceLinesRaw(issue.Component, lineNumFrom, lineNumTo)
		if err != nil {
			return nil, fmt.Errorf("failed to get source lines, component: %s, err: %v", issue.Component, err)
		}
		issueSearchResp.Issues[i].SourceLines = sourceLines
	}

	return issueSearchResp.Issues, nil
}

type IssueSnippetResp map[string]IssueSnippet

type IssueSnippet struct {
	Sources []struct {
		Line int `json:"line"` // 只关心行号
	} `json:"sources"`
}

func (sonar *Sonar) querySonarIssueSnippets(issueKey string) (*IssueSnippet, error) {
	var body bytes.Buffer
	httpResp, err := httpclient.New(httpclient.WithCompleteRedirect()).
		BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
		Get(sonar.Auth.HostURL).
		Path("/api/sources/issue_snippets").
		Param("issueKey", issueKey).
		Do().Body(&body)
	if err != nil {
		return nil, err
	}
	if !httpResp.IsOK() {
		return nil, fmt.Errorf("statusCode: %d, responseBody: %s", httpResp.StatusCode(), body.String())
	}
	var snippetResp IssueSnippetResp
	respBodyStr := body.String()
	if err := json.Unmarshal([]byte(respBodyStr), &snippetResp); err != nil {
		return nil, fmt.Errorf("failed to parse responseBody as json, err: %v, responseBody: %s", err, respBodyStr)
	}
	for _, s := range snippetResp {
		return &s, nil
	}
	return nil, fmt.Errorf("not found issue snippet, issueKey: %s", issueKey)
}

// querySonarSourceLinesRaw 返回代码行
func (sonar *Sonar) querySonarSourceLinesRaw(componentKey string, from, to int) ([]SourceLine, error) {
	var body bytes.Buffer
	httpResp, err := httpclient.New(httpclient.WithCompleteRedirect()).
		BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
		Get(sonar.Auth.HostURL).
		Path("/api/sources/index").
		Param("resource", componentKey).   // my_project:/src/foo/Bar.php
		Param("from", strconv.Itoa(from)). // First line
		Param("to", strconv.Itoa(to+1)).   // Last line (excluded)
		Do().Body(&body)
	if err != nil {
		return nil, err
	}
	if !httpResp.IsOK() {
		return nil, fmt.Errorf("statusCode: %d, responseBody: %s", httpResp.StatusCode(), body.String())
	}
	sourceResp := make([]map[string]string, 0)
	respBodyStr := body.String()
	if err := json.Unmarshal([]byte(respBodyStr), &sourceResp); err != nil {
		return nil, fmt.Errorf("failed to parse responseBody as json, err: %v, responseBody: %s", err, respBodyStr)
	}
	if len(sourceResp) == 0 {
		return nil, fmt.Errorf("not found source lines, componentKey: %s", componentKey)
	}
	sourceLines := make([]SourceLine, 0, len(sourceResp[0]))
	for lineNumStr, content := range sourceResp[0] {
		line, _ := strconv.Atoi(lineNumStr)
		sourceLines = append(sourceLines, SourceLine{Line: line, Code: content})
	}
	sort.SliceStable(sourceLines, func(i, j int) bool { return sourceLines[i].Line < sourceLines[j].Line })
	return sourceLines, nil
}

// querySonarSourceLinesRendered 返回渲染后的代码行
func (sonar *Sonar) querySonarSourceLinesRendered(componentKey string, from, to int) ([]SourceLine, error) {
	queryParams := make(url.Values, 0)
	queryParams.Add("key", componentKey)        // my_project:/src/foo/Bar.php
	queryParams.Add("from", strconv.Itoa(from)) // First line
	if to > 0 {
		queryParams.Add("to", strconv.Itoa(to)) // Last line (inclusive)
	}
	var body bytes.Buffer
	httpResp, err := httpclient.New(httpclient.WithCompleteRedirect()).
		BasicAuth(sonar.Auth.Login, sonar.Auth.Password).
		Get(sonar.Auth.HostURL).
		Path("/api/sources/lines").
		Params(queryParams).
		Do().Body(&body)
	if err != nil {
		return nil, err
	}
	if !httpResp.IsOK() {
		return nil, fmt.Errorf("statusCode: %d, responseBody: %s", httpResp.StatusCode(), body.String())
	}
	var sourceResp struct {
		Sources []SourceLine `json:"sources"`
	}
	respBodyStr := body.String()
	if err := json.Unmarshal([]byte(respBodyStr), &sourceResp); err != nil {
		return nil, fmt.Errorf("failed to parse responseBody as json, err: %v, responseBody: %s", err, respBodyStr)
	}
	sourceLines := make([]SourceLine, 0)
	for _, line := range sourceResp.Sources {
		sourceLines = append(sourceLines, SourceLine{
			Line:       line.Line,
			Code:       line.Code,
			Duplicated: line.Duplicated,
			LineHits:   line.LineHits,
		})
	}
	return sourceLines, nil
}
