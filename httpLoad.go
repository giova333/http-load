package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

func loadTest(definition *LoadDefinition) {
	fmt.Println("Starting load test:")
	c := make(chan *RequestResult)
	var wg sync.WaitGroup
	requestsPerGoroutine := definition.numOfRequests / definition.concurrencyLevel
	remainingRequests := definition.numOfRequests % definition.concurrencyLevel

	for i := 0; i < definition.concurrencyLevel; i++ {
		wg.Add(1)
		times := requestsPerGoroutine
		if remainingRequests > 0 {
			times += 1
			remainingRequests--
		}
		go callTargetUrl(definition, times, &wg, c)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	printResult(c)
}

func printResult(c chan *RequestResult) {
	idx := 1
	successfulReq := 0
	failedReq := 0
	var durationOfAllRequestsNanos int64
	for requestResult := range c {
		durationOfAllRequestsNanos += requestResult.duration.Nanoseconds()
		if requestResult.status == Successful {
			successfulReq++
		} else {
			failedReq++
		}

		fmt.Printf("%d: Status: %s Time: %s\n", idx, requestResult.status, requestResult.duration)
		idx++
	}
	averageDurationOfRequest := time.Duration(durationOfAllRequestsNanos / int64(successfulReq+failedReq))
	fmt.Println("--------------------------------------------")
	fmt.Printf("Number of sucessfull requests: %d\n", successfulReq)
	fmt.Printf("Number of failed requests: %d\n", failedReq)
	fmt.Printf("Average duration of request: %s\n", averageDurationOfRequest)
}

func callTargetUrl(definition *LoadDefinition, times int, wg *sync.WaitGroup, c chan *RequestResult) {
	defer (*wg).Done()
	for i := 0; i < times; i++ {
		request, _ := http.NewRequest(definition.method, definition.targetUrl, bytes.NewBufferString(definition.payload))
		for _, header := range definition.headers {
			nameToValue := strings.Split(header, ":")
			request.Header.Set(nameToValue[0], nameToValue[1])
		}
		client := &http.Client{}
		start := time.Now()
		resp, err := client.Do(request)
		elapsed := time.Since(start)
		if err != nil || resp.StatusCode >= 400 {
			c <- NewRequestResult(Failed, elapsed)
		} else {
			c <- NewRequestResult(Successful, elapsed)
		}
	}
}
