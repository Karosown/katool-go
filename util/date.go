package util

import (
	"time"
)

func BeginEndTimeComputed(run func()) int64 {
	begin := time.Now().UnixNano()
	run()
	return time.Now().UnixNano() - begin
}
