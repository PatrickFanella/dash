package metrics

import (
	"encoding/json"
	"math"
	"testing"
)

func makePromResponse(resultType string, results ...string) *promResponse {
	raw := make([]json.RawMessage, len(results))
	for i, r := range results {
		raw[i] = json.RawMessage(r)
	}
	return &promResponse{
		Status: "success",
		Data:   promData{ResultType: resultType, Result: raw},
	}
}

func TestNormalizeRangeSingleSeries(t *testing.T) {
	resp := makePromResponse("matrix",
		`{"metric":{},"values":[[1700000000,"45.5"],[1700000300,"50.2"],[1700000600,"48.1"]]}`,
	)

	ts := NormalizeRange(resp)
	if len(ts.Timestamps) != 3 {
		t.Fatalf("expected 3 points, got %d", len(ts.Timestamps))
	}
	if ts.Timestamps[0] != 1700000000000 {
		t.Errorf("expected timestamp in ms, got %d", ts.Timestamps[0])
	}
	if ts.Values[0] != 45.5 {
		t.Errorf("expected 45.5, got %f", ts.Values[0])
	}
}

func TestNormalizeRangeMultiSeriesAveraged(t *testing.T) {
	resp := makePromResponse("matrix",
		`{"metric":{"cpu":"0"},"values":[[1700000000,"40"],[1700000300,"60"]]}`,
		`{"metric":{"cpu":"1"},"values":[[1700000000,"60"],[1700000300,"80"]]}`,
	)

	ts := NormalizeRange(resp)
	if len(ts.Values) != 2 {
		t.Fatalf("expected 2 averaged points, got %d", len(ts.Values))
	}
	if ts.Values[0] != 50 {
		t.Errorf("expected average 50, got %f", ts.Values[0])
	}
	if ts.Values[1] != 70 {
		t.Errorf("expected average 70, got %f", ts.Values[1])
	}
}

func TestNormalizeRangeNaNFiltered(t *testing.T) {
	resp := makePromResponse("matrix",
		`{"metric":{},"values":[[1700000000,"45"],[1700000300,"NaN"],[1700000600,"50"]]}`,
	)

	ts := NormalizeRange(resp)
	if len(ts.Values) != 2 {
		t.Fatalf("expected 2 points (NaN filtered), got %d", len(ts.Values))
	}
}

func TestNormalizeRangeEmpty(t *testing.T) {
	resp := makePromResponse("matrix")
	ts := NormalizeRange(resp)
	if ts.Timestamps == nil || len(ts.Timestamps) != 0 {
		t.Errorf("expected empty slice, got %v", ts.Timestamps)
	}
}

func TestNormalizeRangeNilResponse(t *testing.T) {
	ts := NormalizeRange(nil)
	if len(ts.Timestamps) != 0 {
		t.Errorf("expected empty for nil, got %v", ts.Timestamps)
	}
}

func TestNormalizeInstant(t *testing.T) {
	resp := makePromResponse("vector",
		`{"metric":{},"value":[1700000000,"85.5"]}`,
	)

	iv := NormalizeInstant(resp)
	if iv.Value != 85.5 {
		t.Errorf("expected 85.5, got %f", iv.Value)
	}
	if iv.Timestamp != 1700000000000 {
		t.Errorf("expected timestamp in ms, got %d", iv.Timestamp)
	}
}

func TestNormalizeInstantNaN(t *testing.T) {
	resp := makePromResponse("vector",
		`{"metric":{},"value":[1700000000,"NaN"]}`,
	)

	iv := NormalizeInstant(resp)
	if iv.Value != 0 && !math.IsNaN(iv.Value) {
		t.Errorf("expected zero for NaN, got %f", iv.Value)
	}
}

func TestNormalizeInstantEmpty(t *testing.T) {
	resp := makePromResponse("vector")
	iv := NormalizeInstant(resp)
	if iv.Value != 0 {
		t.Errorf("expected zero for empty, got %f", iv.Value)
	}
}
