// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
contract MockToken is ERC20 {
    constructor(address admin) ERC20("MockToken", "MKT")
    {
        _mint(admin, 10**18 );
    }
}
