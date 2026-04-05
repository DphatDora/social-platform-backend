package logger

import (
	"sync"
	"time"
)

type requestBucket struct {
	sec        int64
	count      int64
	latencyNS  int64
	minLatency int64
	maxLatency int64
}

type RequestMetricsSnapshot struct {
	WindowSeconds   int     `json:"window_seconds"`
	RequestCount    int64   `json:"request_count"`
	RequestsPerMin  float64 `json:"requests_per_min"`
	AvgLatencyMS    float64 `json:"avg_latency_ms"`
	MinLatencyMS    float64 `json:"min_latency_ms"`
	MaxLatencyMS    float64 `json:"max_latency_ms"`
	Status2xxCount  int64   `json:"status_2xx_count"`
	Status3xxCount  int64   `json:"status_3xx_count"`
	Status4xxCount  int64   `json:"status_4xx_count"`
	Status5xxCount  int64   `json:"status_5xx_count"`
	LastUpdatedUnix int64   `json:"last_updated_unix"`
}

var (
	requestMetricsMu sync.Mutex
	requestBuckets   [60]requestBucket
	status2xxBuckets [60]requestBucket
	status3xxBuckets [60]requestBucket
	status4xxBuckets [60]requestBucket
	status5xxBuckets [60]requestBucket
)

func ObserveRequest(latency time.Duration, statusCode int) {
	now := time.Now().Unix()
	idx := int(now % int64(len(requestBuckets)))
	latencyNS := latency.Nanoseconds()
	if latencyNS < 0 {
		latencyNS = 0
	}

	requestMetricsMu.Lock()
	defer requestMetricsMu.Unlock()

	b := &requestBuckets[idx]
	if b.sec != now {
		*b = requestBucket{sec: now}
	}
	b.count++
	b.latencyNS += latencyNS
	if b.minLatency == 0 || latencyNS < b.minLatency {
		b.minLatency = latencyNS
	}
	if latencyNS > b.maxLatency {
		b.maxLatency = latencyNS
	}

	incStatusBucket(now, statusCode)
}

func SnapshotRequestMetrics() RequestMetricsSnapshot {
	now := time.Now().Unix()

	requestMetricsMu.Lock()
	defer requestMetricsMu.Unlock()

	var (
		totalCount   int64
		totalLatency int64
		minLatency   int64
		maxLatency   int64
		status2xx    int64
		status3xx    int64
		status4xx    int64
		status5xx    int64
	)

	for i := range requestBuckets {
		if now-requestBuckets[i].sec >= int64(len(requestBuckets)) {
			continue
		}

		totalCount += requestBuckets[i].count
		totalLatency += requestBuckets[i].latencyNS
		if requestBuckets[i].minLatency > 0 && (minLatency == 0 || requestBuckets[i].minLatency < minLatency) {
			minLatency = requestBuckets[i].minLatency
		}
		if requestBuckets[i].maxLatency > maxLatency {
			maxLatency = requestBuckets[i].maxLatency
		}

		if now-status2xxBuckets[i].sec < int64(len(status2xxBuckets)) {
			status2xx += status2xxBuckets[i].count
		}
		if now-status3xxBuckets[i].sec < int64(len(status3xxBuckets)) {
			status3xx += status3xxBuckets[i].count
		}
		if now-status4xxBuckets[i].sec < int64(len(status4xxBuckets)) {
			status4xx += status4xxBuckets[i].count
		}
		if now-status5xxBuckets[i].sec < int64(len(status5xxBuckets)) {
			status5xx += status5xxBuckets[i].count
		}
	}

	snap := RequestMetricsSnapshot{
		WindowSeconds:   len(requestBuckets),
		RequestCount:    totalCount,
		RequestsPerMin:  float64(totalCount),
		Status2xxCount:  status2xx,
		Status3xxCount:  status3xx,
		Status4xxCount:  status4xx,
		Status5xxCount:  status5xx,
		LastUpdatedUnix: now,
	}

	if totalCount > 0 {
		snap.AvgLatencyMS = float64(totalLatency) / float64(totalCount) / float64(time.Millisecond)
		snap.MinLatencyMS = float64(minLatency) / float64(time.Millisecond)
		snap.MaxLatencyMS = float64(maxLatency) / float64(time.Millisecond)
	}

	return snap
}

func incStatusBucket(now int64, statusCode int) {
	buckets := &status5xxBuckets
	switch {
	case statusCode >= 200 && statusCode < 300:
		buckets = &status2xxBuckets
	case statusCode >= 300 && statusCode < 400:
		buckets = &status3xxBuckets
	case statusCode >= 400 && statusCode < 500:
		buckets = &status4xxBuckets
	}

	idx := int(now % int64(len(requestBuckets)))
	b := &(*buckets)[idx]
	if b.sec != now {
		*b = requestBucket{sec: now}
	}
	b.count++
}
