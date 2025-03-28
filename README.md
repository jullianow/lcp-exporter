# LCP Exporter

LCP Exporter is a tool that exports metrics from the Liferay Cloud Platform (LCP) to Prometheus.
This exporter collects data related to cluster infrastructure and status, such as runtime environment information and API availability.

## Prerequisites

- Go 1.18 or higher.
- Prometheus to collect the exported metrics.
- Go dependencies, such as `prometheus/client_golang`.

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/jullianow/lcp-exporter.git
cd lcp-exporter
```

### 2. Install dependencies

If you don't have Go installed yet, you can download it [here](https://golang.org/dl/).

```bash
go mod tidy
```

### 3. Build the project

To build the `lcp-exporter`, run the following command:

```bash
go build -o lcp-exporter ./cmd
```

### 4. Run the project

After building, you can run the exporter like this:

```bash
./lcp-exporter
```

This will start the exporter and expose metrics in the Prometheus-compatible format.

## Configuration

The exporter can be configured to connect to different LCP API endpoints. Make sure to provide the appropriate credentials and configuration for the LCP API you are using.

Example configuration for specifying an endpoint:

```bash
./lcp-exporter --api-url=http://api-endpoint-url
```

## Tests

To run the unit tests, use the following command:

```bash
go test ./...
```

## Monitoring with Prometheus

Once the `lcp-exporter` is running, Prometheus can begin collecting the exposed metrics. Add the following scrape target in the Prometheus configuration file:

```yaml
scrape_configs:
  - job_name: 'lcp-exporter'
    static_configs:
      - targets: ['localhost:9103']
```

After that, Prometheus will start collecting metrics from the `lcp-exporter` at the configured URL.

## Contributing

Contributions are welcome! If you find any issues or would like to add improvements, feel free to open an issue or submit a pull request.

1. Fork this repository.
2. Create a branch (`git checkout -b my-feature`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin my-feature`).
6. Open a pull request.
