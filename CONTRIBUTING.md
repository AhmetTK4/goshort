# Contributing

Use Go 1.24.5 or newer and Node.js 22 with npm. Read the root README for local API, Redis, and frontend setup.

```sh
go test ./...
go vet ./...
cd frontend
npm ci
npm test -- --watchAll=false
npm run build
```

API tests use an in-memory Redis-compatible test server, so no running Redis or Docker is needed for tests. A real Redis server is required to run the API normally.

For bugs, include reproduction steps and expected/actual behavior. Remove tokens, credentials, and private URLs from logs. Discuss substantial changes in an issue first; a PR should explain the change, its scope, and the commands you ran.

This is a learning project. Do not claim production readiness or add unmeasured performance claims. Maintenance is best-effort. Report sensitive security findings privately to ahmettemelkundupoglu@gmail.com.
