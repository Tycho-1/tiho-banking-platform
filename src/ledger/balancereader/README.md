# Balance Reader Service

The balance reader service provides an efficient readable cache of user balances, as read from the `ledger-db`.

The `ledger-db` service holds the source of truth for the system.
The `balance-reader` reads and caches data from the `ledger-db`, but may be out of date when under heavy load.

Implemented in Java with Spring Boot and Guava.

### Endpoints

| Endpoint                | Type  | Auth? | Description                                                             |
| ----------------------- | ----- | ----- | ----------------------------------------------------------------------- |
| `/balances/<accountid>` | GET   | 🔒    |  Get the account balance iff owned by the currently authenticated user. |
| `/healthy`              | GET   |       |  Liveness probe endpoint. Monitors health of background thread.         |
| `/ready`                | GET   |       |  Readiness probe endpoint.                                              |
| `/version`              | GET   |       |  Returns the contents of `$VERSION`                                     |
| `/actuator/prometheus`  | GET   |       |  Prometheus scrape (Micrometer). Independent of Stackdriver/`ENABLE_METRICS`. |

### Environment Variables

- `VERSION`
  - a version string for the service
- `PORT`
  - the port for the webserver
- `POLL_MS`
  - the number of milliseconds to wait in between polls to `ledger-db`
  - optional. Defaults to 100
- `CACHE_SIZE`
  - the max number of account balances to store in the cache
  - optional. Defaults to 1,000,000
- `JVM_OPTS`
  - settings for the JVM. Used to obey container memory limits
- `LOG_LEVEL`
  - service level [log level](https://logging.apache.org/log4j/2.x/manual/customloglevels.html)
- `ENABLE_METRICS`
  - `true` to **push** Micrometer metrics to GCP Cloud Monitoring (Stackdriver). Non-GKE overlays (`disable-gcp-telemetry`) set `false`.
  - Prometheus scrape at `/actuator/prometheus` is always on (pull); it does not use this flag.

- ConfigMap `environment-config`:
  - `LOCAL_ROUTING_NUM`
    - the routing number for our bank
  - `PUB_KEY_PATH`
    - the path to the JWT signer's public key, mounted as a secret

- ConfigMap `ledger-db-config`:
  - `SPRING_DATASOURCE_URL`
    - URL of the `ledger-db` service
  - `SPRING_DATASOURCE_USERNAME`
    - username for the `ledger-db` database
  - `SPRING_DATASOURCE_PASSWORD`
    - password for the `ledger-db` database

### Kubernetes Resources

- [deployment/balancereader](../../../deploy/base/balance-reader.yaml)
- [service/balancereader](../../../deploy/base/balance-reader.yaml)
- Prometheus scrape: [ServiceMonitor](../../../deploy/components/prometheus-servicemonitors/servicemonitor-balancereader.yaml) (`/actuator/prometheus`). Works on any cluster with Prometheus Operator. **`gke-dev`** does not include that component (Cloud Ops instead).
