package aurestbreaker

import (
	"context"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	aurestclientapi "github.com/StephanHCB/go-autumn-restclient/api"
	aurestnontripping "github.com/StephanHCB/go-autumn-restclient/implementation/errors/nontrippingerror"
	"github.com/sony/gobreaker"
	"time"
)

// StateChangeCallbackFunction allows you to instrument the circuit breaker.
type StateChangeCallbackFunction func(circuitBreakerName string, state string)

type CountsCallbackFunction func(circuitBreakerName string, counts gobreaker.Counts)

type Impl struct {
	Wrapped aurestclientapi.Client

	Name           string
	CB             *gobreaker.CircuitBreaker
	RequestTimeout time.Duration

	StateChangeCallback StateChangeCallbackFunction
	CountsCallback      CountsCallbackFunction
}

func New(
	wrapped aurestclientapi.Client,
	circuitBreakerName string,
	maxNumRequestsInHalfOpenState uint32,
	counterClearingIntervalWhileClosed time.Duration,
	timeUntilHalfopenAfterOpen time.Duration,
	requestTimeout time.Duration,
) aurestclientapi.Client {
	instance := &Impl{
		Wrapped:             wrapped,
		Name:                circuitBreakerName,
		RequestTimeout:      requestTimeout,
		StateChangeCallback: doNothingStateChangeCallback,
		CountsCallback:      doNothingCountsCallback,
	}

	settings := gobreaker.Settings{
		Name:        circuitBreakerName,
		MaxRequests: maxNumRequestsInHalfOpenState,
		Interval:    counterClearingIntervalWhileClosed,
		Timeout:     timeUntilHalfopenAfterOpen, // NOT the request timeout
		ReadyToTrip: nil,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			if aulogging.Logger != nil {
				aulogging.Logger.NoCtx().Warn().Printf("circuit breaker %s state change %s -> %s", name, from.String(), to.String())
			}
			instance.StateChangeCallback(name, to.String())
		},
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}
			return aurestnontripping.Is(err)
		},
	}
	instance.CB = gobreaker.NewCircuitBreaker(settings)

	if aulogging.Logger != nil {
		aulogging.Logger.NoCtx().Info().Printf("circuit breaker %s set up", circuitBreakerName)
	}

	return instance
}

// Instrument adds instrumentation.
//
// Either of the callbacks may be nil.
func Instrument(
	client aurestclientapi.Client,
	stateChangeCallback StateChangeCallbackFunction,
	countsCallback CountsCallbackFunction,
) {
	cbClient, ok := client.(*Impl)
	if !ok {
		return
	}

	if stateChangeCallback != nil {
		cbClient.StateChangeCallback = stateChangeCallback
	}
	if countsCallback != nil {
		cbClient.CountsCallback = countsCallback
	}
}

func doNothingStateChangeCallback(_ string, _ string) {

}

func doNothingCountsCallback(_ string, _ gobreaker.Counts) {

}

func (c *Impl) Perform(ctx context.Context, method string, requestUrl string, requestBody interface{}, response *aurestclientapi.ParsedResponse) error {
	_, err := c.CB.Execute(func() (interface{}, error) {
		childCtx, cancel := context.WithTimeout(ctx, c.RequestTimeout)
		defer cancel()

		innerErr := c.Wrapped.Perform(childCtx, method, requestUrl, requestBody, response)
		if innerErr != nil {
			return nil, innerErr
		}

		if response.Status >= 500 {
			// ensure breaking
			return nil, fmt.Errorf("got http status %d", response.Status)
		}

		return nil, nil
	})
	c.CountsCallback(c.Name, c.CB.Counts())
	return err
}
