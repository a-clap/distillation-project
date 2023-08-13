package backend

import (
	"time"
)

func (b *Backend) NTPGet() (bool, error) {
	if b.time == nil {
		return false, ErrNotImplemented
	}
	return b.time.NTP()
}

func (b *Backend) NTPSet(enable bool) error {
	if b.time == nil {
		return ErrNotImplemented
	}
	return b.time.SetNTP(enable)
}

func (b *Backend) TimeSet(ts int64) error {
	if b.time == nil {
		return ErrNotImplemented
	}
	return b.time.SetNow(time.UnixMilli(ts))
}

func (b *Backend) Now() (int64, error) {
	if b.time == nil {
		return 0, ErrNotImplemented
	}
	t, err := b.time.Now()
	if err != nil {
		return 0, err
	}
	return t.UnixMilli(), nil
}
