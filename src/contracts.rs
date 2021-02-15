use ethabi::Token;
use hex_literal::hex;
use rlp::RlpStream;
use web3::{Web3, contract::{Contract, Options, tokens::Tokenize}, ethabi, signing::Key, signing::keccak256, types::{Address, BlockNumber, Bytes, H160, TransactionParameters, TransactionRequest, U256}};
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

    let artifact = Artifact::from_json(include_str!("../contracts/artifacts/contracts/client/FrontendTransfer.sol/FrontendTransfer.json"))?;
    
    let seckey: secp256k1::key::SecretKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80".parse().unwrap();
    let address = (&seckey).address();

    let nonce = web3.eth().transaction_count(address, Some(BlockNumber::Latest)).await?;

    let tx= TransactionParameters {
      to: None,
      gas: 9_000_000.into(),
      data: Bytes::from(hex::decode(&artifact.bytecode()[2..])?),
      nonce: Some(nonce),
      ..Default::default()
    };   

    let signed = web3.accounts().sign_transaction(tx, &seckey).await?;

    let tx_hash = web3.eth().send_raw_transaction(signed.raw_transaction).await?;

    dbg!(tx_hash);


    let receipt = web3.eth().transaction_receipt(tx_hash).await?;
    dbg!(receipt);

    
    let mut rlp_stream = rlp::RlpStream::new();
    rlp_stream.begin_list(2);
    rlp_stream.append(&address);
    rlp_stream.append(&nonce);

    let rlp = rlp_stream.out();
    println!("rlp={}", hex::encode(&rlp));
    
    let expected_address = keccak256(&rlp)[12..].to_vec();
    println!("expected_address={}", hex::encode(expected_address));


    panic!();

    // // Deploying a contract
    // let frontend_transfer = Contract::deploy(web3.eth(), &artifact.abi())?
    //     .confirmations(0)
    //     .options(Options::with(|opt| {
    //         opt.value = Some(5.into());
    //         opt.gas_price = Some(5.into());
    //         opt.gas = Some(3_000_000.into());
    //     }))
    //     .execute(
    //         artifact.bytecode(),
    //         (),
    //         my_account,
    //     )
    //     .await?;

    // Ok(ContractsProvider {
    //   web3,
    //   frontend_transfer,
    // })
  }

  pub async fn accounts(&self) -> Result<Vec<Address>> {
    Ok(self.web3.eth().accounts().await?)
  }
} 
