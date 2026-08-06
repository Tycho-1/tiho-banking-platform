# product-catalog

Read-only Go microservice for bank product offerings (savings, CDs, investments).

## Endpoints

| Path | Auth | Description |
|------|------|-------------|
| `GET /health` | No | Liveness |
| `GET /ready` | No | Readiness |
| `GET /version` | No | Build version |
| `GET /api/products` | JWT | List products |
| `GET /api/products?type=savings\|cd\|investment` | JWT | Filter by type |

## Local run

```bash
cd src/products
go test ./...
go vet ./...

# Demo JWT public key from repo (same as kind-local demo-jwt component)
grep 'jwtRS256.key.pub:' ../../deploy/components/demo-jwt/jwt-secret.yaml \
  | awk '{print $2}' | base64 -d > /tmp/jwt.pub

export PUB_KEY_PATH=/tmp/jwt.pub
go run .
```

If your cluster is up (non-ESO overlay with demo-jwt):

```bash
kubectl get secret jwt-key -n banking-platform \
  -o jsonpath='{.data.jwtRS256.key\.pub}' | base64 -d > /tmp/jwt.pub
```

```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/products
```

## Product data

Edit `data/products.json` to add or update demo products. Rebuild the image after changes.

## CI

`.github/workflows/product-catalog.yml` runs `go test`, `go vet`, and a Docker build on
changes under `src/products/**`. Publishing to GHCR requires the repository variable
`ENABLE_IMAGE_PUSH=true`; the image tag is read from `release.yaml`.
