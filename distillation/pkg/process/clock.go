package process

import (
	"time"
)

type Clock interface {
	Unix() int64 // Unix returns seconds since 01.01.1970 UTC
}

type clock struct {
}

func (*clock) Unix() int64 {
	return time.Now().Unix()
}
