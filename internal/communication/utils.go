package communication

import (
	"encoding/json"
	"fmt"
)

func getIdFromJSONRPC(payload string) (int, error) {
	var jsonObj map[string]interface{}
	err := json.Unmarshal([]byte(payload), &jsonObj)
	if err != nil {
		return 0, err
	}

	id, ok := jsonObj["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("id not found")
	}

	return int(id), nil
}

func buildTimeoutErrorJSONRPC(id int) string {
	return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"error":{"code":-32000,"message":"Timeout (from byzantine-shield)"}}`, id)
}
