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

package main

import (
	"context"
	"log"
	"net/http"
	"otel-demo/app"
	"otel-demo/trace"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

func main() {
	tp := trace.Setup()

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	server := gin.Default()
	server.GET("/fib", fib)
	server.GET("/loop_fib", loop_fib)
	server.Run(":8888")
}

func fib(c *gin.Context) {
	ctx, app := start(c, "fib")
	app.Run(ctx)
}

func loop_fib(c *gin.Context) {
	ctx, app := start(c, "loop_fib")
	app.LoopRun(ctx)
}

func start(c *gin.Context, serviceName string) (context.Context, *app.App) {
	ctx, span := otel.Tracer(serviceName).Start(c, "start")
	defer span.End()

	stringN, exist := c.GetQuery("n")
	if !exist {
		span.SetStatus(codes.Error, "No required query parameter: n")
		c.String(http.StatusBadRequest, "%s", "No required query parameter n")
		return ctx, nil
	}
	n, err := strconv.Atoi(stringN)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.String(http.StatusBadRequest, "error: %v", err)
		return ctx, nil
	}

	return ctx, app.NewApp(c, serviceName, uint(n))
}
