package telemetry

import (
	"context"
	"time"

	"github.com/diogo464/telemetry"
	pb "github.com/libp2p/go-libp2p-kad-dht/pb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
)

var (
	AttributeMessagePutValue     = attribute.String("message_type", "PutValue")
	AttributeMessageGetValue     = attribute.String("message_type", "GetValue")
	AttributeMessageAddProvider  = attribute.String("message_type", "AddProvider")
	AttributeMessageGetProviders = attribute.String("message_type", "GetProviders")
	AttributeMessageFindNode     = attribute.String("message_type", "FindNode")
	AttributeMessagePing         = attribute.String("message_type", "Ping")
	AttributeMessageUnknown      = attribute.String("message_type", "Unknown")
)

type handlerEvent struct {
	MessageType string `json:"message_type"`
	Write       uint64 `json:"write"`
	Handler     uint64 `json:"handler"`
}

type Metrics struct {
	messageIn    syncint64.Counter
	messageOut   syncint64.Counter
	eventHandler telemetry.EventEmitter
}

func NewMetrics(provider metric.MeterProvider) (*Metrics, error) {
	m := provider.Meter("libp2p.io/dht/kad")

	MessageIn, err := m.SyncInt64().Counter(
		"message_in",
		instrument.WithDescription("Number of received messages"),
		instrument.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		return nil, err
	}

	MessageOut, err := m.SyncInt64().Counter(
		"message_out",
		instrument.WithDescription("Number of sent messages"),
		instrument.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		return nil, err
	}

	tprovider := telemetry.DowncastMeterProvider(provider)
	m2 := tprovider.TelemetryMeter("libp2p.io/telemetry")

	eventHandler := m2.Event(
		"kad.handler_timings",
		instrument.WithDescription("Message handler timings"),
	)

	return &Metrics{
		messageIn:    MessageIn,
		messageOut:   MessageOut,
		eventHandler: eventHandler,
	}, nil
}

func (m *Metrics) MessageIn(ctx context.Context, t pb.Message_MessageType) {
	m.messageIn.Add(ctx, 1, messageTypeAttribute(t))
}

func (m *Metrics) MessageOut(ctx context.Context, t pb.Message_MessageType) {
	m.messageOut.Add(ctx, 1, messageTypeAttribute(t))
}

func (m *Metrics) Handler(ctx context.Context, t pb.Message_MessageType, write time.Duration, handler time.Duration) {
	m.eventHandler.Emit(&handlerEvent{
		MessageType: t.String(),
		Write:       uint64(write.Nanoseconds()),
		Handler:     uint64(handler.Nanoseconds()),
	})
}

func messageTypeAttribute(mt pb.Message_MessageType) attribute.KeyValue {
	switch mt {
	case pb.Message_PUT_VALUE:
		return AttributeMessagePutValue
	case pb.Message_GET_VALUE:
		return AttributeMessageGetValue
	case pb.Message_ADD_PROVIDER:
		return AttributeMessageAddProvider
	case pb.Message_GET_PROVIDERS:
		return AttributeMessageGetProviders
	case pb.Message_FIND_NODE:
		return AttributeMessageFindNode
	case pb.Message_PING:
		return AttributeMessagePing
	default:
		return AttributeMessageUnknown
	}
}
