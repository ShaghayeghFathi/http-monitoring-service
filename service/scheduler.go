package service

import (
	"errors"
	"time"
)

type Scheduler struct {
	Mnt  *Monitor
	Stop chan struct{}
}


func NewScheduler(mnt *Monitor) (*Scheduler, error) {
	sch := &Scheduler{Stop: make(chan struct{})}
	if mnt != nil {
		sch.Mnt = mnt
		return sch, nil
	}
	return nil, errors.New("Monitor cannot be null")
}


func (sch *Scheduler) WorkInIntervals(d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		for {
			select {
			case <-ticker.C:
				sch.Mnt.Work()
			case <-sch.Stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func (sch *Scheduler) StopSchedule() {
	close(sch.Stop)
}
