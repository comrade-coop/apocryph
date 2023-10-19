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
