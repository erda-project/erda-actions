package register

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/erda-project/erda-actions/actions/mcp-register/1.0/internal/common"
	"github.com/erda-project/erda-actions/actions/mcp-register/1.0/internal/conf"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/parser/diceyml"
	"github.com/erda-project/erda/pkg/retry"
)

type Register struct {
	cfg *conf.Conf
}

func New(cfg *conf.Conf) *Register {
	return &Register{
		cfg,
	}
}

func (r *Register) Register() error {
	var resp common.GetReleaseResponse

	if err := retry.DoWithInterval(func() error {
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).
			Get(r.cfg.DiceOpenapiPrefix).
			Path(common.ReleaseRequestPath+"/"+r.cfg.ReleaseID).
			Header(common.Authorization, r.cfg.DiceOpenapiToken).Do().JSON(&resp)
		if err != nil {
			return fmt.Errorf("failed to create http client, err: %v", err)
		}

		if !resp.Success {
			respErrs := []string{
				"", // empty line
				fmt.Sprintf("status code: %d", r.StatusCode()),
				fmt.Sprintf("response code: %s", resp.Err.Code),
				fmt.Sprintf("message: %s", resp.Err.Message),
				fmt.Sprintf("context: %s", resp.Err.Ctx),
				fmt.Sprintf("body: %s", string(r.Body())),
			}
			respErr := errors.New(strings.Join(respErrs, "\n"))
			logrus.Errorf("failed to register the release with resp: %v", respErr)
			return nil
		}

		return nil
	}, 2, time.Second*5); err != nil {
		logrus.Errorf("failed to get the release with resp: %v", err)
		return err
	}

	logrus.Infof("the release with resp: %v", resp.Data.DiceYaml)

	var diceYml diceyml.Object
	if err := yaml.Unmarshal([]byte(resp.Data.DiceYaml), &diceYml); err != nil {
		logrus.Errorf("failed to unmarshal the diceYml: %v", err)
		return err
	}

	logrus.Infof("the diceYml: %v", diceYml)

	result := make(map[string]string)

	for name, info := range r.cfg.Services {
		logger := logrus.WithField("service", name)

		service, ok := diceYml.Services[name]
		if !ok {
			logger.Errorf("failed to get the service with name %s", name)
			result[name] = fmt.Sprintf("failed to get the service with name %s", name)
			continue
		}

		tools, err := loadToolCallList(logger, service, info.Host)
		if err != nil {
			result[name] = err.Error()
			continue
		}

		sc := loadServerConfig(logger, service, info.Host)

		req := buildRequest(name, service, info.Host)
		req.Tools = tools
		req.ServerConfig = sc

		err = r.handleRegister(logger, req)
		if err != nil {
			result[name] = err.Error()
			continue
		}
		result[name] = "success"
	}

	logrus.Infof("=============")
	logrus.Infof("mcp server register result: ")

	for name, res := range result {
		logrus.Infof("\t%s: %s", name, res)
	}

	return nil
}

func (r *Register) handleRegister(logger *logrus.Entry, req *common.MCPServerRegisterRequest) error {
	logger = logger.WithField("action", "handleRegister")
	baseUrl := strings.TrimRight(r.cfg.McpProxyUrl, "/")

	url := fmt.Sprintf("%v/api/ai-proxy/mcp/servers/%v/actions/register", baseUrl, req.Name)

	logrus.Infof("mcp proxy url: %s", url)

	reqBody, err := json.Marshal(req)
	if err != nil {
		logger.Errorf("failed to marshal request body, err: %v", err)
		return fmt.Errorf("marshal error: %w", err)
	}

	logger.Infof("mcp register request body: %s", string(reqBody))

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
	if err != nil {
		logger.Errorf("failed to create request, err: %v", err)
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", r.cfg.McpProxyAccessKey)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Errorf("failed to do request, err: %v", err)
		return fmt.Errorf("can't send http request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("failed to register the release with resp: %v", resp.StatusCode)
		return fmt.Errorf("register release with resp: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to register mcp-server: %w", err)
	}
	logger.Infof("http status code: %d", resp.StatusCode)
	logger.Infof(string(body))
	return nil
}

func buildRequest(serviceName string, service *diceyml.Service, host string) *common.MCPServerRegisterRequest {
	endpoint := fmt.Sprintf("http://%s:%s%s", host, service.Labels[common.LabelMcpServicePort], service.Annotations[common.AnnotationMcpConnectURI])
	name := service.Labels[common.LabelMcpName]
	if name == "" {
		name = serviceName
	}

	return &common.MCPServerRegisterRequest{
		Description: service.Annotations[common.AnnotationMcpDescription],
		Endpoint:    endpoint,

		Name:             name,
		Version:          service.Labels[common.LabelMcpVersion],
		IsPublished:      service.Labels[common.LabelMcpIsPublished] != "false",
		IsDefaultVersion: service.Labels[common.LabelMcpIsDefault] == "true",
		TransportType:    service.Labels[common.LabelMcpTransportType],
	}
}

func loadToolCallList(logger *logrus.Entry, service *diceyml.Service, host string) ([]mcp.Tool, error) {
	logger = logger.WithField("action", "loadToolCallList")
	endpoint := fmt.Sprintf("http://%s:%s%s", host, service.Labels[common.LabelMcpServicePort], service.Annotations[common.AnnotationMcpConnectURI])
	transportType := service.Labels[common.LabelMcpTransportType]

	logger.Infof("Endpoint: %v", endpoint)
	logger.Infof("TransportType: %v", transportType)

	mcpClient, err := InitClient(endpoint, transportType)
	if err != nil {
		logger.Errorf("failed to init mcp client, err: %v", err)
		return nil, err
	}

	tools, err := mcpClient.ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		logger.Errorf("failed to list tools, err: %v", err)
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}
	removeAnyOf(tools)

	return tools.Tools, nil
}

func loadServerConfig(logger *logrus.Entry, service *diceyml.Service, host string) string {
	logger = logger.WithField("action", "loadServerConfig")
	resp, err := http.Get(fmt.Sprintf("http://%s:%s/server/config", host, service.Labels[common.LabelMcpServicePort]))
	if err != nil || resp.StatusCode != 200 {
		logger.Warning("get server config error: %v\n", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Infof("failed to read request body: %v\n", err)
		return ""
	}

	return string(body)
}
