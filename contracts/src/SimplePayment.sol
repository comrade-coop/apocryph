// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.22;

import {IERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";
import {Ownable} from "../lib/openzeppelin-contracts/contracts/access/Ownable.sol";

using SafeERC20 for IERC20;

/*
This contract allows the aApp (its owner) to withdraw funds that have been authorized.
It keeps track of the total withdrawn funds per-version in order to facilitate the aApp's internal 

Security-wise, the contract allows its owner to withdraw all and any funds approved to it.
That's why the owner address exists only inside an attested TEE environment, that can be verified to only run code that comes from the aApp's repository:
https://github.com/comrade-coop/s3-aapp
*/

contract SimplePayment is Ownable {
    event Withdraw(address indexed payer, uint64 indexed version, uint256 amount);

    IERC20 public token;
    mapping(address payer => mapping(uint64 version => uint256)) public totalPaid;

    constructor(IERC20 _token) Ownable(msg.sender) {
        token = _token;
    }

    function withdraw(address payerAddress, uint64 version, address withdrawAddress, uint256 amount) public onlyOwner {
        totalPaid[payerAddress][version] += amount;
        emit Withdraw(payerAddress, version, amount);
        token.safeTransferFrom(payerAddress, withdrawAddress, amount);
    }

}
