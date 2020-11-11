# Prometheus Example App

This example app instruments HTTP handlers with [Prometheus](https://prometheus.io/) metrics using the Prometheus [go client](https://github.com/prometheus/client_golang).

Metrics are exposed via the following endpoints:
* `/metrics`: Exposes various prometheus metrics (histograms, gauges, counters etc)
* `/counters`: Exposes a configurable (using --num) number of counters

## Development

Run `make container` to build a docker image.
