package testing

type Status string

const (
	ExtraOsName    = "os.name"
	ExtraOsArch    = "os.arch"
	ExtraOsVersion = "os.ersion"
	TestNgFile     = "testng-results.xml"
	JunitFile      = "TEST-"
	Junit          = "junit"
	TestNg         = "testng"
	// StatusFailed represents a JUnit testcase that was run, but resulted in
	// a failure. Failures are violations of declared test expectations,
	// such as a failed assertion.
	StatusFailed Status = "failed"

	// StatusError represents a JUnit testcase that was run, but resulted in
	// an error. Errors are unexpected violations of the test itself, such as
	// an uncaught exception.
	StatusError Status = "error"

	// props
	ExtraMvnVersion  = "mvn.version"
	ExtraJavaVersion = "java.version"
	ExtraJavaVmName  = "java.vm.name"
)
