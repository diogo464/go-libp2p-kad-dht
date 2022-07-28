module github.com/libp2p/go-libp2p-kad-dht

go 1.16

require (
	github.com/diogo464/ipfs_telemetry v0.0.0-00010101000000-000000000000
	github.com/gogo/protobuf v1.3.2
	github.com/google/gopacket v1.1.19
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/golang-lru v0.5.4
	github.com/ipfs/go-cid v0.2.0
	github.com/ipfs/go-datastore v0.5.1
	github.com/ipfs/go-detect-race v0.0.1
	github.com/ipfs/go-ipfs-util v0.0.2
	github.com/ipfs/go-ipns v0.1.2
	github.com/ipfs/go-log v1.0.5
	github.com/jbenet/goprocess v0.1.4
	github.com/libp2p/go-eventbus v0.2.1
	github.com/libp2p/go-libp2p v0.20.3
	github.com/libp2p/go-libp2p-core v0.16.1
	github.com/libp2p/go-libp2p-kbucket v0.4.7
	github.com/libp2p/go-libp2p-peerstore v0.6.0
	github.com/libp2p/go-libp2p-record v0.1.3
	github.com/libp2p/go-libp2p-routing-helpers v0.2.3
	github.com/libp2p/go-libp2p-swarm v0.11.0
	github.com/libp2p/go-libp2p-testing v0.9.2
	github.com/libp2p/go-libp2p-xor v0.1.0
	github.com/libp2p/go-msgio v0.2.0
	github.com/libp2p/go-netroute v0.2.0
	github.com/multiformats/go-base32 v0.0.4
	github.com/multiformats/go-multiaddr v0.5.0
	github.com/multiformats/go-multibase v0.1.0
	github.com/multiformats/go-multihash v0.2.0
	github.com/multiformats/go-multistream v0.3.3
	github.com/stretchr/testify v1.7.1
	github.com/whyrusleeping/go-keyspace v0.0.0-20160322163242-5b898ac5add1
	go.opencensus.io v0.23.0
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/trace v1.7.0
	go.uber.org/zap v1.21.0
)

replace github.com/libp2p/go-libp2p-kbucket => ../go-libp2p-kbucket/

replace github.com/diogo464/ipfs_telemetry => ../../
