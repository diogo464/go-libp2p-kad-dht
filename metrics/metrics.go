package metrics

import (
	"context"
	"time"

	"github.com/diogo464/telemetry"
	pb "github.com/libp2p/go-libp2p-kad-dht/pb"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
)

// TODO: InstanceID should be added to Metrics.Attributes

var (
	Scope = instrumentation.Scope{
		Name:    "libp2p.io/dht/kad",
		Version: "0.18.0",
	}

	KeyMessageType = attribute.Key("message_type")
	KeyPeerID      = attribute.Key("peer_id")
	KeyInstanceID  = attribute.Key("instance_id")

	AttributeMessagePutValue     = KeyMessageType.String("PutValue")
	AttributeMessageGetValue     = KeyMessageType.String("GetValue")
	AttributeMessageAddProvider  = KeyMessageType.String("AddProvider")
	AttributeMessageGetProviders = KeyMessageType.String("GetProviders")
	AttributeMessageFindNode     = KeyMessageType.String("FindNode")
	AttributeMessagePing         = KeyMessageType.String("Ping")
	AttributeMessageUnknown      = KeyMessageType.String("Unknown")

	ViewBytesDistribution = sdkmetric.NewView(
		sdkmetric.Instrument{
			Kind:  sdkmetric.InstrumentKindSyncHistogram,
			Unit:  unit.Bytes,
			Scope: Scope,
		}, sdkmetric.Stream{
			Aggregation: aggregation.ExplicitBucketHistogram{
				Boundaries: []float64{1024, 2048, 4096, 16384, 65536, 262144, 1048576, 4194304, 16777216, 67108864, 268435456, 1073741824, 4294967296},
				NoMinMax:   false,
			},
		},
	)

	ViewMillisecondsDistribution = sdkmetric.NewView(
		sdkmetric.Instrument{
			Kind:  sdkmetric.InstrumentKindSyncHistogram,
			Unit:  unit.Milliseconds,
			Scope: Scope,
		}, sdkmetric.Stream{
			Aggregation: aggregation.ExplicitBucketHistogram{
				Boundaries: []float64{0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200, 250, 300, 400, 500, 650, 800, 1000, 2000, 5000, 10000, 20000, 50000, 100000},
				NoMinMax:   false,
			},
		},
	)

	Views = []sdkmetric.View{
		ViewBytesDistribution,
		ViewMillisecondsDistribution,
	}
)

type handlerEvent struct {
	MessageType string `json:"message_type"`
	Write       uint64 `json:"write"`
	Handler     uint64 `json:"handler"`
}

type Metrics struct {
	ReceivedMessages       syncint64.Counter
	ReceivedMessageErrors  syncint64.Counter
	ReceivedBytes          syncint64.Counter
	InboundRequestLatency  syncfloat64.Histogram
	OutboundRequestLatency syncfloat64.Histogram
	SentMessages           syncint64.Counter
	SentMessageErrors      syncint64.Counter
	SentRequests           syncint64.Counter
	SentRequestErrors      syncint64.Counter
	SentBytes              syncint64.Counter

	attributes   []attribute.KeyValue
	eventHandler telemetry.EventEmitter
}

func New(provider metric.MeterProvider) (*Metrics, error) {
	m := provider.Meter(Scope.Name, metric.WithInstrumentationVersion(Scope.Version), metric.WithSchemaURL(Scope.SchemaURL))

	ReceivedMessages, err := m.SyncInt64().Counter(
		"received_messages",
		instrument.WithDescription("Total number of messages received per RPC"),
		instrument.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		return nil, err
	}

	ReceivedMessageErrors, err := m.SyncInt64().Counter(
		"received_message_errors",
		instrument.WithDescription("Total number of errors for messages received per RPC"),
		instrument.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		return nil, err
	}

	ReceivedBytes, err := m.SyncInt64().Counter(
		"received_bytes",
		instrument.WithDescription("Total received bytes per RPC"),
		instrument.WithUnit(unit.Bytes),
	)
	if err != nil {
		return nil, err
	}

	InboundRequestLatency, err := m.SyncFloat64().Histogram(
		"inbound_request_latency",
		instrument.WithDescription("Latency per RPC"),
		instrument.WithUnit(unit.Milliseconds),
	)
	if err != nil {
		return nil, err
	}

	OutboundRequestLatency, err := m.SyncFloat64().Histogram(
		"outbound_request_latency",
		instrument.WithDescription("Latency per RPC"),
		instrument.WithUnit(unit.Milliseconds),
	)
	if err != nil {
		return nil, err
	}

	SentMessages, err := m.SyncInt64().Counter(
		"sent_messages",
		instrument.WithDescription("Total number of messages sent per RPC"),
		instrument.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		return nil, err
	}

	SentMessageErrors, err := m.SyncInt64().Counter(
		"sent_message_errors",
		instrument.WithDescription("Total number of errors for messages sent per RPC"),
		instrument.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		return nil, err
	}

	SentRequests, err := m.SyncInt64().Counter(
		"sent_requests",
		instrument.WithDescription("Total number of requests sent per RPC"),
		instrument.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		return nil, err
	}

	SentRequestErrors, err := m.SyncInt64().Counter(
		"sent_request_errors",
		instrument.WithDescription("Total number of errors for requests sent per RPC"),
		instrument.WithUnit(unit.Dimensionless),
	)
	if err != nil {
		return nil, err
	}

	SentBytes, err := m.SyncInt64().Counter(
		"sent_bytes",
		instrument.WithDescription("Total sent bytes per RPC"),
		instrument.WithUnit(unit.Bytes),
	)
	if err != nil {
		return nil, err
	}

	tprovider := telemetry.DowncastMeterProvider(provider)
	m2 := tprovider.TelemetryMeter(Scope.Name, metric.WithInstrumentationVersion(Scope.Version), metric.WithSchemaURL(Scope.SchemaURL))

	eventHandler := m2.Event(
		"handler_timings",
		instrument.WithDescription("Message handler timings"),
	)

	return &Metrics{
		ReceivedMessages:       ReceivedMessages,
		ReceivedMessageErrors:  ReceivedMessageErrors,
		ReceivedBytes:          ReceivedBytes,
		InboundRequestLatency:  InboundRequestLatency,
		OutboundRequestLatency: OutboundRequestLatency,
		SentMessages:           SentMessages,
		SentMessageErrors:      SentMessageErrors,
		SentRequests:           SentRequests,
		SentRequestErrors:      SentRequestErrors,
		SentBytes:              SentBytes,
		eventHandler:           eventHandler,
	}, nil
}

func (m *Metrics) Attributes(attrs ...attribute.KeyValue) []attribute.KeyValue {
	if len(m.attributes) == 0 {
		return attrs
	}
	joined := make([]attribute.KeyValue, 0, len(m.attributes)+len(attrs))
	joined = append(joined, m.attributes...)
	joined = append(joined, attrs...)
	return joined
}

func (m *Metrics) Handler(ctx context.Context, t pb.Message_MessageType, write time.Duration, handler time.Duration) {
	m.eventHandler.Emit(&handlerEvent{
		MessageType: t.String(),
		Write:       uint64(write.Nanoseconds()),
		Handler:     uint64(handler.Nanoseconds()),
	})
}

func MessageTypeAttribute(mt pb.Message_MessageType) attribute.KeyValue {
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
