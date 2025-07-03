package jdk

import (
	"testing"
)

func TestUpdateAlternativesSwitcher_IsVersionSupported(t *testing.T) {
	switcher := NewUpdateAlternativesSwitcher()

	tests := []struct {
		version  string
		expected bool
	}{
		{"8", true},
		{"11", true},
		{"17", true},
		{"21", true},
		{"7", false},
		{"12", false},
		{"", false},
	}

	for _, test := range tests {
		result := switcher.IsVersionSupported(test.version)
		if result != test.expected {
			t.Errorf("IsVersionSupported(%s) = %v; expected %v", test.version, result, test.expected)
		}
	}
}

func TestUpdateAlternativesSwitcher_DeriveJavacPath(t *testing.T) {
	switcher := NewUpdateAlternativesSwitcher()

	tests := []struct {
		javaPath     string
		expectedPath string
	}{
		{
			"/usr/lib/jvm/java-8-openjdk-amd64/jre/bin/java",
			"/usr/lib/jvm/java-8-openjdk-amd64/bin/javac",
		},
		{
			"/usr/lib/jvm/java-11-openjdk-amd64/bin/java",
			"/usr/lib/jvm/java-11-openjdk-amd64/bin/javac",
		},
		{
			"/usr/lib/jvm/java-17-openjdk-arm64/bin/java",
			"/usr/lib/jvm/java-17-openjdk-arm64/bin/javac",
		},
	}

	for _, test := range tests {
		result := switcher.DeriveJavacPath(test.javaPath)
		if result != test.expectedPath {
			t.Errorf("DeriveJavacPath(%s) = %s; expected %s", test.javaPath, result, test.expectedPath)
		}
	}
}

func TestNewUpdateAlternativesSwitcher(t *testing.T) {
	switcher := NewUpdateAlternativesSwitcher()

	if switcher == nil {
		t.Error("NewUpdateAlternativesSwitcher() returned nil")
	}

	if switcher.supportedVersions == nil {
		t.Error("supportedVersions map is nil")
	}

	expectedVersions := []string{"8", "11", "17", "21"}
	for _, version := range expectedVersions {
		if !switcher.IsVersionSupported(version) {
			t.Errorf("Expected version %s to be supported", version)
		}
	}
}
