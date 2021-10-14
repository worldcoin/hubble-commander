# 🚄 Benchmarks & profiling

## After serialization optimizations

**Date:** 2021-05-24

**Commit link**: [https://github.com/Worldcoin/hubble-commander/commit/1dc3341c2b4f5465577115ac4a986a00bc23b759](https://github.com/Worldcoin/hubble-commander/commit/1dc3341c2b4f5465577115ac4a986a00bc23b759)

**Platform:** AWS x86 8 core / 16 GB RAM (to be checked)

**Throughput**: 190 tx/sec

[after-serialization-optimizations.prof](Benchmarks%20&%20profiling%209a570e4851044d9b9fa7528124270415/after-serialization-optimizations.prof)

## After benchmark parallelization

**Date:** 2021-05-25

**Commit link**: TBD

**Platform:** AWS x86 8 core / 16 GB RAM (to be checked)

**Throughput**: 650 tx/sec

[after-bench-update.prof](Benchmarks%20&%20profiling%209a570e4851044d9b9fa7528124270415/after-bench-update.prof)

### Sync benchmark introduced

**Date:** 2021-08-03

**Commit link:** [https://github.com/Worldcoin/hubble-commander/pull/278](https://github.com/Worldcoin/hubble-commander/pull/278) (change to commit link after merged)

**Platform:** AWS x86 8 core / 16 GB RAM (to be checked)

|Test                                   |TPS |Profile                                 |
|---------------------------------------|----|----------------------------------------|
|Creation (1s block timer on local geth)|700 | [creation-1s-block-timer](benchmarks/) |
|Sync (1s block timer on local geth)    |1003| [sync-1s-block-timer](benchmarks/) |
|Creation (no block timer)              |797 | [creation](benchmarks/) |
|Sync (no block timer)                  |1003| [sync](benchmarks/) |
