package p2p

import (
	"fmt"
	"time"

	"github.com/dymensionxyz/dymint/types"
	"github.com/libp2p/go-libp2p/core/peer"

	pb "github.com/dymensionxyz/dymint/p2p/pb"
)

// EventTracer is a generic event tracer interface.
// This is a high level tracing interface which delivers tracing events, as defined by the protobuf
// schema in pb/trace.proto.
type EventTracer interface {
	Trace(evt *pb.TraceEvent)
}

// RawTracer is a low level tracing interface that allows an application to trace the internal
// operation of the pubsub subsystem.
//
// Note that the tracers are invoked synchronously, which means that application tracers must
// take care to not block or modify arguments.
//
// Warning: this interface is not fixed, we may be adding new methods as necessitated by the system
// in the future.
type RawTracer interface {
	// AddPeer is invoked when a new peer is added.
	ReceiveBlock(p peer.ID, Block types.Block)
}

// pubsub tracer details
type blockTracer struct {
	tracer EventTracer
	raw    []RawTracer
	pid    peer.ID
}

func (t *blockTracer) PublishBlock(p peer.ID, block *types.Block) {
	if t == nil {
		return
	}

	if t.tracer == nil {
		return
	}
	now := time.Now().UnixNano()
	evt := &pb.TraceEvent{
		Type:      pb.TraceEvent_PUBLISHED_BLOCk.Enum(),
		PeerID:    []byte(t.pid),
		Timestamp: &now,
		PbMessage: &pb.TraceEvent_PublishedBlock{
			Height: &block.Header.Height,
		},
	}

	t.tracer.Trace(evt)
	fmt.Println("Publishing block ", block.Header.Height, p)
}
func (t *blockTracer) ReceiveBlock(p peer.ID, block *types.Block) {
	if t == nil {
		return
	}

	if t.tracer == nil {
		return
	}
	now := time.Now().UnixNano()
	evt := &pb.TraceEvent{
		Type:      pb.TraceEvent_RECEIVED_BLOCK.Enum(),
		PeerID:    []byte(t.pid),
		Timestamp: &now,
		RbMessage: &pb.TraceEvent_RececeivedBlock{
			Height: &block.Header.Height,
		},
	}

	t.tracer.Trace(evt)
	fmt.Println("Receiving block ", block.Header.Height, p)
}
