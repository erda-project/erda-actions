package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

type actionLogFormatter struct {
	logrus.TextFormatter
}

func (f *actionLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	_bytes, err := f.TextFormatter.Format(entry)
	if err != nil {
		return nil, err
	}
	return append([]byte("[Action Log] "), _bytes...), nil
}

func Init() {
	// set logrus
	logrus.SetFormatter(&actionLogFormatter{
		logrus.TextFormatter{
			ForceColors:            true,
			DisableTimestamp:       true,
			DisableLevelTruncation: true,
		},
	})
	logrus.SetOutput(os.Stdout)
}
