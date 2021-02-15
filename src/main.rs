use contracts::ContractsProvider;
use anyhow::Result;
use tokio;

mod contracts;

#[tokio::main]
async fn main() -> Result<()> {
    let contracts = ContractsProvider::new("http://localhost:8545").await?;

    dbg!(contracts.accounts().await?);

    Ok(())
}
