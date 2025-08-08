package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	wslogger "github.com/thiagozs/go-wslogger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func initTracer(ctx context.Context) (*sdktrace.TracerProvider, error) {
	// >>> Exporter via OTLP HTTP para o Collector em localhost:4318
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:4318"), // sem "http://"
		otlptracehttp.WithInsecure(),                 // HTTP plain
		// otlphttp.WithURLPath("/v1/traces"),   // opcional (default já é /v1/traces)
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("myAppServices"),
			semconv.DeploymentEnvironment("dev"),
			attribute.String("owner", "thiago"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		),
	)
	return tp, nil
}

func main() {
	ctx := context.Background()
	tp, err := initTracer(ctx)
	if err != nil {
		log.Fatalf("init tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("shutdown tracer: %v", err)
		}
	}()

	// ====== Logger ======
	lg := wslogger.NewLogger(
	//  wslogger.WithKind(wslogger.Stdout),
	//      wslogger.WithFormat(wslogger.JSON),
	//      wslogger.WithAppName("poc-wslogger"),
	)

	tr := otel.Tracer("demo")

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, span := tr.Start(ctx, "work-handler")
		defer span.End()

		// Simula algo
		sleep := time.Duration(50+rand.IntN(250)) * time.Millisecond
		time.Sleep(sleep)
		span.SetAttributes(attribute.Int64("work.sleep_ms", sleep.Milliseconds()))

		// IDs do trace/span
		sc := oteltrace.SpanContextFromContext(ctx)
		traceID := sc.TraceID().String()
		spanID := sc.SpanID().String()

		// Log com trace_id/span_id
		lg.Info("processando requisição",
			"trace_id", traceID,
			"span_id", spanID,
			"sleep_ms", sleep.Milliseconds(),
			"path", r.URL.Path,
		)

		fmt.Fprintf(w, "ok\n")
	})

	addr := ":8080"
	lg.Info("HTTP up", "addr", addr, "pid", os.Getpid())
	if err := http.ListenAndServe(addr, nil); err != nil {
		lg.Error("server error", "err", err)
	}
}
