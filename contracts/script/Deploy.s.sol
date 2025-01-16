// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.22;

import {Script, console2} from "forge-std/Script.sol";
import {PaymentV2} from "../src/PaymentV2.sol";
import {MockToken} from "../src/MockToken.sol";

contract DeployScript is Script {
    function setUp() public {}

    function run() public returns (MockToken token, PaymentV2 paymentV2) {
        vm.resetNonce(msg.sender);
        vm.broadcast();
        
        token = new MockToken();
        console2.log("Mock token address: ", address(token));
        vm.broadcast();
        
        new MockToken(); // Advance nonce to maintain compatibility with Apocryph deployment script
        vm.broadcast();
        new MockToken(); // Advance nonce to maintain compatibility with Apocryph deployment script
        vm.broadcast();
        
        paymentV2 = new PaymentV2(token, uint256(30*60));
        console2.log("PaymentV2 address: ", address(paymentV2));
    }
}

