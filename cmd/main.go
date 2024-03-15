package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/1sl4nds/tubalcain/internal/bambulab"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() *trace.TracerProvider {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String("BambuClient"))),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp
}

func main() {
	tp := initTracer()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}()

	ctx := context.Background()
	tracer := otel.Tracer("BambuClientTracer")

	client := bambulab.Client{
		Device: bambulab.DeviceInfo{
			DeviceType: os.Getenv("DEVICE_TYPE"),
			ID:         os.Getenv("DEVICE_ID"),
		},
		Host:       os.Getenv("HOST"),
		Username:   os.Getenv("USERNAME"),
		AuthToken:  os.Getenv("AUTH_TOKEN"),
		AccessCode: os.Getenv("ACCESS_CODE"),
		Tracer:     tracer,
	}

	client.Connect(ctx)
	defer client.Disconnect(ctx)

	// Subscribe to a topic
	client.Subscribe(ctx, "device/01S09A2C0500103/report")

	// Publish a message
	// client.Publish(ctx, "device/01S09A2C0500103/report", map[string]string{"message": "Hello mqtt"})

	// Graceful shutdown on interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
