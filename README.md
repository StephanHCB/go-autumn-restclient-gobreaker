# go-autumn-restclient-circuitbreaker

Adds [sony/gobreaker](https://github.com/sony/gobreaker) as a wrapper for the rest client. 

## About go-autumn

A collection of libraries for [enterprise microservices](https://github.com/StephanHCB/go-mailer-service/blob/master/README.md) in golang that
- is heavily inspired by Spring Boot / Spring Cloud
- is very opinionated
- names modules by what they do
- unlike Spring Boot avoids certain types of auto-magical behaviour
- is not a library monolith, that is every part only depends on the api parts of the other components
  at most, and the api parts do not add any dependencies.  

Fall is my favourite season, so I'm calling it go-autumn.

## About go-autumn-restclient

It's a rest client that also supports x-www-form-urlencoded.

## About go-autumn-restclient-circuitbreaker

This library adds another wrapper to the rest client which provides a circuit breaker.

We currently use [sony/gobreaker](https://github.com/sony/gobreaker), which is a very lightweight
implementation with a small dependency footprint (MIT licensed).

## Usage

TODO

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
