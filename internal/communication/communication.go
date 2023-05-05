package communication

import (
	"bytes"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/BastienFaivre/byzantine-shield/internal/aggregator"
	"github.com/BastienFaivre/byzantine-shield/internal/config"
	"github.com/BastienFaivre/byzantine-shield/internal/jsonrpc"
	"github.com/BastienFaivre/byzantine-shield/internal/types"
	"github.com/BastienFaivre/byzantine-shield/internal/utils"
)

type Proxy struct {
	Config config.Config
}

func NewProxy(config config.Config) *Proxy {
	return &Proxy{
		Config: config,
	}
}

func (p *Proxy) HandleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Request from %s: %s", r.RemoteAddr, string(body))

	rpcCall, err := jsonrpc.ParseCall(string(body))
	if err != nil {
		log.Println(err)
		return
	}

	responses := make(chan types.HttpResponse, len(p.Config.Nodes))
	var wg sync.WaitGroup

	for _, node := range p.Config.Nodes {
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
				Timeout: time.Duration(p.Config.Timeout) * time.Millisecond,
			}
			resp, err := client.Do(req)
			if err != nil {
				// timeout is considered as a valid response
				if err, ok := err.(net.Error); ok && err.Timeout() {
					id, err := jsonrpc.GetIdFromCall(rpcCall)
					if err != nil {
						log.Println(err)
						return
					}
					responses <- types.HttpResponse{
						Node: node,
						Body: []byte(jsonrpc.BuildTimeoutErrorResponse(id)),
					}
				}
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

	var res string
	if utils.Contains(p.Config.NonAggregateMethods, rpcCall.Method) {
		id, err := jsonrpc.GetIdFromCall(rpcCall)
		if err != nil {
			log.Println(err)
			return
		}
		res, err = aggregator.JoinResults(len(p.Config.Nodes), responses, id)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		res, err = aggregator.AggregateResults(len(p.Config.Nodes), responses)
		if err != nil {
			log.Println(err)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}
