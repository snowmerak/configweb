package config

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func Test_IterateData(t *testing.T) {
	sampleYamlData := `STORAGE:
    REDIS: "$SHARED_REDIS"
    POSTGRES: "$SHARED_POSTGRES"
    LOGGING:
        APP_NAME: "practice"
        POLICY: "$LOG_POLICY"
        KAFKA: "$LOG_KAFKA"
NUMBERS:
    - 1
    - 2
`

	value := map[string]any{}

	if err := yaml.Unmarshal([]byte(sampleYamlData), &value); err != nil {
		t.Fatalf("failed to unmarshal data: %v", err)
	}

	data, keys := iterateData(value)

	t.Log("data:")
	for k, v := range data {
		t.Logf("key: %v, value: %v", k, v)
	}

	t.Log("keys:")
	for _, key := range keys {
		t.Logf("key: %v, value: %v", key.First(), key.Second())
	}

	if len(keys) != 4 {
		t.Fatalf("expected 4 keys, got %d", len(keys))
	}
}
