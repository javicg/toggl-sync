package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

type MockHTTPResponse struct {
	response   string
	statusCode int
}

type MockHTTPServer struct {
	validators map[string]func(*http.Request)
	responses  map[string]*MockHTTPResponse
}

func NewHTTPServer() *MockHTTPServer {
	return &MockHTTPServer{
		responses:  make(map[string]*MockHTTPResponse),
		validators: make(map[string]func(*http.Request)),
	}
}

func (server *MockHTTPServer) StubAPI(stubbing *Stubbing) *MockHTTPServer {
	if stubbing.Endpoint != "" {
		server.responses[stubbing.Endpoint] = &MockHTTPResponse{
			response:   stubbing.ResponseBody,
			statusCode: stubbing.ResponseCode,
		}
		if stubbing.RequestValidator != nil {
			server.validators[stubbing.Endpoint] = stubbing.RequestValidator
		}
	}
	return server
}

type Stubbing struct {
	Endpoint         string
	RequestValidator func(*http.Request)
	ResponseCode     int
	ResponseBody     string
}

func (server *MockHTTPServer) Create() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling request: %s %s", r.Method, r.URL.String())
		for endpoint, stub := range server.responses {
			if strings.Contains(r.URL.String(), endpoint) {
				log.Printf("Found matching stub: %s -> (%d, %s)", endpoint, stub.statusCode, stub.response)
				server.validateRequestIfNeeded(endpoint, r)
				server.handleMatchingStub(endpoint, w)
				return
			}
		}
		log.Print("No matching stub found. Returning HTTP 404 (Not Found)")
		w.WriteHeader(http.StatusNotFound)
	}))
}

func (server *MockHTTPServer) validateRequestIfNeeded(endpoint string, r *http.Request) {
	if server.validators[endpoint] != nil {
		server.validators[endpoint](r)
	}
}

func (server *MockHTTPServer) handleMatchingStub(endpoint string, w http.ResponseWriter) {
	stub := server.responses[endpoint]
	if stub.response != "" {
		_, err := fmt.Fprintln(w, stub.response)
		if err != nil {
			log.Printf("Writing response failed [%s]. Returning HTTP 500 (Internal Server Error)", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(stub.statusCode)
	}
}

func AsJSONString(something interface{}) string {
	bytes, _ := json.Marshal(something)
	return string(bytes)
}
