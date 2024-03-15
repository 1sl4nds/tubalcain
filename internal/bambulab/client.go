package bambulab

import (
	"context"
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Client struct {
	Device     DeviceInfo
	Host       string
	Port       int
	LocalMQTT  bool
	Username   string
	AuthToken  string
	AccessCode string
	UsageHours float64
	Client     mqtt.Client
	Connected  bool
	Tracer     trace.Tracer
}

func (c *Client) Connect(ctx context.Context) {
	opts := mqtt.NewClientOptions().AddBroker(c.Host).SetClientID("someID").SetUsername(c.Username).SetPassword(c.AuthToken)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		c.onMessage(ctx, client, msg)
	})

	c.Client = mqtt.NewClient(opts)
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Error connecting to mqtt broker: %v", token.Error())
	}

	log.Println("Connected to mqtt broker")
}

func (c *Client) onMessage(ctx context.Context, client mqtt.Client, msg mqtt.Message) {
	ctx, span := c.Tracer.Start(ctx, "onMessage", trace.WithAttributes(attribute.String("message", string(msg.Payload()))))
	defer span.End()

	// Add message as baggage
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier{})
	log.Printf("Received message on topic: %s Message: %s\n", msg.Topic(), msg.Payload())
}

func (c *Client) Publish(ctx context.Context, topic string, payload interface{}) {
	bytesPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error marshalling payload: %v", err)
	}

	ctx, span := c.Tracer.Start(ctx, "Publish", trace.WithAttributes(attribute.String("payload", string(bytesPayload))))
	defer span.End()

	// Add payload as baggage
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier{})

	token := c.Client.Publish(topic, 0, false, bytesPayload)
	token.Wait()
}

func (c *Client) Subscribe(ctx context.Context, topic string) {
	ctx, span := c.Tracer.Start(ctx, "Subscribe", trace.WithAttributes(attribute.String("topic", topic)))
	defer span.End()

	// Add topic as baggage
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier{})

	if token := c.Client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatalf("Subscribe error: %v", token.Error())
	}
	log.Printf("Subscribed to topic: %s\n", topic)
}

func (c *Client) Disconnect(ctx context.Context) {
	_, span := c.Tracer.Start(ctx, "Disconnect")
	defer span.End()

	c.Client.Disconnect(250)
	log.Println("Disconnected from mqtt broker")
}
