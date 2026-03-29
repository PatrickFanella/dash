package metrics

// QueryType indicates whether a metric uses instant or range queries.
type QueryType int

const (
	RangeQuery   QueryType = iota
	InstantQuery
)

// MetricDef defines a named PromQL query and its type.
type MetricDef struct {
	Query string
	Type  QueryType
}

// Queries maps metric names to their PromQL definitions.
var Queries = map[string]MetricDef{
	// CPU usage as percentage (0-100). Averages across all cores.
	"cpu": {
		Query: `100 - (avg(rate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)`,
		Type:  RangeQuery,
	},
	// Memory usage as percentage (0-100).
	"memory": {
		Query: `(1 - node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * 100`,
		Type:  RangeQuery,
	},
	// Network receive throughput in bytes/sec.
	"network_rx": {
		Query: `rate(node_network_receive_bytes_total{device=~"eth.*|en.*"}[5m])`,
		Type:  RangeQuery,
	},
	// Network transmit throughput in bytes/sec.
	"network_tx": {
		Query: `rate(node_network_transmit_bytes_total{device=~"eth.*|en.*"}[5m])`,
		Type:  RangeQuery,
	},
	// Disk usage as percentage (0-100) for root filesystem.
	"disk": {
		Query: `(1 - node_filesystem_avail_bytes{mountpoint="/"} / node_filesystem_size_bytes{mountpoint="/"}) * 100`,
		Type:  InstantQuery,
	},
	// CPU temperature in Celsius.
	"temperature": {
		Query: `node_hwmon_temp_celsius`,
		Type:  RangeQuery,
	},
	// System boot time as Unix seconds. Subtract from now() to get uptime.
	"uptime": {
		Query: `node_boot_time_seconds`,
		Type:  InstantQuery,
	},
}
