package p2p

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	pb "github.com/dymensionxyz/dymint/p2p/pb"

	"github.com/libp2p/go-msgio/protoio"
)

var TraceBufferSize = 1 << 16 // 64K ought to be enough for everyone; famous last words.
var MinTraceBatchSize = 16

// rejection reasons
const (
	RejectBlacklstedPeer      = "blacklisted peer"
	RejectBlacklistedSource   = "blacklisted source"
	RejectMissingSignature    = "missing signature"
	RejectUnexpectedSignature = "unexpected signature"
	RejectUnexpectedAuthInfo  = "unexpected auth info"
	RejectInvalidSignature    = "invalid signature"
	RejectValidationQueueFull = "validation queue full"
	RejectValidationThrottled = "validation throttled"
	RejectValidationFailed    = "validation failed"
	RejectValidationIgnored   = "validation ignored"
	RejectSelfOrigin          = "self originated message"
)

type basicTracer struct {
	ch     chan struct{}
	mx     sync.Mutex
	buf    []*pb.TraceEvent
	lossy  bool
	closed bool
}

func (t *basicTracer) Trace(evt *pb.TraceEvent) {
	t.mx.Lock()
	defer t.mx.Unlock()

	if t.closed {
		return
	}

	if t.lossy && len(t.buf) > TraceBufferSize {
		//log.Debug("trace buffer overflow; dropping trace event")
	} else {
		t.buf = append(t.buf, evt)
	}

	select {
	case t.ch <- struct{}{}:
	default:
	}
}

func (t *basicTracer) Close() {
	t.mx.Lock()
	defer t.mx.Unlock()
	if !t.closed {
		t.closed = true
		close(t.ch)
	}
}

// JSONTracer is a tracer that writes events to a file, encoded in ndjson.
type JSONTracer struct {
	basicTracer
	w io.WriteCloser
}

// NewJsonTracer creates a new JSONTracer writing traces to file.
func NewJSONTracer(file string) (*JSONTracer, error) {
	return OpenJSONTracer(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
}

// OpenJSONTracer creates a new JSONTracer, with explicit control of OpenFile flags and permissions.
func OpenJSONTracer(file string, flags int, perm os.FileMode) (*JSONTracer, error) {
	f, err := os.OpenFile(file, flags, perm)
	if err != nil {
		return nil, err
	}

	tr := &JSONTracer{w: f, basicTracer: basicTracer{ch: make(chan struct{}, 1)}}
	go tr.doWrite()

	return tr, nil
}

func (t *JSONTracer) doWrite() {
	var buf []*pb.TraceEvent
	enc := json.NewEncoder(t.w)
	for {
		_, ok := <-t.ch

		t.mx.Lock()
		tmp := t.buf
		t.buf = buf[:0]
		buf = tmp
		t.mx.Unlock()

		for i, evt := range buf {
			err := enc.Encode(evt)
			if err != nil {
				//log.Warnf("error writing event trace: %s", err.Error())
			}
			buf[i] = nil
		}

		if !ok {
			t.w.Close()
			return
		}
	}
}

var _ EventTracer = (*JSONTracer)(nil)

// PBTracer is a tracer that writes events to a file, as delimited protobufs.
type PBTracer struct {
	basicTracer
	w io.WriteCloser
}

func NewPBTracer(file string) (*PBTracer, error) {
	return OpenPBTracer(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
}

// OpenPBTracer creates a new PBTracer, with explicit control of OpenFile flags and permissions.
func OpenPBTracer(file string, flags int, perm os.FileMode) (*PBTracer, error) {
	f, err := os.OpenFile(file, flags, perm)
	if err != nil {
		return nil, err
	}

	tr := &PBTracer{w: f, basicTracer: basicTracer{ch: make(chan struct{}, 1)}}
	go tr.doWrite()

	return tr, nil
}

func (t *PBTracer) doWrite() {
	var buf []*pb.TraceEvent
	w := protoio.NewDelimitedWriter(t.w)
	for {
		_, ok := <-t.ch

		t.mx.Lock()
		tmp := t.buf
		t.buf = buf[:0]
		buf = tmp
		t.mx.Unlock()

		for i, evt := range buf {
			err := w.WriteMsg(evt)
			if err != nil {
				//log.Warnf("error writing event trace: %s", err.Error())
			}
			buf[i] = nil
		}

		if !ok {
			t.w.Close()
			return
		}
	}
}

var _ EventTracer = (*PBTracer)(nil)
