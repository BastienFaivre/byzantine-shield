package aggregator

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/BastienFaivre/byzantine-shield/internal/types"
)

func findKey(responseCount map[string]int, jsonObj map[string]interface{}) string {
	for key := range responseCount {
		var keyObj map[string]interface{}
		err := json.Unmarshal([]byte(key), &keyObj)
		if err == nil && reflect.DeepEqual(keyObj, jsonObj) {
			return key
		}
	}
	return ""
}

func AggregateResults(nbNodes int, responses chan types.HttpResponse) (string, error) {
	responseCount := make(map[string]int)

	for response := range responses {
		var jsonObj map[string]interface{}
		err := json.Unmarshal(response.Body, &jsonObj)
		if err != nil {
			log.Println(err)
		} else {
			key := findKey(responseCount, jsonObj)
			if key != "" {
				responseCount[key]++
			} else {
				// Marshal to get a compact string representation
				key_, err := json.Marshal(jsonObj)
				if err != nil {
					log.Println(err)
				} else {
					key = string(key_)
					responseCount[key] = 1
				}
			}
		}
	}

	maxCount := 0
	var aggregatedResponse string

	for response, count := range responseCount {
		if count > maxCount {
			maxCount = count
			aggregatedResponse = response
		}
	}

	log.Printf("Aggregated response (mode ratio: %.2f%% (%v/%v)): %s", float64(maxCount)/float64(nbNodes)*100, maxCount, nbNodes, aggregatedResponse)

	return aggregatedResponse, nil
}
