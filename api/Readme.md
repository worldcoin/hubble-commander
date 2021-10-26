## API Error code ranges

- `999` - Unknown Errors
- `10XXX` - Transaction Errors
- `20XXX` - Commitment Errors
- `30XXX` - Batch Errors
- `40XXX` - Badger Errors
- `50XXX` - Proof Errors
- `99XXX` - Uncategorized Errors like NetworkInfo, BLS, UserStates etc.

## API Errors

|  Error code  |                                Message                              |
| ------------ | ------------------------------------------------------------------- |
| `999`        | `unknown error: COMMANDER_ERROR`                                    |
| `10000`      | `transaction not found`                                             |
| `10001`      | `transactions not found`                                            |
| `10002`      | `some field is missing, verify the transfer/create2transfer object` |
| `10003`      | `invalid recipient, cannot send funds to yourself`                  |
| `10004`      | `nonce too low`                                                     |
| `10005`      | `nonce too high`                                                    |
| `10006`      | `not enough balance`                                                |
| `10007`      | `amount must be greater than 0`                                     |
| `10008`      | `fee too low`                                                       |
| `10009`      | `invalid signature`                                                 |
| `10010`      | `amount is not encodable as multi-precission decimal`               |
| `10011`      | `fee is not encodable as multi-precission decimal`                  |
| `20000`      | `commitment not found`                                              |
| `30000`      | `batch not found`                                                   |
| `30001`      | `batches not found`                                                 |
| `40000`      | `an error occurred while saving data to the Badger database`        |
| `40001`      | `an error occurred while iterating over Badger database`            |
| `50000`      | `proof endpoints disabled`                                          |
| `50001`      | `commitment proof not found`                                        |
| `50002`      | `public key proof not found`                                        |
| `50003`      | `user state proof not found`                                        |
| `99000`      | `an error occurred while fetching the account count`                |
| `99001`      | `public key not found`                                              |
| `99002`      | `user state not found`                                              |
| `99003`      | `user states not found`                                             |
| `99004`      | `an error occurred while fetching the domain for signing`           |
