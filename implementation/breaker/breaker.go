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

type Impl struct {
	Wrapped aurestclientapi.Client

	CB *gobreaker.CircuitBreaker
	RequestTimeout time.Duration
}

func New(
	wrapped aurestclientapi.Client,
	circuitBreakerName string,
	maxNumRequestsInHalfOpenState uint32,
	counterClearingIntervalWhileClosed time.Duration,
	timeUntilHalfopenAfterOpen time.Duration,
	requestTimeout time.Duration,
) aurestclientapi.Client {
	settings := gobreaker.Settings{
		Name:          circuitBreakerName,
		MaxRequests:   maxNumRequestsInHalfOpenState,
		Interval:      counterClearingIntervalWhileClosed,
		Timeout:       timeUntilHalfopenAfterOpen, // NOT the request timeout
		ReadyToTrip:   nil,
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			if aulogging.Logger != nil {
				aulogging.Logger.NoCtx().Warn().Printf("circuit breaker %s state change %s -> %s", name, from.String(), to.String())
			}
		},
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}
			return aurestnontripping.Is(err)
		},
	}
	instance := &Impl{
		Wrapped:        wrapped,
		CB:             gobreaker.NewCircuitBreaker(settings),
		RequestTimeout: requestTimeout,
	}

	if aulogging.Logger != nil {
		aulogging.Logger.NoCtx().Info().Printf("circuit breaker %s set up", circuitBreakerName)
	}

	return instance
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
	return err
}
