package util

import (
	"time"
)

func BeginEndTimeComputed(run func()) int64 {
	begin := time.Now().UnixNano()
	run()
	return time.Now().UnixNano() - begin
}

type PeriodTimeStamp struct {
	Start int64
	End   int64
}

type PeriodTime struct {
	Start time.Time
	End   time.Time
}

// GetPeriods 获取一些时间段，时间段的范围在start - end 之间，时间段的长度为 duration
func GetPeriods(start, end time.Time, duration time.Duration) []*PeriodTime {
	rv := make([]*PeriodTime, 0)
	temp := start
	for {
		if temp.After(end) || temp.Equal(end) {
			break
		}
		period := &PeriodTime{
			Start: temp,
			End:   temp.Add(duration),
		}
		temp = period.End
		rv = append(rv, period)
	}
	return rv
}
