package dice

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func DownloadFile(url string, destPath string) (err error) {
	logrus.Infof("begin download file: %s -> %s", url, destPath)
	defer func() {
		if err != nil {
			logrus.Errorf("end download file: %s -> %s, failed, err: %v", url, destPath, err)
		} else {
			logrus.Infof("end download file: %s -> %s, success", url, destPath)
		}
	}()
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	counter := &WriteCounter{}
	_, err = io.Copy(f, io.TeeReader(res.Body, counter))
	counter.PrintProgress(true)
	return err
}

type WriteCounter struct {
	Total          uint64
	lastReportTime int64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress(false)
	return n, nil
}

func (wc *WriteCounter) PrintProgress(forcePrint bool) {
	now := time.Now().UnixNano()
	if time.Duration(now-wc.lastReportTime) > (time.Second*3) || forcePrint {
		fmt.Printf("data transport  %s\n", ByteCountIEC(int64(wc.Total)))
		wc.lastReportTime = now
	}
}

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
