package wallets

import (
	"testing"
	"time"
)

const layout = "2006-01-02T15:04:05.000000Z"

type MockTimeProvider struct {
	MockTime time.Time
}

func (m *MockTimeProvider) Now() time.Time {
	return m.MockTime
}

func (m *MockTimeProvider) NowMicroString() string {
	return m.Now().Format(layout)
}

func SetupTest(m *testing.M) {
	m.Run()
}
