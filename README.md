# tubalcain

## TL;DR

- A MQTT client implementation for IoT devices.
- Integrates with OpenTelemetry for observability.
- Docker and Docker Compose support for easy setup and deployment.
- Configurations for Prometheus, Grafana, OpenSearch, and more included.

## Getting Started

### Prerequisites

- [Install Docker](https://docs.docker.com/engine/install/)
- [Install Go](https://go.dev/doc/install)

### Installation

1. Clone the repository:

```sh
git clone https://github.com/organization/printfarm.git
cd printfarm
```

2. Build the Docker image:

```sh
docker compose build
```

## Usage

To run the application along with its dependencies (Prometheus, Grafana, OpenTelemetry Collector, etc.), use Docker Compose:

```sh
docker compose up
```

## Contributing

Please read [CONTRIBUTING.md](https://github.com/organization/go-bambulab/blob/main/CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Deployment

The `docker-compose.yml` file included in the project root defines the deployment for local development and testing. For production environments, ensure you configure the environment variables securely and consider scaling the services as necessary.

## Built-with

- Go - The programming language used.
- Eclipse Paho MQTT Go Client - For MQTT communication.
- OpenTelemetry - For tracing and observability.
- Docker & Docker Compose - For containerization and orchestration.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/organization/go-bambulab/tags).

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## References

- [Eclipse Paho MQTT Go Client documentation](https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang)
- [OpenTelemetry Go SDK](https://github.com/open-telemetry/opentelemetry-go)
- [Docker documentation](https://docs.docker.com/)
- [Docker Compose documentation](https://docs.docker.com/compose/)
