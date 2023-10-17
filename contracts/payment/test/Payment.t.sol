// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {Payment} from "../src/Payment.sol";
import {MockToken} from "../src/MockToken.sol";

contract PaymentTest is Test {
    Payment public channel;
    MockToken public mktn;
    address provider;
    address client;

    function setUp() public {
        channel = new Payment();
        client = vm.createWallet("client").addr;
        provider = vm.createWallet("provider").addr;
        mktn = new MockToken();
        vm.startPrank(client);
        mktn.ClaimTokens(1000);
    }

    function test_Supply() public {
        uint256 supply = mktn.balanceOf(client);
        assertEq(supply, 1000);
    }

    function test_CreateChannel() public {
        vm.startPrank(client);
        // aprove the payment contract to withdraw 500 of client mktn token
        mktn.approve(address(channel), 500);
        vm.expectRevert("Deadline Expired");
        channel.createChannel(provider, address(mktn), 500, 1, 5, 5);

        channel.createChannel(provider, address(mktn), 500, 2, 5, 5);

        uint256 supply = mktn.balanceOf(client);
        uint256 supplySC = mktn.balanceOf(address(channel));
        assertEq(500, supply, "failed to lock funds");
        assertEq(500, supplySC, "smart contract did not receive the funds");

        mktn.approve(address(channel), 500);
        vm.expectRevert("Channel already created");
        channel.createChannel(provider, address(mktn), 500, 3, 5, 5);

        mktn.approve(address(channel), 0);
        vm.expectRevert("allowance does not match specified amount");
        channel.createChannel(provider, address(mktn), 500, 3, 5, 5);
    }

    function test_UnlockFunds() public {
        vm.startPrank(client);

        // aprove the contract to withdraw 500 of client mktn token
        mktn.approve(address(channel), 500);

        channel.createChannel(provider, address(mktn), 500, 2, 5, 5);
        // unlock the funds before deadline expires
        vm.expectRevert("Deadline not reached yet");
        channel.unclockFunds(address(mktn), provider);
        // advance the block timestamp
        vm.warp(2);
        // withdraw the funds
        channel.unclockFunds(address(mktn), provider);
        uint256 supply = mktn.balanceOf(client);
        assertEq(1000, supply, "failed to withdraw the funds");
        // withdraw empty funds
        vm.expectRevert("Empty Channel");
        channel.unclockFunds(address(mktn), provider);
    }

    function test_LockFunds() public {
		vm.startPrank(client);
        mktn.approve(address(channel), 500);
        channel.createChannel(provider, address(mktn), 500, 2, 5, 5);
        mktn.approve(address(channel), 500);
        channel.lockFunds(provider, address(mktn), 500);
        uint256 balance = mktn.balanceOf(address(channel));
        assertEq(balance, 1000);
    }

    function test_withdraw() public {
        vm.startPrank(client);

        mktn.approve(address(channel), 500);
        channel.createChannel(provider, address(mktn), 500, 2, 5, 5);
        vm.stopPrank();

        vm.startPrank(provider);
        channel.uploadMetrics(client, address(mktn),5);
        channel.withdraw(address(mktn), client);
        uint256 balance = mktn.balanceOf(provider);
        assertEq(balance, 25);
    }

    function test_UpdatePrice() public {
        vm.startPrank(client);
        mktn.approve(address(channel), 500);
        channel.createChannel(provider, address(mktn), 500, 2, 5, 5);

        vm.startPrank(provider);
        channel.updatePrice(client, address(mktn), 10);

        vm.startPrank(client);
        channel.acceptNewPrice(provider, address(mktn));

        vm.startPrank(provider);
        channel.uploadMetrics(client, address(mktn), 5);
        channel.withdraw(address(mktn), client);

        uint256 balance = mktn.balanceOf(provider);
        assertEq(balance, 50);
    }
}
