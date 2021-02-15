use serde::Deserialize;
use serde_json::{Value, from_str, to_vec};
use anyhow::Result;

#[derive(Deserialize)]
pub struct Artifact {
  abi: Value,
  bytecode: String,
}

impl Artifact {
  pub fn from_json(json: &str) -> Result<Self> {
    Ok(from_str(json)?)
  }

  pub fn abi(&self) -> Vec<u8> {
    to_vec(&self.abi).unwrap()
  }

  pub fn bytecode(&self) -> &str {
    &self.bytecode
  }
}

#[cfg(test)]
mod tests {
    use anyhow::Result;
    use super::Artifact;

  #[test]
  fn artifact() -> Result<()> {
    let artifact = Artifact::from_json(include_str!("../../contracts/artifacts/contracts/client/FrontendTransfer.sol/FrontendTransfer.json"))?;
    
    assert!(artifact.abi().len() > 0);
    assert!(artifact.bytecode().len() > 0);

    Ok(())
  }
}
