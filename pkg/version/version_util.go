package version

var agentHistoryVersions = []string{
	"3.16",
	"3.17",
	"3.18",
	"3.19",
	"3.20",
	"3.21",
	"4.0",
	"1.1",
	"1.2",
	"1.3",
	"1.4",
	"1.5",
	"1.6",
	"2.0",
}

func IsHistoryVersion(version string) bool {
	for _, historyVersion := range agentHistoryVersions {
		if historyVersion == version {
			return true
		}
	}
	return false
}
