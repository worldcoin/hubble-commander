import os
import json
import shutil
from pathlib import Path

source_path = Path('../hubble-contracts/artifacts')
output_path = Path("../contracts")

def generate_bindings(path, type, pkg, filename):
    with open(path, 'r') as f:
        data = f.read()
        obj = json.loads(data)
        with open('tmp.abi', 'w') as out:
            out.write(json.dumps(obj['abi']))
        with open('tmp.bin', 'w') as out:
            out.write(obj['bytecode'])
    os.system(
        f'abigen --abi tmp.abi --bin tmp.bin --pkg {pkg} --type {type} --out {filename}')


def generate(artifact, name):
    path = source_path / artifact
    lower = name.lower()
    os.makedirs(output_path / lower, exist_ok=True)
    filename = (output_path / lower / lower).with_suffix('.go')
    generate_bindings(path, name, lower, filename)


def generate_subdir(subdir, artifact, name):
    path = source_path / artifact
    prefix_len = len(subdir)
    pkg = name[prefix_len:].lower()
    os.makedirs(output_path / subdir / pkg, exist_ok=True)
    filename = (output_path / subdir / pkg / pkg).with_suffix('.go')
    generate_bindings(path, name, pkg, filename)


os.system('rm -rf contracts')

generate('contracts/test/POB.sol/ProofOfBurn.json', 'ProofOfBurn')
generate('contracts/proposers/Chooser.sol/Chooser.json', 'Chooser')
generate('contracts/TokenRegistry.sol/TokenRegistry.json', 'TokenRegistry')
generate('contracts/SpokeRegistry.sol/SpokeRegistry.json', 'SpokeRegistry')
generate('contracts/Vault.sol/Vault.json', 'Vault')
generate('contracts/WithdrawManager.sol/WithdrawManager.json', 'WithdrawManager')
generate('contracts/DepositManager.sol/DepositManager.json', 'DepositManager')
generate('contracts/BLSAccountRegistry.sol/BLSAccountRegistry.json', 'AccountRegistry')
generate('contracts/Transfer.sol/Transfer.json', 'Transfer')
generate('contracts/MassMigrations.sol/MassMigration.json', 'MassMigration')
generate('contracts/Create2Transfer.sol/Create2Transfer.json', 'Create2Transfer')
generate('contracts/rollup/Rollup.sol/Rollup.json', 'Rollup')

generate_subdir('frontend', 'contracts/client/FrontendGeneric.sol/FrontendGeneric.json', 'FrontendGeneric')
generate_subdir('frontend', 'contracts/client/FrontendTransfer.sol/FrontendTransfer.json', 'FrontendTransfer')
generate_subdir('frontend', 'contracts/client/FrontendMassMigration.sol/FrontendMassMigration.json', 'FrontendMassMigration')
generate_subdir('frontend', 'contracts/client/FrontendCreate2Transfer.sol/FrontendCreate2Transfer.json', 'FrontendCreate2Transfer')

generate_subdir('test', 'contracts/test/TestTx.sol/TestTx.json', 'TestTx')
generate_subdir('test', 'contracts/test/TestTypes.sol/TestTypes.json', 'TestTypes')
generate_subdir('test', 'contracts/test/TestBLS.sol/TestBLS.json', 'TestBLS')
generate_subdir('test', 'contracts/test/CustomToken.sol/CustomToken.json', 'TestCustomToken')

generate_subdir('libs', 'contracts/libs/BNPairingPrecompileCostEstimator.sol/BNPairingPrecompileCostEstimator.json', 'CostEstimator')

generate('@openzeppelin/contracts/token/ERC20/IERC20.sol/IERC20.json', 'ERC20')

os.remove('tmp.abi')
os.remove('tmp.bin')
