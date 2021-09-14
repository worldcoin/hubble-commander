# State Rent

### Background

- Hubble state tree is limited in size
    - Making it deeper, makes it more expensive to run fraud proofs and more expensive to run a commander

### Possible Approaches

- If the sender account is below some (dynamically calculated) threshold balance, anybody can submit a transaction that “steals” those funds
    - We can do this rather easily by requiring a minimum balance in the *signature* fraud proof
    - We also need to adapt the *transfer* fraud proof to allow to replace 0 balance leaves by an empty leaf
- Of course it doesn’t help anything, if the threshold balance is static
    - Expected behavior is that you pay a fixed amount of the given token per batch (during the time your state leaf exists in the tree)
    - Calculate: `min_balance = rent_per_batch * (current_batch_number - creation_batch_number)`
- Possible problems:
    - Operator might censor transactions of users close to the minimum to be able to claim the rent
        - It would be better if the operator has to burn it
        - However, then there is not incentive for the operator as he has to pay for call data of the transaction
    - Token rent in USD will heavily fluctuate
        - It doesn’t have to, we can get the current token price and denominate rent fees in USD
        - Might be harder for commander to track
        - Token price fluctuations might have dangerous side effects of liquidating a lot of accounts