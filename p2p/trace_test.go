package p2p

import (
	"context"
	"crypto/rand"
	"os"
	"testing"
	"time"

	"github.com/dymensionxyz/dymint/config"
	"github.com/dymensionxyz/dymint/types"
	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-msgio/protoio"
)

type traceStats struct {
	publish, received int
}

func (t *traceStats) process(test *testing.T, evt *pb.TraceEvent) {
	// fmt.Printf("process event %s\n", evt.GetType())
	switch evt.GetType() {
	case pb.TraceEvent_PUBLISH_MESSAGE:
		t.publish++
	case pb.TraceEvent_DELIVER_MESSAGE:
		t.received++
	}
}

func (ts *traceStats) check(t *testing.T) {
	if ts.publish == 0 {
		t.Fatal("expected non-zero count")
	}
	if ts.received == 0 {
		t.Fatal("expected non-zero count")
	}
}

func TestPBTracer(t *testing.T) {
	tracer, err := NewPBTracer("/tmp/trace.out.pb")
	if err != nil {
		t.Fatal(err)
	}

	opts := []Option{
		WithEventTracer(tracer),
	}
	privKey, _, _ := crypto.GenerateEd25519Key(rand.Reader)
	p2pClient, err := NewClient(config.P2PConfig{}, privKey, "TestChain", log.TestingLogger(), opts...)
	p2pClient.Start(context.Background())
	header := &types.Header{Height: 1}
	block := &types.Block{Header: *header}
	t.Log(block.Header.Height)
	t.Log(block.Header.Height)

	p2pClient.BlockPublished(block)
	p2pClient.BlockReceived(block)
	//testWithTracer(t, tracer)
	time.Sleep(time.Second)
	tracer.Close()

	var stats traceStats
	var evt pb.TraceEvent

	f, err := os.Open("/tmp/trace.out.pb")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	r := protoio.NewDelimitedReader(f, 1<<20)
	for {
		evt.Reset()
		err := r.ReadMsg(&evt)
		if err != nil {
			break
		}

		stats.process(t, &evt)
	}

	stats.check(t)
}
