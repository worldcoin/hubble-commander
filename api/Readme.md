## API Error code ranges

- `999` - Unknown Errors
- `10XXX` - Transaction Errors
- `20XXX` - Commitment Errors
- `30XXX` - Batch Errors
- `40XXX` - Badger Errors
- `50XXX` - Proof Errors
- `99XXX` - Uncategorized Errors like NetworkInfo, BLS, UserStates etc.
- `-32XXX` - JSON-RPC library errors

## Commander API errors

| Error code | Message                                                                                                   |
|------------|-----------------------------------------------------------------------------------------------------------|
| `999`      | `unknown error: COMMANDER_ERROR`                                                                          |
| `10000`    | `transaction not found`                                                                                   |
| `10001`    | _DEPRECATED_:`transactions not found`                                                                     |
| `10002`    | `some field is missing, verify the transfer/create2transfer object`                                       |
| `10003`    | `invalid recipient, cannot send funds to yourself`                                                        |
| `10004`    | `nonce too low`                                                                                           |
| `10005`    | _DEPRECATED_: `nonce too high`                                                                            |
| `10006`    | `not enough balance`                                                                                      |
| `10007`    | `amount must be greater than 0`                                                                           |
| `10008`    | `fee too low`                                                                                             |
| `10009`    | `invalid signature`                                                                                       |
| `10010`    | `amount is not encodable as multi-precission decimal`                                                     |
| `10011`    | `fee is not encodable as multi-precission decimal`                                                        |
| `10012`    | `sender state ID does not exist`                                                                          |
| `10013`    | _DEPRECATED_:`spoke ID must be greater than 0`                                                            |
| `10014`    | `cannot update mined transaction`                                                                         |
| `10015`    | `transaction already exists`                                                                              |
| `10016`    | `spoke with given ID does not exist`                                                                      |
| `10017`    | `commander instance is not accepting transactions`                                                                  |
| `20000`    | `commitment not found`                                                                                    |
| `30000`    | `batch not found`                                                                                         |
| `30001`    | `batches not found`                                                                                       |
| `40000`    | `an error occurred while saving data to the Badger database`                                              |
| `40001`    | `an error occurred while iterating over Badger database`                                                  |
| `50000`    | `proof methods disabled`                                                                                  |
| `50001`    | `commitment inclusion proof could not be generated`                                                       |
| `50002`    | `public key inclusion proof could not be generated`                                                       |
| `50003`    | `user state inclusion proof could not be generated`                                                       |
| `50004`    | `mass migration commitment inclusion proof could not be generated`                                        |
| `50005`    | `withdraw proof could not be calculated for a given batch`                                                |
| `50006`    | `invalid batch type, only mass migration batches are supported`                                           |
| `50007`    | `mass migration with given transaction hash was not found in a given commitment`                          |
| `50008`    | `commitment inclusion proof can only be generated for Transfer/Create2Transfer commitments`               |
| `50009`    | `mass migration commitment inclusion proof cannot be generated for different type of commitments`         |
| `99000`    | `an error occurred while fetching the account count`                                                      |
| `99001`    | `public key not found`                                                                                    |
| `99002`    | `user state not found`                                                                                    |
| `99003`    | `user states not found`                                                                                   |
| `99004`    | `an error occurred while fetching the domain for signing`                                                 |

## JSON-RPC library errors

| Error code           | Type               | Example messages                                       |
|----------------------|--------------------|--------------------------------------------------------|
| `-32700`             | `Parse error`      | `"parse error"`                                        |
| `-32600`             | `Invalid request`  | `"empty batch"` or `"invalid request"`                 |
| `-32601`             | `Method not found` | `"the method foo_bar does not exist/is not available"` |
| `-32602`             | `Invalid params`   | `"invalid argument 0: error unmarshalling Uint256"`    |
| `-32603`             | `Internal error`   | *reserved error code (probably unused)*                |
| `-32099` to `-32000` | `Internal error`   | *reserved error code range (probably unused)*          |

Read more about default errors in JSON-RPC 2.0 specification [here](https://www.jsonrpc.org/specification).
