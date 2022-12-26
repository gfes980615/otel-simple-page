# Open Telemetry Go Demo

This repository is an open-telemetry demo based on [open-telemetry/opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go/tree/main/example/fib).

## Getting Started

1. Run app to generate tracing.

    ```bash
    go run main.go
    ```

## Fibonacci

The main program reads input `n` and returns the `n`-th Fibonacci number. If `n` > 93, it will return an error for testing purpose.

For more details, see [docs](https://opentelemetry.io/docs/instrumentation/go/getting-started/).
