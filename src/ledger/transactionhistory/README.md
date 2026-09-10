# Transaction History Service

The transaction history service provides an efficient readable cache of past transactions, as read from the `ledger-db`.

The `ledger-db` service holds the source of truth for the system.
The `transaction-history`  reads and caches data from the `ledger-db`, but may be out of date when under heavy load.

Implemented in Java with Spring Boot and Guava.

### Endpoints

| Endpoint                    | Type  | Auth? | Description                                                                   |
| --------------------------- | ----- | ----- | ----------------------------------------------------------------------------- |
| `/healthy`                  | GET   |       |  Liveness probe endpoint. Monitors health of background thread.               |
| `/ready`                    | GET   |       |  Readiness probe endpoint.                                                    |
| `/transactions/<accountid>` | GET   | 🔒    |  Return the account transaction list iff authenticated to access the account. |
| `/version`                  | GET   |       |  Returns the contents of `$VERSION`                                           |
| `/actuator/prometheus`      | GET   |       |  Prometheus scrape (Micrometer). Independent of Stackdriver/`ENABLE_METRICS`. |

### Environment Variables

- `VERSION`
  - a version string for the service
- `PORT`
  - the port for the webserver
- `POLL_MS`
  - the number of milliseconds to wait in between polls to `ledger-db`
  - optional. Defaults to 100
- `CACHE_SIZE`
  - the max number of history lists to store in the cache
  - optional. Defaults to 1,000
- `CACHE_MINUTES`
  - the expiry time for the cache in minutes
  - optional. Defaults to 60
- `HISTORY_LIMIT`
  - the number of past transactions to store for each user
  - optional. Defaults to 100
- `JVM_OPTS`
  - settings for the JVM. Used to obey container memory limits
- `EXTRA_LATENCY_MILLIS`
  - add fake extra latency in milliseconds to transaction history requests
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

- [deployment/transactionhistory](../../../deploy/base/transaction-history.yaml)
- [service/transactionhistory](../../../deploy/base/transaction-history.yaml)
- Prometheus scrape: [ServiceMonitor](../../../deploy/components/prometheus-servicemonitors/servicemonitor-transactionhistory.yaml) (`/actuator/prometheus`). Works on any cluster with Prometheus Operator. **`gke-dev`** does not include that component (Cloud Ops instead).
