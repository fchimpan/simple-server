package clock

import "time"

type Clocker interface {
	Now() time.Time
}

type RealClocker struct{}

func (RealClocker) Now() time.Time {
	return time.Now()
}

type FakeClocker struct{}

func (FakeClocker) Now() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}
