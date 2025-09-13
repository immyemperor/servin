---
layout: default
title: Logging and Monitoring
---

# Logging and Monitoring

Comprehensive guide to observability, logging, and monitoring with Servin Container Runtime.

## Monitoring Architecture

### Observability Stack

Complete monitoring and observability solution:

```
Observability Architecture:
┌─────────────────────────────────────────────────────────┐
│                    Dashboards & Alerts                  │
│                   (Grafana, Kibana)                     │
├─────────────────────────────────────────────────────────┤
│                 Metrics & Log Storage                   │
│              (Prometheus, Elasticsearch)               │
├─────────────────────────────────────────────────────────┤
│               Collection & Aggregation                  │
│             (Node Exporter, Fluentd, Beats)            │
├─────────────────────────────────────────────────────────┤
│                   Servin Runtime                        │
│              (Containers, Images, Volumes)              │
└─────────────────────────────────────────────────────────┘
```

### Monitoring Components

- **Metrics Collection**: Prometheus, InfluxDB, DataDog
- **Log Aggregation**: Fluentd, Filebeat, Logstash
- **Visualization**: Grafana, Kibana, Custom dashboards
- **Alerting**: AlertManager, PagerDuty, Slack integration
- **Tracing**: Jaeger, Zipkin, OpenTelemetry

## Container Logging

### Basic Logging

Access container logs with various options:

```bash
# View container logs
servin logs nginx-container

# Follow logs in real-time
servin logs --follow nginx-container

# Show timestamps
servin logs --timestamps nginx-container

# Show last N lines
servin logs --tail 50 nginx-container

# Show logs since specific time
servin logs --since 2024-01-01T00:00:00Z nginx-container
servin logs --since 1h nginx-container

# Show logs until specific time
servin logs --until 2024-01-01T12:00:00Z nginx-container

# Filter logs by timestamp range
servin logs --since 1h --until 30m nginx-container

# Show logs for multiple containers
servin logs web-server db-server cache-server
```

### Advanced Logging Options

Configure detailed logging behavior:

```bash
# Show logs with details
servin logs --details nginx-container

# Limit log output size
servin logs --tail 100 --since 1h nginx-container

# Output logs in JSON format
servin logs --format json nginx-container

# Save logs to file
servin logs nginx-container > container-logs.txt

# Continuous log monitoring
servin logs --follow --tail 0 nginx-container | tee -a monitor.log

# Filter logs with grep
servin logs nginx-container | grep ERROR

# Show logs for all containers
servin logs $(servin ps -q)
```

### Log Drivers

Configure different logging drivers:

```bash
# JSON file driver (default)
servin run --log-driver json-file nginx:latest

# Syslog driver
servin run --log-driver syslog \
  --log-opt syslog-address=udp://logs.company.com:514 \
  nginx:latest

# Fluentd driver
servin run --log-driver fluentd \
  --log-opt fluentd-address=fluentd.company.com:24224 \
  --log-opt tag=nginx.access \
  nginx:latest

# Journald driver
servin run --log-driver journald nginx:latest

# Splunk driver
servin run --log-driver splunk \
  --log-opt splunk-token=your-token \
  --log-opt splunk-url=https://splunk.company.com:8088 \
  nginx:latest

# AWS CloudWatch driver
servin run --log-driver awslogs \
  --log-opt awslogs-group=myapp \
  --log-opt awslogs-region=us-west-2 \
  nginx:latest
```

## Log Aggregation

### Centralized Logging with Fluentd

Deploy Fluentd for log collection:

```bash
# Fluentd configuration
# fluent.conf
<source>
  @type forward
  port 24224
  bind 0.0.0.0
</source>

<filter servin.**>
  @type parser
  key_name log
  <parse>
    @type json
  </parse>
</filter>

<match servin.**>
  @type elasticsearch
  host elasticsearch
  port 9200
  index_name servin-logs
  type_name container
</match>

# Deploy Fluentd
servin run -d \
  --name fluentd \
  -p 24224:24224 \
  -v $(pwd)/fluent.conf:/fluentd/etc/fluent.conf \
  -v /var/lib/servin/containers:/var/lib/servin/containers:ro \
  fluent/fluentd:latest

# Configure containers to use Fluentd
servin run -d \
  --name web-app \
  --log-driver fluentd \
  --log-opt fluentd-address=localhost:24224 \
  --log-opt tag=webapp.access \
  nginx:latest
```

### ELK Stack Integration

Set up Elasticsearch, Logstash, and Kibana:

```bash
# Elasticsearch
servin run -d \
  --name elasticsearch \
  -p 9200:9200 \
  -p 9300:9300 \
  -e "discovery.type=single-node" \
  -e "ES_JAVA_OPTS=-Xms512m -Xmx512m" \
  -v es-data:/usr/share/elasticsearch/data \
  elasticsearch:7.15.2

# Logstash
# logstash.conf
input {
  beats {
    port => 5044
  }
}

filter {
  if [fields][container_name] {
    mutate {
      add_field => { "container" => "%{[fields][container_name]}" }
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "servin-logs-%{+YYYY.MM.dd}"
  }
}

servin run -d \
  --name logstash \
  -p 5044:5044 \
  -v $(pwd)/logstash.conf:/usr/share/logstash/pipeline/logstash.conf \
  logstash:7.15.2

# Kibana
servin run -d \
  --name kibana \
  -p 5601:5601 \
  -e ELASTICSEARCH_HOSTS=http://elasticsearch:9200 \
  kibana:7.15.2

# Filebeat for log shipping
# filebeat.yml
filebeat.inputs:
- type: container
  paths:
    - '/var/lib/servin/containers/*/*.log'
  processors:
    - add_docker_metadata:
        host: "unix:///var/run/servin.sock"

output.logstash:
  hosts: ["logstash:5044"]

servin run -d \
  --name filebeat \
  --user=root \
  -v $(pwd)/filebeat.yml:/usr/share/filebeat/filebeat.yml \
  -v /var/lib/servin/containers:/var/lib/servin/containers:ro \
  -v /var/run/servin.sock:/var/run/servin.sock:ro \
  elastic/filebeat:7.15.2
```

## Metrics Collection

### Prometheus Integration

Set up Prometheus for metrics collection:

```bash
# Prometheus configuration
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "rules/*.yml"

scrape_configs:
  - job_name: 'servin'
    static_configs:
      - targets: ['localhost:9323']
    metrics_path: '/metrics'
    scrape_interval: 5s

  - job_name: 'node-exporter'
    static_configs:
      - targets: ['node-exporter:9100']

  - job_name: 'cadvisor'
    static_configs:
      - targets: ['cadvisor:8080']

alerting:
  alertmanagers:
    - static_configs:
        - targets:
          - alertmanager:9093

# Deploy Prometheus
servin run -d \
  --name prometheus \
  -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  -v prometheus-data:/prometheus \
  prom/prometheus:latest

# Node Exporter for host metrics
servin run -d \
  --name node-exporter \
  -p 9100:9100 \
  --pid=host \
  -v "/:/host:ro,rslave" \
  prom/node-exporter:latest \
  --path.rootfs=/host

# cAdvisor for container metrics
servin run -d \
  --name cadvisor \
  -p 8080:8080 \
  --privileged \
  --device=/dev/kmsg \
  -v /:/rootfs:ro \
  -v /var/run:/var/run:ro \
  -v /sys:/sys:ro \
  -v /var/lib/servin/:/var/lib/servin:ro \
  -v /dev/disk/:/dev/disk:ro \
  gcr.io/cadvisor/cadvisor:latest
```

### Custom Metrics

Expose application metrics:

```bash
# Application with Prometheus metrics
# app.py
from prometheus_client import Counter, Histogram, generate_latest
import time

REQUEST_COUNT = Counter('app_requests_total', 'Total requests')
REQUEST_LATENCY = Histogram('app_request_duration_seconds', 'Request latency')

@REQUEST_LATENCY.time()
def process_request():
    REQUEST_COUNT.inc()
    time.sleep(0.1)
    return "OK"

# Metrics endpoint
@app.route('/metrics')
def metrics():
    return generate_latest()

# Deploy application with metrics
servin run -d \
  --name app-with-metrics \
  -p 8000:8000 \
  -p 8001:8001 \
  myapp:latest

# Scrape application metrics
# Add to prometheus.yml:
  - job_name: 'myapp'
    static_configs:
      - targets: ['app-with-metrics:8001']
```

## Visualization and Dashboards

### Grafana Setup

Deploy Grafana for visualization:

```bash
# Deploy Grafana
servin run -d \
  --name grafana \
  -p 3000:3000 \
  -e GF_SECURITY_ADMIN_PASSWORD=admin123 \
  -v grafana-data:/var/lib/grafana \
  grafana/grafana:latest

# Configure Prometheus datasource
curl -X POST http://admin:admin123@localhost:3000/api/datasources \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Prometheus",
    "type": "prometheus",
    "url": "http://prometheus:9090",
    "access": "proxy",
    "isDefault": true
  }'
```

### Container Dashboard

Create comprehensive container monitoring dashboard:

```json
{
  "dashboard": {
    "title": "Servin Container Monitoring",
    "panels": [
      {
        "title": "Container CPU Usage",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(container_cpu_usage_seconds_total[5m])",
            "legendFormat": "{{name}}"
          }
        ]
      },
      {
        "title": "Container Memory Usage",
        "type": "graph",
        "targets": [
          {
            "expr": "container_memory_usage_bytes",
            "legendFormat": "{{name}}"
          }
        ]
      },
      {
        "title": "Container Network I/O",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(container_network_receive_bytes_total[5m])",
            "legendFormat": "{{name}} RX"
          },
          {
            "expr": "rate(container_network_transmit_bytes_total[5m])",
            "legendFormat": "{{name}} TX"
          }
        ]
      },
      {
        "title": "Container Disk I/O",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(container_fs_reads_bytes_total[5m])",
            "legendFormat": "{{name}} Read"
          },
          {
            "expr": "rate(container_fs_writes_bytes_total[5m])",
            "legendFormat": "{{name}} Write"
          }
        ]
      }
    ]
  }
}
```

## Alerting

### AlertManager Configuration

Set up alert management:

```yaml
# alertmanager.yml
global:
  smtp_smarthost: 'mail.company.com:587'
  smtp_from: 'alerts@company.com'

route:
  group_by: ['alertname']
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: 'web.hook'

receivers:
  - name: 'web.hook'
    email_configs:
      - to: 'admin@company.com'
        subject: 'Alert: {{ .GroupLabels.alertname }}'
        body: |
          {{ range .Alerts }}
          Alert: {{ .Annotations.summary }}
          Description: {{ .Annotations.description }}
          {{ end }}
    
    slack_configs:
      - api_url: 'https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK'
        channel: '#alerts'
        title: 'Servin Alert'
        text: '{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}'

inhibit_rules:
  - source_match:
      severity: 'critical'
    target_match:
      severity: 'warning'
    equal: ['alertname', 'dev', 'instance']

# Deploy AlertManager
servin run -d \
  --name alertmanager \
  -p 9093:9093 \
  -v $(pwd)/alertmanager.yml:/etc/alertmanager/alertmanager.yml \
  prom/alertmanager:latest
```

### Alert Rules

Define alerting rules:

```yaml
# container-alerts.yml
groups:
  - name: container.rules
    rules:
      - alert: ContainerHighCPU
        expr: rate(container_cpu_usage_seconds_total[5m]) * 100 > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Container {{ $labels.name }} high CPU usage"
          description: "Container {{ $labels.name }} CPU usage is above 80%"

      - alert: ContainerHighMemory
        expr: container_memory_usage_bytes / container_spec_memory_limit_bytes * 100 > 90
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Container {{ $labels.name }} high memory usage"
          description: "Container {{ $labels.name }} memory usage is above 90%"

      - alert: ContainerDown
        expr: up{job="servin"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Servin daemon is down"
          description: "Servin daemon has been down for more than 1 minute"

      - alert: ContainerRestarting
        expr: increase(container_restart_count[1h]) > 5
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "Container {{ $labels.name }} restarting frequently"
          description: "Container {{ $labels.name }} has restarted {{ $value }} times in the last hour"

      - alert: ContainerVolumeUsage
        expr: container_fs_usage_bytes / container_fs_limit_bytes * 100 > 90
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Container {{ $labels.name }} volume usage high"
          description: "Container {{ $labels.name }} volume usage is above 90%"
```

## Health Checks and Monitoring

### Container Health Checks

Implement comprehensive health monitoring:

```bash
# Container with health check
servin run -d \
  --name web-app \
  --health-cmd "curl -f http://localhost:8080/health" \
  --health-interval 30s \
  --health-timeout 10s \
  --health-retries 3 \
  --health-start-period 60s \
  nginx:latest

# Custom health check script
servin run -d \
  --name app-with-health \
  --health-cmd "/app/health-check.sh" \
  --health-interval 30s \
  myapp:latest

# Monitor health status
servin inspect app-with-health --format "{{.State.Health.Status}}"

# View health check logs
servin inspect app-with-health --format "{{json .State.Health.Log}}"

# List unhealthy containers
servin ps --filter health=unhealthy
```

### Service Discovery

Implement service discovery for monitoring:

```bash
# Consul for service discovery
servin run -d \
  --name consul \
  -p 8500:8500 \
  -e CONSUL_BIND_INTERFACE=eth0 \
  consul:latest

# Register services with Consul
curl -X PUT http://localhost:8500/v1/agent/service/register \
  -d '{
    "Name": "web-app",
    "Address": "192.168.1.100",
    "Port": 8080,
    "Check": {
      "HTTP": "http://192.168.1.100:8080/health",
      "Interval": "30s"
    }
  }'

# Prometheus with Consul discovery
# prometheus.yml
scrape_configs:
  - job_name: 'consul'
    consul_sd_configs:
      - server: 'consul:8500'
    relabel_configs:
      - source_labels: [__meta_consul_service]
        target_label: job
```

## Performance Monitoring

### Resource Monitoring

Monitor system and container resources:

```bash
# Real-time container stats
servin stats

# Historical resource usage
servin run --rm \
  -v /var/run/servin.sock:/var/run/servin.sock \
  monitoring/container-stats:latest

# System resource monitoring
servin run -d \
  --name resource-monitor \
  --privileged \
  -v /proc:/host/proc:ro \
  -v /sys:/host/sys:ro \
  -v /:/rootfs:ro \
  monitoring/system-stats:latest

# Network monitoring
servin run -d \
  --name network-monitor \
  --net=host \
  --cap-add=NET_ADMIN \
  monitoring/network-stats:latest
```

### Application Performance Monitoring

Monitor application performance:

```bash
# APM with Elastic APM
servin run -d \
  --name apm-server \
  -p 8200:8200 \
  -e output.elasticsearch.hosts=elasticsearch:9200 \
  elastic/apm-server:7.15.2

# Application with APM agent
servin run -d \
  --name instrumented-app \
  -e ELASTIC_APM_SERVER_URL=http://apm-server:8200 \
  -e ELASTIC_APM_SERVICE_NAME=myapp \
  -e ELASTIC_APM_ENVIRONMENT=production \
  myapp:instrumented

# Custom performance metrics
servin run -d \
  --name perf-monitor \
  -v /var/run/servin.sock:/var/run/servin.sock \
  -v performance-data:/data \
  monitoring/performance:latest
```

## Distributed Tracing

### Jaeger Integration

Set up distributed tracing:

```bash
# Jaeger all-in-one
servin run -d \
  --name jaeger \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 14250:14250 \
  jaegertracing/all-in-one:latest

# Application with tracing
servin run -d \
  --name traced-app \
  -e JAEGER_ENDPOINT=http://jaeger:14268/api/traces \
  -e JAEGER_SERVICE_NAME=myapp \
  -e JAEGER_SAMPLER_TYPE=const \
  -e JAEGER_SAMPLER_PARAM=1 \
  myapp:traced

# OpenTelemetry collector
# otel-config.yml
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:

exporters:
  jaeger:
    endpoint: jaeger:14250
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch]
      exporters: [jaeger]

servin run -d \
  --name otel-collector \
  -p 4317:4317 \
  -p 4318:4318 \
  -v $(pwd)/otel-config.yml:/etc/otel-collector-config.yml \
  otel/opentelemetry-collector:latest \
  --config=/etc/otel-collector-config.yml
```

## Log Analytics

### Advanced Log Analysis

Implement sophisticated log analysis:

```bash
# Log analysis with ELK
# logstash-advanced.conf
input {
  beats {
    port => 5044
  }
}

filter {
  if [fields][container_name] {
    grok {
      match => { 
        "message" => "%{COMBINEDAPACHELOG}" 
      }
    }
    
    date {
      match => [ "timestamp", "dd/MMM/yyyy:HH:mm:ss Z" ]
    }
    
    mutate {
      convert => { "response" => "integer" }
      convert => { "bytes" => "integer" }
    }
    
    if [response] >= 400 {
      mutate {
        add_tag => [ "error" ]
      }
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "servin-logs-%{+YYYY.MM.dd}"
  }
}

# Real-time log analysis
servin run -d \
  --name log-analyzer \
  -v log-analysis-rules:/etc/rules \
  -e ELASTICSEARCH_URL=http://elasticsearch:9200 \
  monitoring/log-analyzer:latest

# Anomaly detection
servin run -d \
  --name anomaly-detector \
  -e ML_MODEL_PATH=/models/anomaly-model.pkl \
  -v anomaly-models:/models \
  monitoring/anomaly-detector:latest
```

## Automation and Integration

### Monitoring Automation

Automate monitoring deployment and management:

```bash
#!/bin/bash
# deploy-monitoring.sh

# Deploy monitoring stack
servin-compose -f monitoring-stack.yml up -d

# Wait for services to be ready
sleep 30

# Configure Grafana datasources
curl -X POST http://admin:admin123@localhost:3000/api/datasources \
  -H "Content-Type: application/json" \
  -d @datasource-config.json

# Import dashboards
for dashboard in dashboards/*.json; do
  curl -X POST http://admin:admin123@localhost:3000/api/dashboards/db \
    -H "Content-Type: application/json" \
    -d @"$dashboard"
done

# Setup alerts
curl -X POST http://localhost:9093/api/v1/alerts \
  -H "Content-Type: application/json" \
  -d @alert-rules.json

echo "Monitoring stack deployed successfully"
```

### Integration with CI/CD

Integrate monitoring with deployment pipelines:

```yaml
# .github/workflows/deploy-with-monitoring.yml
name: Deploy with Monitoring

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy application
        run: |
          servin run -d --name myapp myapp:${{ github.sha }}
          
      - name: Setup monitoring
        run: |
          # Add monitoring labels
          servin update myapp \
            --label monitoring.enabled=true \
            --label monitoring.service=myapp \
            --label monitoring.version=${{ github.sha }}
          
      - name: Configure health checks
        run: |
          servin update myapp \
            --health-cmd "curl -f http://localhost:8080/health" \
            --health-interval 30s
          
      - name: Register with service discovery
        run: |
          curl -X PUT http://consul:8500/v1/agent/service/register \
            -d @service-definition.json
```

This comprehensive logging and monitoring guide covers all aspects of observability with Servin, from basic logging to advanced distributed tracing and automated monitoring solutions.
