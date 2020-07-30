FROM scratch

ADD prometheus-example-app /bin/prometheus-example-app

ADD ssl/prom-example.key /ssl/prom-example.key
ADD ssl/prom-example.pem /ssl/prom-example.pem

ADD kube_state_metrics /bin/kube_state_metrics

ENTRYPOINT ["/bin/prometheus-example-app"]
