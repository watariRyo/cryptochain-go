package configs

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoad(t *testing.T) {
	expectedEnv := "TEST"
	expectedHost := "http://backend-1:8080"
	expectedServerDefaultPort := "3000"
	expectedServerPort := "4000"
	expectedServerCorsOrigins := "5000"
	expectedRedisHost := "redis-test"
	expectedRedisPort := "6379"
	expectedRedisPassword := "password"
	expectedRedisKey := "my-cryptochain"

	want := &Config{
		Env:  expectedEnv,
		Host: expectedHost,
		Server: Server{
			DefaultPort: expectedServerDefaultPort,
			Port:        expectedServerPort,
			CorsOrigins: expectedServerCorsOrigins,
		},
		Redis: Redis{
			Host:     expectedRedisHost,
			Port:     expectedRedisPort,
			Password: expectedRedisPassword,
			Key:      expectedRedisKey,
		},
	}

	os.Setenv("ENV", expectedEnv)
	os.Setenv("SERVER_DEFAULTPORT", expectedServerDefaultPort)
	os.Setenv("SERVER_PORT", expectedServerPort)
	os.Setenv("SERVER_CORSORIGINS", expectedServerCorsOrigins)
	os.Setenv("REDIS_HOST", expectedRedisHost)
	os.Setenv("REDIS_PASSWORD", expectedRedisPassword)

	cfg, err := Load()
	if err != nil {
		t.Errorf("Load failed: %v", err)
	}

	if d := cmp.Diff(cfg, want); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}
