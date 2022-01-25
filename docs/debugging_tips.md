# Debugging Tips

## Stacktraces for Rollup Loop Panics

1. Start geth in a new terminal: 
   ```bash
   make start-geth-locally
   ```

2. Set up your environment:
   ```bash
   export HUBBLE_ETHEREUM_CHAIN_ID=1337
   export HUBBLE_ETHEREUM_RPC_URL=ws://localhost:8546
   export HUBBLE_ETHEREUM_PRIVATE_KEY=ee79b5f6e221356af78cf4c36f4f7885a11b67dfcc81c34d80249947330c0f82
   ```

3. Deploy the contracts to the new geth instance (in a terminal with the above env!)

   ```bash
   make deploy
   ```

   You will need to do this every time you restart geth (`start-geth-locally` erases
   previous state)

4. Start debugging the commander:

   ```bash
   HUBBLE_BOOTSTRAP_CHAIN_SPEC_PATH=chain-spec.yaml HUBBLE_BOOTSTRAP_PRUNE=true HUBBLE_ROLLUP_DISABLE_SIGNATURES=true dlv debug ./main -- start
   ```

   From inside delve you'll need to set a breakpoint to catch the panic. As of writing
   the correct spot is in `(*Commander).startWorker`, in commander/commander.go, if you
   set your breakpoint after the call to `recover` then you will be able to inspect the
   stack trace as of the panic.

   The interaction looks something like this:

   ```
   (dlv) break commander/commander.go:164
   Breakpoint 1 set at 0x117eb39 for github.com/Worldcoin/hubble-commander/commander.(*Commander).startWorker.func1.1() ./commander/commander.go:164
   (dlv) continue
   [a lot of log output while you perform step 5]
   (dlv) bt
   [the desired stack trace!]
   ```

5. Do the thing which causes the commander to panic. There's a good chance this looks like
   one of the following (in a shell with the environment from step 2):

   ```bash
   make test-e2e-in-process
   TEST=TestMeasureDisputeGasUsage make test-e2e-locally
   ```
