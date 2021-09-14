# API Package

Uses Go-Ethereum API library.

All API calls prefixed by `hubble_`

## TO DO

* Generate API documentation / specification
  * See PR for OpenRPC
  * See JSON RPC notion file  HubbleDocs/JsonRpc
  * See Postman collection

* Dev deployments (see postman)
  * AWS - Crypto (Used for benchmarking)
    * Check if it is a bare-metal machine
  * Kubernetes
    * Primary (building batches) (Single proof of authority master)
    * Secondary (only syncing)
