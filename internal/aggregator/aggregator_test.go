package aggregator

import (
	"testing"

	"github.com/BastienFaivre/byzantine-shield/internal/types"
)

func TestAggregateResults(t *testing.T) {
	results := make(chan types.HttpResponse, 3)
	results <- types.HttpResponse{Node: "node1", Body: `{"id":1,"jsonrpc":"2.0","result":42}`}
	results <- types.HttpResponse{Node: "node2", Body: `{"id":1,"jsonrpc":"2.0","result":42}`}
	results <- types.HttpResponse{Node: "node3", Body: `{"id":1,"jsonrpc":"2.0","result":24}`}
	close(results)

	aggregatedResponse, err := AggregateResults(3, results)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expected := `{"id":1,"jsonrpc":"2.0","result":42}`
	if aggregatedResponse != expected {
		t.Errorf("expected %s, got %s", expected, aggregatedResponse)
	}
}
