package communication

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/BastienFaivre/byzantine-shield/internal/aggregator"
	"github.com/BastienFaivre/byzantine-shield/internal/config"
	"github.com/BastienFaivre/byzantine-shield/internal/types"
)

type Proxy struct {
	Nodes   []string
	Timeout int
}

func NewProxy(config config.Config) *Proxy {
	return &Proxy{
		Nodes:   config.Nodes,
		Timeout: config.Timeout,
	}
}

func (p *Proxy) HandleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Request from %s: %s", r.RemoteAddr, string(body))

	responses := make(chan types.HttpResponse, len(p.Nodes))
	var wg sync.WaitGroup

	for _, node := range p.Nodes {
		wg.Add(1)
		go func(node string) {
			defer wg.Done()

			req, err := http.NewRequest(r.Method, node, bytes.NewBuffer(body))
			if err != nil {
				log.Println(err)
				return
			}

			req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

			client := &http.Client{
				Timeout: time.Duration(p.Timeout) * time.Millisecond,
			}
			resp, err := client.Do(req)
			if err != nil {
				log.Println(err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				return
			}

			responses <- types.HttpResponse{
				Node: node,
				Body: body,
			}
		}(node)
	}

	wg.Wait()
	close(responses)

	res, err := aggregator.AggregateResults(len(p.Nodes), responses)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// add response body
	w.Write([]byte(res))
}
