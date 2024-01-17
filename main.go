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
	server.GET("/ubtv2", ubtV2)
	server.Run(":8888")
}

func ubt(c *gin.Context) {
	trafficPerMin, _ := c.GetQuery("traffic_per_min")
	startTime, _ := c.GetQuery("start_time")
	endTime, _ := c.GetQuery("end_time")
	cnt, _ := strconv.ParseInt(trafficPerMin, 10, 64)
	if cnt == 0 {
		cnt = 70000
	}
	if startTime == "" || endTime == "" {
		c.Error(fmt.Errorf("start_time and end_time are required"))
		return
	}
	ubtTracingData(startTime, endTime, int(cnt))
}

func ubtV2(c *gin.Context) {
	trafficPerMin, _ := c.GetQuery("traffic_per_min")
	duration, _ := c.GetQuery("duration")
	cnt, _ := strconv.ParseInt(trafficPerMin, 10, 64)
	if cnt == 0 {
		cnt = 70000
	}
	dcnt, _ := strconv.ParseInt(duration, 10, 64)
	if dcnt == 0 {
		dcnt = 3600
	}
	for i := 0; i < int(dcnt); i++ {
		sendUbtTraceCurrentTimeToIPP(time.Now(), int(cnt))
		time.Sleep(1 * time.Minute)
	}
}

func ubtTracingData(startTime, endTime string, trafficPerMin int) {
	st := time.Now()
	timeStart, _ := time.Parse("2006-01-02T15:04:05.000Z", startTime)
	timeEnd, _ := time.Parse("2006-01-02T15:04:05.000Z", endTime)

	ts := timeStart
	for {
		t := ts
		sendUbtTraceToIPP(t, trafficPerMin)
		ts = ts.Add(60 * time.Second)
		if ts.After(timeEnd) {
			break
		}
	}

	fmt.Println("Done: ", time.Since(st))
	fmt.Println("--------------------")
}

func sendUbtTraceToIPP(t time.Time, trafficPerMin int) {
	for i := 0; i < trafficPerMin; i++ {
		traceData := lib.GenRawData(t)
		operationName := "local-test"
		st := traceData.Timestamp.Add(time.Duration(rand.Intn(60)) * time.Second)
		et := st.Add(time.Duration(traceData.ResponseTime) * time.Millisecond)
		_, span := otel.Tracer(traceData.ServiceName).Start(context.Background(), operationName,
			go_trace.WithTimestamp(st),
		)

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

		span.End(go_trace.WithTimestamp(et))

		time.Sleep(10 * time.Millisecond)
	}
}

func sendUbtTraceCurrentTimeToIPP(t time.Time, trafficPerMin int) {
	for i := 0; i < trafficPerMin; i++ {
		traceData := lib.GenRawData(t)
		operationName := "local-test"
		_, span := otel.Tracer(traceData.ServiceName).Start(context.Background(), operationName)

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

		span.End()
	}
}
