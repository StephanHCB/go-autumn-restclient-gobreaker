package examplefullstack

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	aurestbreaker "github.com/StephanHCB/go-autumn-restclient-circuitbreaker/implementation/breaker"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	auresthttpclient "github.com/StephanHCB/go-autumn-restclient/implementation/httpclient"
	"net/http"
	"time"
)

func example() {
	// assumes you set up logging by importing one of the go-autumn-logging-xxx dependencies
	//
	// for this example, let's set up a logger that does nothing, so we don't pull in these dependencies here
	//
	// This of course makes the requestLoggingClient not work.
	aulogging.SetupNoLoggerForTesting()

	// ----- setup (can be done once during application startup) ----

	// 1. set up http client
	var timeout time.Duration = 0
	var customCACert []byte = nil
	var requestManipulator aurestclientapi.RequestManipulatorCallback = nil

	httpClient, _ := auresthttpclient.New(timeout, customCACert, requestManipulator)

	// 2. circuit breaker

	var circuitBreakerName string = "some-name"
	var maxNumRequestsInHalfOpenState uint32 = 100
	var counterClearingIntervalWhileClosed time.Duration = 5 * time.Minute
	var timeUntilHalfopenAfterOpen time.Duration = 60 * time.Second
	var requestTimeout time.Duration = 15 * time.Second

	circuitBreakerClient := aurestbreaker.New(httpClient, circuitBreakerName, maxNumRequestsInHalfOpenState, counterClearingIntervalWhileClosed, timeUntilHalfopenAfterOpen, requestTimeout)

	// ----- now make a request -----

	bodyDto := make(map[string]interface{})

	response := aurestclientapi.ParsedResponse{
		Body: &bodyDto,
	}
	err := circuitBreakerClient.Perform(context.Background(), http.MethodGet, "https://some.rest.api", nil, &response)
	if err != nil {
		return
	}

	// now bodyDto is filled with the response and response.Status and response.Header are also set.
}
