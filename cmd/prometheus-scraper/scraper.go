package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	dto "github.com/prometheus/client_model/go"

	"github.com/prometheus/common/expfmt"
)

// Represents a single point in Wavefront metric format.
type MetricPoint struct {
	Metric    string
	Value     float64
	Timestamp int64
	Source    string
	Tags      map[string]string

	// Embed Prometheus labels directly to avoid memory allocations
	Labels []*dto.LabelPair

	SrcTags map[string]string
}

func main() {
	t := time.NewTicker(10 * time.Second)
	for _ = range t.C {
		scrape()
	}
}

func scrape() {
	client := http.Client{}
	response, _ := client.Get("http://localhost:8443/metrics")

	promMetrics := parse(response.Body)
	points := buildPoints(promMetrics)
	report(points)

}

func parse(reader io.ReadCloser) map[string]*dto.MetricFamily {
	parser := expfmt.TextParser{}
	metrics, _ := parser.TextToMetricFamilies(reader)
	return metrics
}

func buildPoints(families map[string]*dto.MetricFamily) []*MetricPoint {
	//TODO: refine as needed
	points := make([]*MetricPoint, 0)
	for metricName, mf := range families {
		for _, m := range mf.Metric {
			pts := buildPoint(metricName, m, 0)
			if len(pts) > 0 {
				points = append(points, pts...)
			}
		}
	}
	return points
}

func buildPoint(name string, m *dto.Metric, now int64) []*MetricPoint {
	//TODO: why does this return a slice in the collector (can returning just the point make things better?)

	var result []*MetricPoint
	if m.Gauge != nil {
		if !math.IsNaN(m.GetGauge().GetValue()) {
			point := metricPoint(name+".gauge", m.GetGauge().GetValue(), now, "source", nil)
			result = filterAppend(result, point, m)
		}
	} else if m.Counter != nil {
		if !math.IsNaN(m.GetCounter().GetValue()) {
			point := metricPoint(name+".counter", m.GetCounter().GetValue(), now, "source", nil)
			result = filterAppend(result, point, m)
		}
	} else if m.Untyped != nil {
		if !math.IsNaN(m.GetUntyped().GetValue()) {
			point := metricPoint(name+".value", m.GetUntyped().GetValue(), now, "source", nil)
			result = filterAppend(result, point, m)
		}
	}
	return result
}

func filterAppend(slice []*MetricPoint, point *MetricPoint, m *dto.Metric) []*MetricPoint {
	//TODO: refine as needed
	return append(slice, point)
}

func report(points []*MetricPoint) {
	fmt.Printf("reporting %d metrics to wavefront\n", len(points))
}

func metricPoint(name string, value float64, ts int64, source string, tags map[string]string) *MetricPoint {
	return &MetricPoint{
		Metric:    "prefix." + strings.Replace(name, "_", ".", -1),
		Value:     value,
		Timestamp: ts,
		Source:    source,
		Tags:      tags,
	}
}
