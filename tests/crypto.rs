use amcl::rand::RAND;
use amcl::bn254::bls;
use amcl::arch;

pub fn printbinary(array: &[u8]) {
  for i in 0..array.len() {
      print!("{:02X}", array[i])
  }
  println!("")
}


#[test]
fn bn254() {

  let mut raw: [u8; 100] = [0; 100];

  let mut rng = RAND::new();
  rng.clean();
  for i in 0..100 {
      raw[i] = i as u8
  }

  rng.seed(100, &raw);

  println!("{} bit build", arch::CHUNK);


  const BFS: usize = bls::BFS;
  const BGS: usize = bls::BGS;

  let mut s: [u8; BGS] = [0; BGS];

  const G1S: usize = BFS + 1; /* Group 1 Size */
  const G2S: usize = 4 * BFS; /* Group 2 Size */

  let mut w: [u8; G2S] = [0; G2S];
  let mut sig: [u8; G1S] = [0; G1S];

  let m = String::from("This is a test message");

  bls::key_pair_generate(&mut rng, &mut s, &mut w);
  print!("Private key : 0x");
  printbinary(&s);
  print!("Public  key : 0x");
  printbinary(&w);

  bls::sign(&mut sig, &m, &s);
  print!("Signature : 0x");
  printbinary(&sig);

  let res = bls::verify(&sig, &m, &w);

  assert_eq!(0, res);
}
