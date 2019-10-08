package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	libhoney "github.com/honeycombio/libhoney-go"
)

func TestMain(t *testing.T) {
	logger := log.New(os.Stdout, "", 1)
	libhoney.Init(libhoney.Config{
		WriteKey:   "abcd",
		Dataset:    "testd",
		APIHost:    "http://localhost:9997",
		SampleRate: 1,
		Logger:     logger,
	})

	cache := traceCache{}
	r := chi.NewRouter()
	setupRouter(r, &cache)
	s := &http.Server{
		Addr:    "localhost:9997",
		Handler: r,
	}

	go func() {
		fmt.Println(s.ListenAndServe())
	}()

	// Send an event to our mock Honeycomb API
	ev := libhoney.NewEvent()
	ev.AddField("trace.trace_id", "20348102384")
	ev.AddField("trace.span_id", "20348102384")
	ev.AddField("service_name", "api-server")
	ev.AddField("name", "testop")
	ev.Send()
	libhoney.Close()

	// Get a dump of all events
	resp, err := http.Get("http://localhost:9997/traces")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var dump []trace
	json.Unmarshal(bodyBytes, &dump)
	fmt.Println(dump)
	assert.Equal(t, len(dump), 1)
	assert.Equal(t, dump[0].Name, "testop")
	assert.Equal(t, dump[0].Service, "api-server")
	s.Close()
}
