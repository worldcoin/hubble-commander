# BLS Signature library

Implements the [BLS digital signature](https://en.wikipedia.org/wiki/BLS_digital_signature) scheme. This signature scheme supports aggregation.

Uses the Ethereum compatible BN254 pairing curve (aka alt-BN128; see [EIP-196](https://eips.ethereum.org/EIPS/eip-196)). We may prefer BLS12-381 or another pairing curve once Ethereum gains support for it (see [EIP-2539](https://eips.ethereum.org/EIPS/eip-2539)).

The BN254 pairing curve and BLS implementation are from <https://github.com/kilic/bn254/blob/master/bls/bls.go>.
