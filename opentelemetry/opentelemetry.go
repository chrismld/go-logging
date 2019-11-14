package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"

	"go.opentelemetry.io/otel/exporter/trace/stdout"
	"go.opentelemetry.io/otel/global"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	file, err := os.OpenFile("info.log", os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.WarnLevel)

	initTracer()
	tracer := global.TraceProvider().GetTracer("ex.com/basic")
	tracer.WithSpan(context.Background(), "main",
		func(ctx context.Context) error {
			tracer.WithSpan(ctx, "microservice",
				func(ctx context.Context) error {
					tracer.WithSpan(ctx, "database",
						func(ctx context.Context) error {
							return nil
						},
					)
					return nil
				},
			)
			return nil
		},
	)
}

func initTracer() {
	exporter, err := stdout.NewExporter(stdout.Options{PrettyPrint: true})
	if err != nil {
		log.Fatal(err)
	}
	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter))
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)
}
