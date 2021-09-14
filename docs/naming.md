# Naming

## Tests

We use two styles of writing tests: **individual tests** and **test suites** using `testify` library.

### Individual tests

Standalone functions: `func TestFunctionName_OptionalDescription(t *testing.T)`

Methods on structs: `func TestStructName_MethodName_OptionalDescription(t *testing.T)`

### Test suites

Test suite name is derived from the Go file name or struct name. Therefore we don't repeat the struct name in method test names.

Standalone functions: 

`func (s *FileNameTestSuite) TestFunctionName_OptionalDescription()`

Methods on structs: 

`func (s *FileNameTestSuite) TestMethodName_OptionalDescription()`

## Fields

There are different names used for the same things across commander and contracts codebases. We will unify naming as part of a future PR choosing the names in green.

Account Tree

- **Value:** Public Key
- **Key:** Account Index = **PubKey ID** = Account ID

State Tree

- **Value:** User State (PubKey ID, token id, nonce, balance)
- **Key:** State Index = **State ID** = Leaf Merkle Path

Transaction:

- **From** and **to** fields are state IDs