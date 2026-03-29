package metrics

import (
	"encoding/json"
	"math"
	"strconv"
)

// NormalizeRange converts a Prometheus matrix response into a TimeSeries.
// If multiple series exist, they are averaged per timestamp.
func NormalizeRange(resp *promResponse) TimeSeries {
	if resp == nil || len(resp.Data.Result) == 0 {
		return TimeSeries{Timestamps: []int64{}, Values: []float64{}}
	}

	var allSeries []promMatrixSample
	for _, raw := range resp.Data.Result {
		var s promMatrixSample
		if err := json.Unmarshal(raw, &s); err != nil {
			continue
		}
		allSeries = append(allSeries, s)
	}

	if len(allSeries) == 0 {
		return TimeSeries{Timestamps: []int64{}, Values: []float64{}}
	}

	if len(allSeries) == 1 {
		return extractTimeSeries(allSeries[0].Values)
	}

	// Multiple series — average by timestamp
	type bucket struct {
		sum   float64
		count int
	}
	buckets := make(map[int64]*bucket)
	var orderedTS []int64

	for _, s := range allSeries {
		for _, raw := range s.Values {
			ts, val := parseSamplePair(raw)
			if math.IsNaN(val) || math.IsInf(val, 0) {
				continue
			}
			b, exists := buckets[ts]
			if !exists {
				b = &bucket{}
				buckets[ts] = b
				orderedTS = append(orderedTS, ts)
			}
			b.sum += val
			b.count++
		}
	}

	result := TimeSeries{
		Timestamps: make([]int64, 0, len(orderedTS)),
		Values:     make([]float64, 0, len(orderedTS)),
	}
	for _, ts := range orderedTS {
		b := buckets[ts]
		result.Timestamps = append(result.Timestamps, ts)
		result.Values = append(result.Values, b.sum/float64(b.count))
	}
	return result
}

// NormalizeInstant converts a Prometheus vector response into an InstantValue.
func NormalizeInstant(resp *promResponse) InstantValue {
	if resp == nil || len(resp.Data.Result) == 0 {
		return InstantValue{}
	}

	var s promVectorSample
	if err := json.Unmarshal(resp.Data.Result[0], &s); err != nil {
		return InstantValue{}
	}

	ts, val := parseSamplePair(s.Value)
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return InstantValue{}
	}
	return InstantValue{Value: val, Timestamp: ts}
}

func extractTimeSeries(values []json.RawMessage) TimeSeries {
	ts := TimeSeries{
		Timestamps: make([]int64, 0, len(values)),
		Values:     make([]float64, 0, len(values)),
	}
	for _, raw := range values {
		tsMs, val := parseSamplePair(raw)
		if math.IsNaN(val) || math.IsInf(val, 0) {
			continue
		}
		ts.Timestamps = append(ts.Timestamps, tsMs)
		ts.Values = append(ts.Values, val)
	}
	return ts
}

// parseSamplePair parses a Prometheus [timestamp, "value"] pair.
// The pair is a JSON array with a number (timestamp) and a string (value).
func parseSamplePair(raw json.RawMessage) (int64, float64) {
	var pair [2]any
	if err := json.Unmarshal(raw, &pair); err != nil {
		return 0, math.NaN()
	}

	// Timestamp is a JSON number → float64
	tsFloat, ok := pair[0].(float64)
	if !ok {
		return 0, math.NaN()
	}
	tsMs := int64(tsFloat * 1000)

	// Value is a JSON string → parse as float
	valStr, ok := pair[1].(string)
	if !ok {
		return tsMs, math.NaN()
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return tsMs, math.NaN()
	}
	return tsMs, val
}
