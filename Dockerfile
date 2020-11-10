FROM scratch

ADD prometheus-example-app /bin/prometheus-example-app

ADD sample_prom_metrics /bin/sample_metrics

ENTRYPOINT ["/bin/prometheus-example-app"]
