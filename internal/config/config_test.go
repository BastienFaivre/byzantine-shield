package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config-test.yml")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testConfig := `nodes:
  - 192.168.1.1:8545
  - 192.168.1.2:8545
  - 192.168.1.3:8545
timeout: 1000
`

	if _, err := tmpFile.Write([]byte(testConfig)); err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}

	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if len(config.Nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(config.Nodes))
	}

	expectedNodes := []string{
		"192.168.1.1:8545",
		"192.168.1.2:8545",
		"192.168.1.3:8545",
	}
	for i, expectedNode := range expectedNodes {
		if config.Nodes[i] != expectedNode {
			t.Errorf("expected node %d to be %s, got %s", i, expectedNode, config.Nodes[i])
		}
	}

	if config.Timeout != 1000 {
		t.Errorf("expected timeout of 1000, got %d", config.Timeout)
	}
}
