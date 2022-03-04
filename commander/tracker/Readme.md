# Tracker

### Tracker.TrackTxs
* Receives transaction requests and sends them in order to avoid incorrect nonce.
* Tracks the status of sent transactions and returns an error if the transaction failed.
