package worker

import (
	"time"
)

type Options struct {
	Interval        time.Duration
	RunImmediately  bool
	RunConsequently bool
}

func DefaultOptions(interval time.Duration) *Options {
	return &Options{
		Interval:        interval,
		RunImmediately:  true,
		RunConsequently: false,
	}
}
