# Payment channel
payment channel smart contract setup that supports all ERC20 tokens

## Prerequisite

* [foundry](https://book.getfoundry.sh/)

## Test
An ERC20 token called MockToken is setup for testing
```
forge install foundry-rs/forge-std --no-commit
forge install openzeppelin/openzeppelin-contracts --no-commit
forge test
```

## Build 

```
forge build -o build
```

## Deployment

To deploy on testnet/mainnet, you can use the `scripts/Deploy.s.sol` script -- see [forge documentation](https://book.getfoundry.sh/tutorials/solidity-scripting) for details on using scripts. Afterwards, hardcode the contract addresses as defaults for the Publisher client (and/or provider client), or, better yet commit the `./broadcast` folder and take the contract addresses from there -- or, potentially even better, use ENS to resolve the payment/registry contracts.
