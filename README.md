### Stader Node and EthX

ETHx is an innovative liquid staking token developed by Stader, designed to revolutionize Ethereum staking. Our vision for ETHx is to transform the staking experience, providing users with the unprecedented freedom to move and utilize their staked ETH while continuing to earn rewards and engage with the growing DeFi ecosystem.

ETHx is built to reduce technical and capital barriers to running nodes on Ethereum and empowering smaller node operators is of the highest importance for Stader. Stader’s ETHx permissionless pool lets anyone operate a node with 4.4 ETH of asset collateral [4 ETH + 0.4ETH worth of SD (Stader’s governance token)]. 

This repo contains code for the stader-cli which allows users to easily join EthX permissionless pool and become a crucial part in Stader's missions to revolutionize Ethereum Staking!

## Documentation

NOs can find documentation w.r.t setting a system requirements, how to set a node up, the latest binaries etc here https://staderlabs.gitbook.io/ethereum/

## Integration testing

Upcoming

## Safety version features
- VC containers will not be launched if the `allowVCContainers` setting is empty or set to `false`.
- `stader-cli validator deposit` returns error if `createNewValidators` settings param is empty or `false`

## New `user-settings.yaml` params
```yaml
root:
  createNewValidators: "false|true"
  allowVCContainers: "false|true"
```

## Build and run safety run version
```bash
# build cli
./build-release.sh -c -v v1.6.1-safety-run
# run cli
./build/v1.6.1-safety-run/stader-cli-[arch] [args] 
```
