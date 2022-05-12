package coverage

import (
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/erda-project/erda-proto-go/dop/qa/unittest/pb"
	"github.com/erda-project/erda/apistructs"
)

func GetJacocoFiles(path, prefix string) []string {
	var results []string
	filepath.Walk(path, func(fullPath string, info fs.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if (strings.HasSuffix(info.Name(), ".xml") &&
			strings.HasPrefix(info.Name(), prefix)) || info.Name() == "jacoco.xml" {
			results = append(results, fullPath)
		}
		return nil
	})
	return results
}

func BatchGenCoverageReport(paths []string) ([]*pb.CodeCoverageNode, error) {
	var nodes []*pb.CodeCoverageNode
	for _, path := range paths {
		node, err := GenCoverageReport(path)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node...)
	}
	return nodes, nil
}

func GenCoverageReport(path string) ([]*pb.CodeCoverageNode, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	all, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	codeReport, err := apistructs.ConvertXmlToReport(all)
	if err != nil {
		return nil, err
	}
	codeReport.ProjectName = codeReport.Name
	nodes, _ := apistructs.ConvertReportToTree(codeReport)
	pbNodes := ConvertReportToPb(nodes)
	return pbNodes, nil
}

func ConvertReportToPb(nodes []*apistructs.CodeCoverageNode) []*pb.CodeCoverageNode {
	var pbNodes []*pb.CodeCoverageNode
	for _, n := range nodes {
		values := make([]float32, 0)
		for _, value := range n.Value {
			values = append(values, float32(value))
		}
		pbNodes = append(pbNodes, &pb.CodeCoverageNode{
			Value:    values,
			Name:     n.Name,
			Path:     n.Path,
			Children: ConvertReportToPb(n.Nodes),
		})
	}
	return pbNodes
}

func ConvertToolTipToPb(toolTip apistructs.ToolTip) *pb.ToolTip {
	return &pb.ToolTip{
		Formatter: toolTip.Formatter,
	}
}
