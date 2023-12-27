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
	"fmt"
	"log"
	"math/rand"
	"otel-demo/lib"
	"otel-demo/trace"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	go_trace "go.opentelemetry.io/otel/trace"
)

func main() {
	tp := trace.Setup()

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	server := gin.Default()
	server.GET("/ubt", ubt)
	server.Run(":8888")
}

func ubt(c *gin.Context) {
	count, _ := c.GetQuery("count")
	cnt, _ := strconv.ParseInt(count, 10, 64)
	ubtTracingData(cnt)
}

func ubtTracingData(count int64) {
	st := time.Now()
	timeStart, _ := time.Parse("2006-01-02T15:04:05.000Z", "2023-01-01T00:00:00.000Z")
	timeEnd, _ := time.Parse("2006-01-02T15:04:05.000Z", "2023-01-01T23:59:00.000Z")

	ts := timeStart
	for {
		t := ts
		sendUbtTraceToIPP(t, 10)
		ts = ts.Add(60 * time.Second)
		if ts.After(timeEnd) {
			break
		}
	}

	fmt.Println("Done: ", time.Since(st))
	fmt.Println("--------------------")
}

func sendUbtTraceToIPP(t time.Time, count int) {
	for i := 0; i < 70000; i++ {
		traceData := lib.GenRawData(t)
		operationName := "local-test"
		_, span := otel.Tracer(traceData.ServiceName).Start(context.Background(), operationName,
			go_trace.WithTimestamp(traceData.Timestamp.Add(time.Duration(rand.Intn(60))*time.Second)))

		span.SetAttributes(attribute.KeyValue{
			Key:   attribute.Key("x-customer-id"),
			Value: attribute.StringValue(traceData.CustomerID),
		})
		span.SetAttributes(attribute.KeyValue{
			Key:   attribute.Key("x-cpid"),
			Value: attribute.StringValue(traceData.Cpid),
		})
		span.SetAttributes(attribute.KeyValue{
			Key:   attribute.Key("http.url"),
			Value: attribute.StringValue(traceData.API),
		})
		span.SetAttributes(attribute.KeyValue{
			Key:   attribute.Key("http.method"),
			Value: attribute.StringValue(traceData.HTTPMethod),
		})
		span.SetAttributes(attribute.KeyValue{
			Key:   attribute.Key("http.status_code"),
			Value: attribute.StringValue(traceData.HTTPStatusCode),
		})

		span.End(go_trace.WithTimestamp(traceData.Timestamp.Add(time.Duration(traceData.ResponseTime) * time.Millisecond)))
	}
	time.Sleep(1 * time.Second)
}
