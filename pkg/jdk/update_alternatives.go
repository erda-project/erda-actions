package jdk

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var defaultSupportedVersions = map[string]bool{
	"8":  true,
	"11": true,
	"17": true,
	"21": true,
}

// UpdateAlternativesSwitcher update-alternatives Java version switcher implementation
type UpdateAlternativesSwitcher struct {
	supportedVersions map[string]bool
}

func NewUpdateAlternativesSwitcher() *UpdateAlternativesSwitcher {
	return &UpdateAlternativesSwitcher{
		supportedVersions: defaultSupportedVersions,
	}
}

// SwitchToVersion switches to the specified Java version
func (s *UpdateAlternativesSwitcher) SwitchToVersion(version string) error {
	if !s.IsVersionSupported(version) {
		return fmt.Errorf("unsupported Java version: %s", version)
	}

	fmt.Printf("Switching to Java %s...\n", version)

	// 1. Find available Java version
	javaPath, err := s.findJavaPath(version)
	if err != nil {
		return fmt.Errorf("java %s version not found: %v", version, err)
	}

	fmt.Printf("Found Java %s path: %s\n", version, javaPath)

	// 2. Set java
	if err := s.runCommand(fmt.Sprintf("update-alternatives --set java \"%s\"", javaPath)); err != nil {
		return fmt.Errorf("failed to set java: %v", err)
	}

	// 3. Intelligently derive and set javac path
	javacPath := s.deriveJavacPath(javaPath)
	if err := s.runCommand(fmt.Sprintf("update-alternatives --set javac \"%s\"", javacPath)); err != nil {
		return fmt.Errorf("failed to set javac: %v", err)
	}

	fmt.Printf("Successfully switched to Java %s\n", version)
	return nil
}

// GetCurrentJavaHome gets the JAVA_HOME of the current Java version
func (s *UpdateAlternativesSwitcher) GetCurrentJavaHome() (string, error) {
	cmd := exec.Command("/bin/bash", "-c", "java -XshowSettings:properties -version 2>&1 | grep 'java.home' | awk -F= '{print $2}' | tr -d ' '")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get JAVA_HOME: %v", err)
	}

	javaHome := strings.TrimSpace(string(output))
	if javaHome == "" {
		return "", errors.New("java.home is empty")
	}

	return javaHome, nil
}

// IsVersionSupported checks if the specified version is supported
func (s *UpdateAlternativesSwitcher) IsVersionSupported(version string) bool {
	return s.supportedVersions[version]
}

// ListAvailableVersions lists all available Java versions
func (s *UpdateAlternativesSwitcher) ListAvailableVersions() ([]string, error) {
	cmd := exec.Command("/bin/bash", "-c", "update-alternatives --list java")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get Java version list: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	var versions []string
	versionSet := make(map[string]bool)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract version number from path, e.g.: /usr/lib/jvm/java-8-openjdk-amd64/jre/bin/java -> 8
		for version := range s.supportedVersions {
			if strings.Contains(line, fmt.Sprintf("java-%s", version)) && !versionSet[version] {
				versions = append(versions, version)
				versionSet[version] = true
				break
			}
		}
	}

	return versions, nil
}

// findJavaPath finds the Java path for the specified version
func (s *UpdateAlternativesSwitcher) findJavaPath(version string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("update-alternatives --list java | grep 'java-%s' | head -n 1", version))
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	javaPath := strings.TrimSpace(string(output))
	if javaPath == "" {
		return "", fmt.Errorf("java %s path is empty", version)
	}

	return javaPath, nil
}

// DeriveJavacPath derives javac path from java path (exported for testing)
func (s *UpdateAlternativesSwitcher) DeriveJavacPath(javaPath string) string {
	return s.deriveJavacPath(javaPath)
}

// deriveJavacPath derives javac path from java path
func (s *UpdateAlternativesSwitcher) deriveJavacPath(javaPath string) string {
	if strings.Contains(javaPath, "/jre/bin/java") {
		// JRE path: /usr/lib/jvm/java-8-xxx/jre/bin/java -> /usr/lib/jvm/java-8-xxx/bin/javac
		return strings.Replace(javaPath, "/jre/bin/java", "/bin/javac", 1)
	} else {
		// JDK path: /usr/lib/jvm/java-11-xxx/bin/java -> /usr/lib/jvm/java-11-xxx/bin/javac
		return strings.Replace(javaPath, "/bin/java", "/bin/javac", 1)
	}
}

// runCommand executes command
func (s *UpdateAlternativesSwitcher) runCommand(cmd string) error {
	command := exec.Command("/bin/bash", "-c", cmd)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
