# go-autumn-restclient-circuitbreaker

Adds [sony/gobreaker](https://github.com/sony/gobreaker) as a wrapper for the rest client. 

## About go-autumn-restclient

It's a rest client that also supports x-www-form-urlencoded.

## About go-autumn-restclient-circuitbreaker

This library adds another wrapper to the rest client which provides a circuit breaker.

We currently use [sony/gobreaker](https://github.com/sony/gobreaker), which is a very lightweight
implementation with a small dependency footprint (MIT licensed).

## Usage

Change the setup of your rest client like this:

```
// [...]

var circuitBreakerName string = "some-name"
var maxNumRequestsInHalfOpenState uint32 = 100
var counterClearingIntervalWhileClosed time.Duration = 5 * time.Minute
var timeUntilHalfopenAfterOpen time.Duration = 60 * time.Second
var requestTimeout time.Duration = 15 * time.Second

circuitBreakerClient := aurestbreaker.New(httpClient, circuitBreakerName, maxNumRequestsInHalfOpenState, counterClearingIntervalWhileClosed, timeUntilHalfopenAfterOpen, requestTimeout)
```

You should usually insert the cbClient above the request logger and below the retryer.

## Logging

This library uses the [StephanHCB/go-autumn-logging](https://github.com/StephanHCB/go-autumn-logging) api for
logging framework independent logging.

### Library Authors

If you are writing a library, do NOT import any of the go-autumn-logging-* modules that actually bring in a logging library.
You will deprive application authors of their chance to pick the logging framework of their choice.

In your testing code, call `aulogging.SetupNoLoggerForTesting()` to avoid the nil pointer dereference.

### Application Authors

If you are writing an application, import one of the modules that actually bring in a logging library,
such as go-autumn-logging-zerolog. These modules will provide an implementation and place it in the Logger variable.

Of course, you can also provide your own implementation of the `LoggingImplementation` interface, just
set the `Logger` global singleton to an instance of your implementation.

Then just use the Logger, both during application runtime and tests.
