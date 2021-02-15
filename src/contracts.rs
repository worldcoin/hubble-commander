use web3::{Web3, types::Address};
use anyhow::Result;

pub struct ContractsProvider {
  web3: Web3<web3::transports::Http>,
}

impl ContractsProvider {
  pub fn new(rpc: &str) -> Result<Self> {
    let transport = web3::transports::Http::new(rpc)?;
    let web3 = Web3::new(transport);


    Ok(ContractsProvider {
      web3,
    })
  }

  pub async fn accounts(&self) -> Result<Vec<Address>> {
    Ok(self.web3.eth().accounts().await?)
  }
} 
