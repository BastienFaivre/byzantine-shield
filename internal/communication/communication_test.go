package communication

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BastienFaivre/byzantine-shield/internal/config"
)

func TestHandleRequest(t *testing.T) {
	nodes := make([]string, 3)

	// Start mock servers
	mockServers := make([]*httptest.Server, len(nodes))
	for i := range nodes {
		mockServers[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id":1,"jsonrpc":"2.0","result":"0x1"}`))
		}))
		defer mockServers[i].Close()
	}

	// Update nodes to use mock server URLs
	for i, server := range mockServers {
		nodes[i] = server.URL
	}

	// Initialize Proxy and test request
	proxy := NewProxy(config.Config{Nodes: nodes, Timeout: 5000})
	testRequest := `{"jsonrpc":"2.0","id":1,"method":"eth_getBalance","params":["0xc94770007dda54cF92009BFF0dE90c06F603a09f", "latest"]}`
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(testRequest)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(proxy.HandleRequest)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := `{"id":1,"jsonrpc":"2.0","result":"0x1"}`
	response := rr.Body.String()
	if response != expectedResponse {
		t.Errorf("Handler returned unexpected body: got %v want %v", response, expectedResponse)
	}
}
