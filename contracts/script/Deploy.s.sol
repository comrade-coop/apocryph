// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.22;

import {Script, console2} from "forge-std/Script.sol";
import {MockToken} from "../src/MockToken.sol";

contract DeployScript is Script {
    function setUp() public {}

    function run() public returns (MockToken token) {
        vm.resetNonce(msg.sender);
        vm.broadcast();
        
        token = new MockToken();
        console2.log("Mock token address: ", address(token));
        vm.broadcast();
    }
}

