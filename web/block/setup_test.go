package block

import (
	"testing"
	"time"
)

type MockTimeProvider struct {
	MockTime time.Time
}

func (m *MockTimeProvider) Now() time.Time {
	return m.MockTime
}

func SetupTest(m *testing.M) {
	m.Run()
}
