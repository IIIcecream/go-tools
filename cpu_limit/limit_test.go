package limit

import (
	"context"
	"testing"
	"time"
)

func TestDoWithCPULimit(t *testing.T) {
	const singleTaskDur = time.Second
	const totalTaskDur = time.Second * 10

	ctx := context.Background()

	task := NewCPULimitTask(0.02)

	start := time.Now()
	for {
		singleStart := time.Now()
		err := task.Do(ctx, func() (done bool, err error) {
			for i := 0; i < 10; i++ {
			}
			if time.Since(singleStart) > singleTaskDur {
				return true, nil
			}
			return false, nil
		})
		if err != nil {
			t.Fatal(err)
		}

		if time.Since(start) > totalTaskDur {
			break
		}
	}
}
