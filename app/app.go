// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// name is the Tracer name used to identify this instrumentation library.
const name = "fib"

// App is an Fibonacci computation application.
type App struct {
	lib string
	c   *gin.Context
	r   uint
}

// NewApp returns a new App.
func NewApp(c *gin.Context, svcName string, r uint) *App {
	return &App{
		lib: svcName,
		c:   c,
		r:   r,
	}
}

func (a *App) Run(ctx context.Context) error {
	newCtx, span := otel.Tracer(a.lib).Start(ctx, "Run")
	defer span.End()

	n, err := a.Poll(newCtx)
	if err != nil {
		return err
	}

	a.Write(newCtx, n)

	return nil
}

func (a *App) LoopRun(ctx context.Context) error {
	newCtx, span := otel.Tracer(a.lib).Start(ctx, "LoopRun")
	defer span.End()
	span.SetAttributes(attribute.String("loop to request: ", strconv.FormatUint(uint64(a.r), 10)))
	for i := 0; i <= int(a.r); i++ {
		a.Write(newCtx, uint(i))
	}

	return nil
}

// Poll asks a user for input and returns the request.
func (a *App) Poll(ctx context.Context) (uint, error) {
	_, span := otel.Tracer(a.lib).Start(ctx, "Poll")
	defer span.End()
	n := a.r
	// Store n as a string to not overflow an int64.
	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nStr))

	return n, nil
}

// Write writes the n-th Fibonacci number back to the user.
func (a *App) Write(ctx context.Context, n uint) {
	newCtx, span := otel.Tracer(a.lib).Start(ctx, "Write")
	defer span.End()

	f, err := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer(a.lib).Start(ctx, "Fibonacci")
		defer span.End()
		f, err := Fibonacci(n)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		nStr := strconv.FormatUint(uint64(f), 10)
		span.SetAttributes(attribute.String("result", nStr))
		return f, err
	}(newCtx)
	if err != nil {
		a.c.String(http.StatusOK, "Fibonacci(%d): %v, TraceID: %s\n", n, err, span.SpanContext().TraceID())
	} else {
		a.c.String(http.StatusOK, "Fibonacci(%d) = %d, TraceID: %s\n", n, f, span.SpanContext().TraceID())
	}
}
