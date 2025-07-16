package calendar

import "time"

type ServiceIface interface {
	NowUTC() time.Time
	Now() time.Time
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (c *Service) NowUTC() time.Time {
	return time.Now().UTC()
}

func (c *Service) Now() time.Time {
	return time.Now()
}
