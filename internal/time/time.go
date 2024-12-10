package time

import "time"

type TimeProvider interface {
	Now() time.Time
	NowMicroString() string
}

type RealTimeProvider struct{}

const layout = "2006-01-02T15:04:05.000000Z"

func (r *RealTimeProvider) Now() time.Time {
	return time.Now()
}

func (r *RealTimeProvider) NowMicroString() string {
	return r.Now().Format(layout)
}

func MicroParseString(timestamp time.Time) string {
	return timestamp.Format(layout)
}

func MicroParse(str string) (time.Time, error) {
	tm, err := time.Parse(layout, str)
	if err != nil {
		return tm, err
	}
	return tm, nil
}
