package aggregator

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/BastienFaivre/byzantine-shield/internal/jsonrpc"
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

func getResponseCount(responses chan types.HttpResponse) map[string]int {
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

	return responseCount
}

func AggregateResults(nbNodes int, responses chan types.HttpResponse) (string, error) {
	responseCount := getResponseCount(responses)

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

type joinResponse struct {
	ModeRatio string `json:"modeRatio"`
	Response  string `json:"response"`
}

func JoinResults(nbNodes int, responses chan types.HttpResponse, id int) (string, error) {
	responsesCount := getResponseCount(responses)

	var result []joinResponse
	for response, count := range responsesCount {
		result = append(result, joinResponse{
			ModeRatio: fmt.Sprintf("%.2f%% (%v/%v)", float64(count)/float64(nbNodes)*100, count, nbNodes),
			Response:  response,
		})
	}

	res := jsonrpc.BuildResponse(id, result)

	log.Printf("Aggregated response (join): %s", res)

	return res, nil
}
