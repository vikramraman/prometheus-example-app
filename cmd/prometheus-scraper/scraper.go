package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/prometheus/common/log"
	"io"
	"io/ioutil"
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
	t := time.NewTicker(5 * time.Second)
	for _ = range t.C {
		scrape()
	}
}

func scrape() {
	client := http.Client{}
	resp, err := client.Get("http://localhost:8443/metrics")

	if err != nil {
		fmt.Errorf("%v", err)
	}

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("incorrect status: %d", resp.StatusCode)
	}

	body, err := readResponse(resp.Body, resp.ContentLength)
	if err != nil {
		fmt.Errorf("incorrect read: %v", err)
	}

	// reduce intermediate storage during parsing
	// reduce prometheus format memory during parsing
	// reduce flush buffer size - string management

	promMetrics := chunkedParsing(body)
	//promMetrics := parse(resp.Body)
	points := buildPoints(promMetrics)
	report(points)
}

func readResponse(body io.ReadCloser, cLen int64) ([]byte, error) {
	if cLen > 0 {
		buf := bytes.NewBuffer(make([]byte, 0, cLen))
		_, err := buf.ReadFrom(body)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	return ioutil.ReadAll(body)
}

func chunkedParsing(body []byte) map[string]*dto.MetricFamily {
	// split on lines
	// loopp to build metric until lookahead sees another comment
	// parse chunk
	// add results to results map
	//TODO: chunks
	// allocate buffer of size X



	return nil
}

func readChunks(body io.ReadCloser) ([]byte, error) {
	scanner := bufio.NewScanner(body)
	scanner.Split(bufio.ScanLines)
	var buf bytes.Buffer

	// read 1000 lines at a time or till EOF
	// after reading 1000 lines delegate to parseMetrics(buffer)
	// make sure to reuse the same buffer
	line := 0
	lookAhead := scanner.Scan()
	for scanner.Scan() {
		if line < 1000 {
			log.Infof("text: " + scanner.Text())
			buf.Write(scanner.Bytes())
			line++
			continue
		}
		log.Infof("text: " + scanner.Text())
		buf.Write(scanner.Bytes())
		//buf.WriteString(scanner.Text())
	}
	fmt.Println(buf.String())
}

func readMetrics(scanner bufio.Scanner) {
	buf := make([]byte, 0)
	begin := true
	for scanner.Scan() {
		//TODO: check if we need to append EOL
		line := scanner.Bytes()
		if !begin && isComment(line) {
			break
		}
		begin = false
		buf = append(buf, line...)
	}
}

func isComment() bool {

}

//func readLinesTill(n int, scanner *bufio.Scanner) {
//	var buf bytes.Buffer
//	for line := 0; line < n; line++ {
//		if scanner.Scan() {
//			buf.Write(scanner.Bytes())
//		}
//	}
//	scanner.
//}

func parse(buf []byte) map[string]*dto.MetricFamily {
	var parser expfmt.TextParser

	// parse even if the buffer begins with a newline
	buf = bytes.TrimPrefix(buf, []byte("\n"))
	// Read raw data
	buffer := bytes.NewBuffer(buf)
	reader := bufio.NewReader(buffer)

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
