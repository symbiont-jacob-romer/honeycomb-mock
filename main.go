package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/klauspost/compress/zstd"
)

type trace struct {
	Service string `json:"service_name"`
	Name    string `json:"name"`
	TraceID string `json:"trace.trace_id"`
	SpanID  string `json:"trace.span_id"`
}

type event struct {
	Data trace `json:"data"`
}

type traceCache struct {
	mu     sync.Mutex
	traces []trace
}

func (c *traceCache) Add(t trace) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Printf("appending %+v. New cache size: %v\n", t, len(c.traces)+1)
	c.traces = append(c.traces, t)
}

func (c *traceCache) Dump() []byte {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Printf("dumping %v traces from cache\n", len(c.traces))
	m, _ := json.Marshal(c.traces)
	return m
}

func setupRouter(r chi.Router, cache *traceCache) chi.Router {
	r.Get("/traces", func(w http.ResponseWriter, r *http.Request) {
		w.Write(cache.Dump())
	})
	r.Post("/1/batch/{dataset}", func(w http.ResponseWriter, r *http.Request) {
		var decodedBody []byte
		var jsonDecodedEvents []event
		bodyBytes, _ := ioutil.ReadAll(r.Body)

		switch r.Header.Get("Content-Encoding") {
		case "zstd":
			var decoder, _ = zstd.NewReader(nil)
			decodedBody, _ = decoder.DecodeAll(bodyBytes, nil)
		case "gzip":
			gr, _ := gzip.NewReader(bytes.NewBuffer(bodyBytes))
			defer gr.Close()
			decodedBody, _ = ioutil.ReadAll(gr)
		default:
			decodedBody = bodyBytes
		}

		if err := json.Unmarshal(decodedBody, &jsonDecodedEvents); err != nil {
			panic(err)
		}

		for _, ev := range jsonDecodedEvents {
			cache.Add(ev.Data)
		}
	})
	return r
}

func main() {
	cache := traceCache{}
	cache.traces = []trace{}
	r := chi.NewRouter()
	setupRouter(r, &cache)
	s := &http.Server{
		Addr:    "localhost:8126",
		Handler: r,
	}
	fmt.Println(s.ListenAndServe())
}
