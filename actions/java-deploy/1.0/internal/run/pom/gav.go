package pom

import (
	"encoding/xml"
	"os"
)

const (
	GroupID    = "groupId"
	ArtifactID = "artifactId"
	Version    = "version"
)

type GAV struct {
	GroupID    string
	ArtifactID string
	Version    string
}

type Project struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

func GetGAV(pomFilePath string) (*GAV, error) {
	f, err := os.Open(pomFilePath)
	if err != nil {
		return nil, err
	}
	decoder := xml.NewDecoder(f)

	var p Project
	if err := decoder.Decode(&p); err != nil {
		return nil, err
	}
	gav := GAV{
		GroupID:    p.GroupID,
		ArtifactID: p.ArtifactID,
		Version:    p.Version,
	}
	return &gav, nil
}
