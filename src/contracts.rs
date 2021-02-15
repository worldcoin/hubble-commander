use hex_literal::hex;
use web3::{Web3, contract::{Contract, Options}, types::{Address, U256}};
use anyhow::Result;

mod artifact;
use artifact::Artifact;

pub struct ContractsProvider {
  web3: Web3<web3::transports::Http>,
  frontend_transfer: Contract<web3::transports::Http>,
}

impl ContractsProvider {
  pub async fn new(rpc: &str) -> Result<Self> {
    let transport = web3::transports::Http::new(rpc)?;
    let web3 = Web3::new(transport);

    let my_account = hex!("d028d24f16a8893bd078259d413372ac01580769").into();

    let artifact = Artifact::from_json(include_str!("../contracts/artifacts/contracts/client/FrontendTransfer.sol/FrontendTransfer.json"))?;
    // Deploying a contract
    let frontend_transfer = Contract::deploy(web3.eth(), &artifact.abi())?
        .confirmations(0)
        .options(Options::with(|opt| {
            opt.value = Some(5.into());
            opt.gas_price = Some(5.into());
            opt.gas = Some(3_000_000.into());
        }))
        .execute(
            artifact.bytecode(),
            (),
            my_account,
        )
        .await?;

    Ok(ContractsProvider {
      web3,
      frontend_transfer,
    })
  }

  pub async fn accounts(&self) -> Result<Vec<Address>> {
    Ok(self.web3.eth().accounts().await?)
  }
} 
