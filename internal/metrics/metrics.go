package metrics

import (
	"sync/atomic"
	"time"
)

type TunnelMetrics struct {
	Connections    uint64
	BytesSent      uint64
	BytesReceived  uint64
	LastConnection time.Time
}

var Metrics = &TunnelMetrics{}

func AddConnection() {
	atomic.AddUint64(&Metrics.Connections, 1)
	Metrics.LastConnection = time.Now()
}

func AddBytesSent(n uint64) {
	atomic.AddUint64(&Metrics.BytesSent, n)
}

func AddBytesReceived(n uint64) {
	atomic.AddUint64(&Metrics.BytesReceived, n)
}
