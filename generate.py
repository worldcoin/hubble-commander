import os
import json
import shutil
from pathlib import Path


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
    base = Path('hubble-contracts/artifacts')
    path = base / artifact
    lower = name.lower()
    os.makedirs(Path("contracts") / lower, exist_ok=True)
    filename = (Path('contracts') / lower / lower).with_suffix('.go')
    generate_bindings(path, name, lower, filename)


def generate_frontend(artifact, name):
    base = Path('hubble-contracts/artifacts')
    path = base / artifact
    pkg = name[8:].lower()
    os.makedirs(Path("contracts/frontend") / pkg, exist_ok=True)
    filename = (Path('contracts/frontend') / pkg / pkg).with_suffix('.go')
    generate_bindings(path, name, pkg, filename)


os.system('rm -rf contracts')

generate('contracts/proposers/POB.sol/ProofOfBurn.json', 'ProofOfBurn')
generate('contracts/TokenRegistry.sol/TokenRegistry.json', 'TokenRegistry')
generate('contracts/SpokeRegistry.sol/SpokeRegistry.json', 'SpokeRegistry')
generate('contracts/Vault.sol/Vault.json', 'Vault')
generate('contracts/DepositManager.sol/DepositManager.json', 'DepositManager')
generate('contracts/BLSAccountRegistry.sol/BLSAccountRegistry.json', 'AccountRegistry')
generate('contracts/Transfer.sol/Transfer.json', 'Transfer')
generate('contracts/MassMigrations.sol/MassMigration.json', 'MassMigration')
generate('contracts/Create2Transfer.sol/Create2Transfer.json', 'Create2Transfer')
generate('contracts/rollup/Rollup.sol/Rollup.json', 'Rollup')

generate_frontend('contracts/client/FrontendGeneric.sol/FrontendGeneric.json', 'FrontendGeneric')
generate_frontend('contracts/client/FrontendTransfer.sol/FrontendTransfer.json', 'FrontendTransfer')
generate_frontend('contracts/client/FrontendMassMigration.sol/FrontendMassMigration.json', 'FrontendMassMigration')
generate_frontend('contracts/client/FrontendCreate2Transfer.sol/FrontendCreate2Transfer.json', 'FrontendCreate2Transfer')

generate('@openzeppelin/contracts/token/ERC20/IERC20.sol/IERC20.json', 'ERC20')

os.remove('tmp.abi')
os.remove('tmp.bin')
