package aggregator

import (
	"byzantine-shield/internal/types"
	"log"
)

func AggregateResults(nbNodes int, results chan types.HttpResponse) (string, error) {
	responseCount := make(map[string]int)

	for result := range results {
		responseCount[result.Body]++
	}

	maxCount := 0
	var aggregatedResponse string

	for response, count := range responseCount {
		if count > maxCount {
			maxCount = count
			aggregatedResponse = response
		}
	}

	log.Printf("Aggregated response: %s", aggregatedResponse)
	log.Printf("Response count: %v/%v (%.2f%%)", maxCount, nbNodes, float64(maxCount)/float64(nbNodes)*100)

	return aggregatedResponse, nil
}
