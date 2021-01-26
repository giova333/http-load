package main

import (
	"strings"
	"time"
)

type LoadDefinition struct {
	numOfRequests    int
	concurrencyLevel int
	targetUrl        string
	method           string
	payload          string
	headers          headers
}

func NewLoadDefinition(numOfRequests int, concurrencyLevel int, targetUrl string, method string, payload string, headers headers) *LoadDefinition {
	return &LoadDefinition{numOfRequests: numOfRequests, concurrencyLevel: concurrencyLevel, targetUrl: targetUrl, method: method, payload: payload, headers: headers}
}

const (
	Successful = "Successful"
	Failed     = "Failed"
)

type RequestResult struct {
	status   string
	duration time.Duration
}

func NewRequestResult(status string, duration time.Duration) *RequestResult {
	return &RequestResult{status: status, duration: duration}
}

type headers []string

func (h *headers) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func (h *headers) String() string {
	return strings.Join(*h, ",")
}
