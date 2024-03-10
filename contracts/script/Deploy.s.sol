// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.22;

import {Script, console2} from "forge-std/Script.sol";
import {Payment} from "../src/Payment.sol";
import {MockToken} from "../src/MockToken.sol";
import {Registry} from "../src/Registry.sol";

contract DeployScript is Script {
    function setUp() public {}

    function run() public returns (MockToken token, Payment payment, Registry registry) {
        vm.resetNonce(msg.sender);
        vm.broadcast();
        token = new MockToken();
        console2.log("Token address: ", address(token));
        vm.broadcast();
        payment = new Payment(token);
        console2.log("Payment address: ", address(payment));
        vm.broadcast();
        registry = new Registry();
        console2.log("Registry address: ", address(registry));
    }
}
