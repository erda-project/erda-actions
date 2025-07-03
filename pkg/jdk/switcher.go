package jdk

// VersionSwitcher Java version switcher interface
type VersionSwitcher interface {
	// SwitchToVersion switches to the specified Java version
	SwitchToVersion(version string) error

	// GetCurrentJavaHome gets the JAVA_HOME of the current Java version
	GetCurrentJavaHome() (string, error)

	// IsVersionSupported checks if the specified version is supported
	IsVersionSupported(version string) bool

	// ListAvailableVersions lists all available Java versions
	ListAvailableVersions() ([]string, error)
}
