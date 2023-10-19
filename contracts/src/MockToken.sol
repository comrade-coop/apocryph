// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {ERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";

contract MockToken is ERC20 {
    constructor() ERC20("MockToken", "MKT") {
        _mint(address(this), 10 ** 18);
    }

    function ClaimTokens(uint256 amount) public {
        require(balanceOf(address(this)) > amount, "Contract does not have Enough tokens");
        _transfer(address(this), msg.sender, amount);
    }
}
