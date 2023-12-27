package main

import (
	"compress/gzip"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	pb "github.com/libp2p/go-libp2p-pubsub/pb"
	//pb "github.com/dymensionxyz/dymint/p2p/pb"
	"github.com/libp2p/go-libp2p/core/peer"

	//"github.com/libp2p/go-libp2p/core/peer"

	ggio "github.com/gogo/protobuf/io"
)

// tracestat is a program that parses a pubsub tracer dump and calculates
// statistics. By default, stats are printed to stdout, but they can be
// optionally written to a JSON file for further processing.
func main() {
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	summary := flag.Bool("summary", true, "print trace summary")
	cdf := flag.Bool("cdf", false, "print propagation delay CDF")
	avg := flag.Bool("avg", false, "pring average propagation delay per message")
	flag.Parse()

	stat := &tracestat{
		peers:    make(map[peer.ID]*msgstat),
		msgs:     make(map[uint64][]int64),
		delays:   make(map[uint64][]int64),
		avgDelay: make(map[uint64]int),
	}

	for _, f := range flag.Args() {
		err = load(f, stat.addEvent)
		if err != nil {
			return
		}
	}

	if *cdf || *avg /*||dup || *hops || *jsonOut != "" */ {
		stat.compute()
	}

	if *summary {
		stat.printSummary()
	}
	if *cdf {
		stat.printCDF()
	}
	if *avg {
		stat.printAverageDelay()
	}

}

// tracestat is the tree that's populated as we parse the dump.
type tracestat struct {
	// peers summarizes per-peer stats.
	peers map[peer.ID]*msgstat

	// aggregate stats.
	aggregate msgstat

	// msgs contains per-message propagation traces: timestamps from published
	// and delivered messages.
	msgs map[uint64][]int64

	// delays is the propagation delay per message, in millis.
	delays map[uint64][]int64

	// delayCDF is the computed propagation delay distribution across all
	// messages.
	delayCDF []blocks

	avgDelay map[uint64]int

	msgsOrder []uint64
}

// msgstat holds message statistics.
type msgstat struct {
	publish int
	deliver int
}

// sample represents a CDF bucket.
type blocks struct {
	delay int
	count int
}

func load(f string, addEvent func(*pb.TraceEvent)) error {
	r, err := os.Open(f)
	if err != nil {
		return fmt.Errorf("error opening trace file %s: %w", f, err)
	}
	defer r.Close()

	gzipR, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("error opening gzip reader for %s: %w", f, err)
	}
	defer gzipR.Close()

	var evt pb.TraceEvent
	pbr := ggio.NewDelimitedReader(gzipR, 1<<20)

	for {
		evt.Reset()

		switch err = pbr.ReadMsg(&evt); err {
		case nil:
			addEvent(&evt)
		case io.EOF:
			return nil
		default:
			return fmt.Errorf("error decoding trace event from %s: %w", f, err)
		}
	}
}

func (ts *tracestat) addEvent(evt *pb.TraceEvent) {

	var (
		peer      = peer.ID(evt.GetPeerID())
		timestamp = evt.GetTimestamp()
	)

	ps, ok := ts.peers[peer]
	if !ok {
		ps = &msgstat{}
		ts.peers[peer] = ps
	}

	//fmt.Printf("message peer %s\n", peer)
	switch evt.GetType() {
	case pb.TraceEvent_PUBLISH_MESSAGE:
		b := evt.GetPublishMessage().GetMessageID()
		height := uint64(binary.LittleEndian.Uint64(b))

		ps.publish++
		ts.aggregate.publish++

		_, ok := ts.msgs[height]

		if !ok {
			//ts.count++
			//ts.msgsOrder[ts.count] = mid
			ts.msgsOrder = append(ts.msgsOrder, height)
			//fmt.Println("new message", mid)
		}
		fmt.Println("Publish block", height, timestamp)
		ts.msgs[height] = append(ts.msgs[height], timestamp)

	case pb.TraceEvent_DELIVER_MESSAGE:
		b := evt.GetDeliverMessage().GetMessageID()
		height := uint64(binary.LittleEndian.Uint64(b))
		ps.deliver++
		ts.aggregate.deliver++

		//peer := string(evt.GetDuplicateMessage().GetReceivedFrom())

		_, ok := ts.msgs[height]

		if !ok {
			//ts.count++
			//ts.msgsOrder[ts.count] = mid
			ts.msgsOrder = append(ts.msgsOrder, height)
			//fmt.Println("new deliver message", mid, ts.count, timestamp)
		}
		fmt.Println("Received block", height, timestamp)

		ts.msgs[height] = append(ts.msgs[height], timestamp)

		//ts.msgsPeer[key{peer, mid}] = append(ts.msgsPeer[key{peer, mid}], timestamp)

	}
}

func (ts *tracestat) compute() {
	// sort the message publish/delivery timestamps and transform to delays
	//fmt.Println("Computing CDF")
	for mid, timestamps := range ts.msgs {
		sort.Slice(timestamps, func(i, j int) bool {
			return timestamps[i] < timestamps[j]
		})

		delays := make([]int64, len(timestamps)-1)
		t0 := timestamps[0]
		for i, t := range timestamps[1:] {
			delays[i] = t - t0
		}
		fmt.Println("mid", mid)
		ts.delays[mid] = delays
	}

	// compute the CDF rounded to millisecond precision
	samples := make(map[int]int)
	for _, delays := range ts.delays {
		for _, dt := range delays {
			mdt := int((dt + 499999) / 1000000)
			samples[mdt]++
		}
	}

	xsamples := make([]blocks, 0, len(samples))
	for dt, count := range samples {
		xsamples = append(xsamples, blocks{dt, count})
	}
	sort.Slice(xsamples, func(i, j int) bool {
		return xsamples[i].delay < xsamples[j].delay
	})
	for i := 1; i < len(xsamples); i++ {
		xsamples[i].count += xsamples[i-1].count
	}
	ts.delayCDF = xsamples

	for mid, timestamps := range ts.msgs {
		sum := 0
		count := 0
		sort.Slice(timestamps, func(i, j int) bool {
			return timestamps[i] < timestamps[j]
		})

		for b, delay := range ts.delays[mid] {
			miliDelay := int((delay + 499999) / 1000000)
			fmt.Println("Block", b, miliDelay)
			sum += miliDelay
			if miliDelay > 0 {
				count += 1
			}
		}
		//	avgDelay := sum / count
		avgDelay := 0
		if count > 0 {
			avgDelay = sum / count
		}
		ts.avgDelay[mid] = avgDelay
	}
}

func (ts *tracestat) printSummary() {
	fmt.Printf("=== Trace Summary ===\n")
	fmt.Printf("Peers: %d\n", len(ts.peers))
	fmt.Printf("Published Blocks: %d\n", ts.aggregate.publish)
	fmt.Printf("Received Blocks: %d\n", ts.aggregate.deliver)
}

func (ts *tracestat) printCDF() {
	fmt.Printf("=== Propagation Delay CDF (ms) ===\n")
	for _, sample := range ts.delayCDF {
		fmt.Printf("%d %d\n", sample.delay, sample.count)
	}
}

func (ts *tracestat) printAverageDelay() {
	fmt.Printf("=== Average Delay (ms) ===\n")
	for _, sample := range ts.msgsOrder {
		fmt.Printf("%d %d\n", int(sample), ts.avgDelay[sample])
	}
}
