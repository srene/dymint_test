package p2p

import (
	"encoding/binary"
	"time"

	"github.com/dymensionxyz/dymint/types"
	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/libp2p/go-libp2p/core/peer"
	//pb "github.com/dymensionxyz/dymint/p2p/pb"
)

// EventTracer is a generic event tracer interface.
// This is a high level tracing interface which delivers tracing events, as defined by the protobuf
// schema in pb/trace.proto.
type EventTracer interface {
	Trace(evt *pb.TraceEvent)
}

// RawTracer is a low level tracing interface that allows an application to trace the internal
// operation of the subsystem.
//
// Note that the tracers are invoked synchronously, which means that application tracers must
// take care to not block or modify arguments.
//
// Warning: this interface is not fixed, we may be adding new methods as necessitated by the system
// in the future.
type RawTracer interface {
	// AddPeer is invoked when a new peer is added.
	PublishBlock(p peer.ID, block *types.Block)
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
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(block.Header.Height))

	evt := &pb.TraceEvent{
		Type:           pb.TraceEvent_PUBLISH_MESSAGE.Enum(),
		PeerID:         []byte(t.pid),
		Timestamp:      &now,
		PublishMessage: &pb.TraceEvent_PublishMessage{MessageID: b},
		/*PbMessage: &pb.TraceEvent_PublishedBlock{
			Height: &block.Header.Height,
		},*/
	}

	t.tracer.Trace(evt)
	//fmt.Println("Publishing block ", block.Header.Height, p)
}
func (t *blockTracer) ReceiveBlock(p peer.ID, block *types.Block) {
	if t == nil {
		return
	}

	if t.tracer == nil {
		return
	}
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(block.Header.Height))

	now := time.Now().UnixNano()
	evt := &pb.TraceEvent{
		Type:           pb.TraceEvent_DELIVER_MESSAGE.Enum(),
		PeerID:         []byte(t.pid),
		Timestamp:      &now,
		DeliverMessage: &pb.TraceEvent_DeliverMessage{MessageID: b},
		/*RbMessage: &pb.TraceEvent_DeliverMessage{
			Height: &block.Header.Height,
		},*/
	}

	t.tracer.Trace(evt)
	//fmt.Println("Receiving block ", block.Header.Height, p)
}
