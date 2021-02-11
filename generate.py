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
    os.system(
        f'abigen --abi tmp.abi --pkg {pkg} --type {type} --out {filename}')


def generate(artifact):
    name = Path(artifact).with_suffix('').name
    base = Path('hubble-contracts/artifacts/contracts')

    path = base / artifact
    lower = name.lower()
    filename = (Path('contracts') / lower / lower).with_suffix('.go')
    generate_bindings(path, name, lower, filename)


os.system('rm -rf contracts')

generate('rollup/Rollup.sol/Rollup.json')
generate('client/FrontendGeneric.sol/FrontendGeneric.json')
generate('DepositManager.sol/DepositManager.json')
generate('TokenRegistry.sol/TokenRegistry.json')
generate('proposers/BurnAuction.sol/BurnAuction.json')
generate('client/FrontendTransfer.sol/FrontendTransfer.json')
generate('client/FrontendMassMigration.sol/FrontendMassMigration.json')
generate('client/FrontendCreate2Transfer.sol/FrontendCreate2Transfer.json')
generate('BLSAccountRegistry.sol/BLSAccountRegistry.json')

os.remove('tmp.abi')
