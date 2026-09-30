# Ledger Writer Service

The ledger writer service accepts and validates incoming transactions before writing them to the ledger.

Implemented in Java with Spring Boot.

### Endpoints

| Endpoint           | Type  | Auth? | Description                                          |
| ------------------ | ----- | ----- | ---------------------------------------------------- |
| `/ready`           | GET   |       |  Readiness probe endpoint.                           |
| `/transactions`    | POST  | 🔒    |  Submits a transaction to be appended to the ledger. |
| `/version`         | GET   |       |  Returns the contents of `$VERSION`                  |
| `/actuator/prometheus` | GET |    |  Prometheus scrape (Micrometer). Independent of Stackdriver/`ENABLE_METRICS`. |

### Environment Variables

- `VERSION`
  - a version string for the service
- `PORT`
  - the port for the webserver
- `JVM_OPTS`
  - settings for the JVM. Used to obey container memory limits
- `LOG_LEVEL`
  - the service-wide [log level](https://logging.apache.org/log4j/2.x/manual/customloglevels.html) (default: INFO)
- `ENABLE_METRICS`
  - `true` to **push** Micrometer metrics to GCP Cloud Monitoring (Stackdriver). Non-GKE overlays (`disable-gcp-telemetry`) set `false`.
  - Prometheus scrape at `/actuator/prometheus` is always on (pull); it does not use this flag.
  
- ConfigMap `environment-config`:
  - `LOCAL_ROUTING_NUM`
    - the routing number for our bank
  - `PUB_KEY_PATH`
    - the path to the JWT signer's public key, mounted as a secret

- ConfigMap `service-api-config`
  - `BALANCES_API_ADDR`
    - the address and port of the `balancereader` service

- ConfigMap `ledger-db-config`:
  - `SPRING_DATASOURCE_URL`
    - URL of the `ledger-db` service
  - `SPRING_DATASOURCE_USERNAME`
    - username for the `ledger-db` database
  - `SPRING_DATASOURCE_PASSWORD`
    - password for the `ledger-db` database

### Kubernetes Resources

- [deployment/ledgerwriter](../../../deploy/base/ledger-writer.yaml)
- [service/ledgerwriter](../../../deploy/base/ledger-writer.yaml)
- Prometheus scrape: [ServiceMonitor](../../../deploy/components/prometheus-servicemonitors/servicemonitor-ledgerwriter.yaml) (`/actuator/prometheus`). Works on any cluster with Prometheus Operator. **`gke-dev`** does not include that component (Cloud Ops instead).
