# API Package

Uses Go-Ethereum API library.

All API calls prefixed by `hubble_`

TODO:

* Generate API documentation / specificiation
  * See PR for OpenRPC
  * See JSON RPC notion file  HubbleDocs/JsonRpc
  * See Postman collection

* Dev deployments (see postman)
  * AWS - Crypto (Used for benchmarking)
    * Check if it is a bare-metal machine
  * Kubernetes
    * Primary (building batches) (Single proof of authority master)
    * Secondary (only syncing)

* Desirable renames
  * ProofOfBurn -> ProofOfAuthority
  * Create2 -> Anything better
