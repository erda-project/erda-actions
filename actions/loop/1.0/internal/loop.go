package main

import (
	"time"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/pkg/loop"
)

var loopTimes = 1

func Loop(o *Object) error {
	l := loop.New(
		loop.WithMaxTimes(o.conf.LoopMaxTimes),
		loop.WithInterval(o.conf.LoopInterval),
		loop.WithDeclineRatio(o.conf.LoopDeclineRatio),
		loop.WithDeclineLimit(o.conf.LoopDeclineLimit),
	)

	switch o.conf.Type {
	case HTTP:
		if err := l.Do(wrapDoOnce(o.doOnceHTTP, o.successFlag)); err != nil {
			return err
		}
		if o.successFlag != nil && *o.successFlag {
			return nil
		}
		return errors.Errorf("loop HTTP not success")

	case CMD:
		if err := l.Do(wrapDoOnce(o.doOnceCMD, o.successFlag)); err != nil {
			return err
		}
		if o.successFlag != nil && *o.successFlag {
			return nil
		}
		return errors.Errorf("loop CMD not success")

	default:
		return nil
	}
}

func wrapDoOnce(fn func() (bool, error), successFlag *bool) func() (bool, error) {
	return func() (bool, error) {
		logger.Printf("begin try %d times, timestamp: %s\n", loopTimes, time.Now().Format(time.RFC3339Nano))
		defer func() {
			logger.Printf("end try %d times, timestamp: %s\n", loopTimes, time.Now().Format(time.RFC3339Nano))
			logger.Println()
			loopTimes++
		}()
		success, err := fn()
		if success {
			logger.Println("loop success!")
			*successFlag = true
			return true, nil
		}
		if err != nil {
			logger.Printf("loop err: %v", err)
			return false, err
		}
		return false, nil
	}
}
