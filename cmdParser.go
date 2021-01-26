package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func parseCmd() *LoadDefinition {
	numOfRequests := flag.Int("n", 1, "Number of requests to send.")
	concurrencyLevel := flag.Int("c", 1, "Parameter that describes concurrency level of requests")
	method := flag.String("m", "GET", "Method name")
	payloadRow := flag.String("p", "", "Payload raw content")
	fileWithPayload := flag.String("f", "", "File with payload")

	var headers headers
	flag.Var(&headers, "h", "Header (Repeatable)")

	flag.Parse()

	targetUrl := flag.Arg(0)

	if *payloadRow == "" && *fileWithPayload != "" {
		content, err := ioutil.ReadFile(*fileWithPayload)
		if err != nil {
			fmt.Printf("Invalid file path [%s]", *fileWithPayload)
			os.Exit(1)
		}
		*payloadRow = string(content)
	}

	validateRequestNum(*numOfRequests, *concurrencyLevel)
	validateUrl(targetUrl)
	validateMethod(*method)

	loadDefinition := NewLoadDefinition(*numOfRequests, *concurrencyLevel, targetUrl, *method, *payloadRow, headers)
	printLoadDefinition(loadDefinition)
	return loadDefinition
}

func printLoadDefinition(definition *LoadDefinition) {
	fmt.Println("--------------------------------------------")
	fmt.Println("Load description:")
	fmt.Printf("Target url: %s\n", definition.targetUrl)
	fmt.Printf("Number of requests: %d\n", definition.numOfRequests)
	fmt.Printf("Concurrency level: %d\n", definition.concurrencyLevel)
	fmt.Printf("Method: %s\n", definition.method)
	if definition.payload != "" {
		fmt.Printf("Payload: %s\n", definition.payload)
	}
	if len(definition.headers) != 0 {
		fmt.Printf("Headers: %s\n", definition.headers)
	}
	fmt.Println("--------------------------------------------")
}

func validateRequestNum(numOfRequests int, concurrencyLevel int) {
	if numOfRequests < 0 {
		fmt.Printf("Number of requests cannot be negative [%d]\n", numOfRequests)
		os.Exit(1)
	}
	if concurrencyLevel < 0 {
		fmt.Printf("Concurrency level [%d]\n", concurrencyLevel)
		os.Exit(1)
	}
	if numOfRequests < concurrencyLevel {
		fmt.Printf("Request concurrencyLevel [%d] must be less then or equal to total number of requests [%d]\n", concurrencyLevel, numOfRequests)
		os.Exit(1)
	}
}

func validateMethod(method string) {
	notAValidMethod := method != http.MethodGet &&
		method != http.MethodPost &&
		method != http.MethodPut &&
		method != http.MethodPatch &&
		method != http.MethodDelete &&
		method != http.MethodHead &&
		method != http.MethodOptions
	if notAValidMethod {
		fmt.Printf("Http method [%s] is not supported", method)
		os.Exit(1)
	}
}

func validateUrl(target string) {
	_, err := url.ParseRequestURI(target)
	if err != nil {
		fmt.Printf("Invalid target url [%s]", target)
		os.Exit(1)
	}
	u, err := url.Parse(target)
	if err != nil || u.Scheme == "" || u.Host == "" {
		fmt.Printf("Invalid target url [%s]", target)
		os.Exit(1)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		fmt.Printf("Invalid target schema [%s]. Supported schemas are [http/https]", u.Scheme)
		os.Exit(1)
	}
}
