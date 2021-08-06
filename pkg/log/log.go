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

func AddLineDelimiter(prefix ...string) {
	var _prefix string
	if len(prefix) > 0 {
		_prefix = prefix[0]
	}
	logrus.Printf("%s==========", _prefix)
}

func AddNewLine(num ...int) {
	_num := 1
	if len(num) > 0 {
		_num = num[0]
	}
	for i := 0; i < _num; i++ {
		logrus.Println()
	}
}
