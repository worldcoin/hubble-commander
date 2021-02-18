import os
import json
import shutil
from pathlib import Path


def generate_bindings(path, type, pkg, filename):
    os.makedirs(Path("contracts") / pkg, exist_ok=True)
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
    filename = (Path('contracts') / lower / lower).with_suffix('.go')
    generate_bindings(path, name, lower, filename)


os.system('rm -rf contracts')

generate('contracts/rollup/Rollup.sol/Rollup.json', 'Rollup')
generate('contracts/client/FrontendGeneric.sol/FrontendGeneric.json', 'Frontend')
generate('contracts/DepositManager.sol/DepositManager.json', 'Deposit')
generate('contracts/TokenRegistry.sol/TokenRegistry.json', 'TokenRegistry')
generate('contracts/proposers/BurnAuction.sol/BurnAuction.json', 'BurnAuction')
generate('contracts/client/FrontendTransfer.sol/FrontendTransfer.json', 'Transfer')
generate('contracts/client/FrontendMassMigration.sol/FrontendMassMigration.json', 'MassMigration')
generate('contracts/client/FrontendCreate2Transfer.sol/FrontendCreate2Transfer.json', 'Create2Transfer')
generate('contracts/BLSAccountRegistry.sol/BLSAccountRegistry.json', 'AccountRegistry')
generate('@openzeppelin/contracts/token/ERC20/IERC20.sol/IERC20.json', 'ERC20')

os.remove('tmp.abi')
os.remove('tmp.bin')
