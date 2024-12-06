package limit

import (
	"context"
	"time"
)

type SpanFunc func() (done bool, err error)

const minBusy = time.Millisecond * 20

type CPULimitTask struct {
	start         time.Time
	busy          time.Duration
	busySpanStart time.Time
	busySpan      time.Duration
	cpuLimit      float32
}

func NewCPULimitTask(cpuLimit float32) *CPULimitTask {
	now := time.Now()
	task := &CPULimitTask{
		start:         now,
		busySpanStart: now,
		cpuLimit:      cpuLimit,
	}
	return task
}

func (task *CPULimitTask) FillIdle(ctx context.Context) error {
	task.busy += time.Since(task.busySpanStart)
	expectTotal := time.Duration(float32(task.busy) / task.cpuLimit)
	if idle := expectTotal - time.Since(task.start); idle > 0 {
		select {
		case <-ctx.Done():
			task.busySpanStart = time.Now()
			return ctx.Err()
		case <-time.After(idle):
			task.busySpanStart = time.Now()
			return nil
		}
	}
	return nil
}

func (task *CPULimitTask) Do(ctx context.Context, span SpanFunc) error {
	for {
		var done bool
		var spanErr error
		select {
		case <-ctx.Done():
			spanErr = ctx.Err()
		default:
		}
		if spanErr == nil {
			done, spanErr = span()
		}

		busySpan := time.Since(task.busySpanStart)
		if busySpan < minBusy {
			if done || spanErr != nil {
				return spanErr
			}
			continue
		}

		task.busy += busySpan
		expectTotal := time.Duration(float32(task.busy) / task.cpuLimit)
		if idle := expectTotal - time.Since(task.start); idle > 0 {
			select {
			case <-ctx.Done():
				done, spanErr = false, ctx.Err()
			case <-time.After(idle):
			}
			task.busySpanStart = time.Now()
		}

		if done || spanErr != nil {
			return spanErr
		}
	}
}
