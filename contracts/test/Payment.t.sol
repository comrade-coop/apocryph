// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {Payment} from "../src/Payment.sol";
import {MockToken} from "../src/MockToken.sol";
import {IERC20Errors} from "openzeppelin-contracts/contracts/interfaces/draft-IERC6093.sol";

contract PaymentTest is Test {
    Payment public payment;
    MockToken public token;
    address provider;
    address publisher;

    function setUp() public {
        payment = new Payment();
        publisher = vm.createWallet("publisher").addr;
        provider = vm.createWallet("provider").addr;
        token = new MockToken();
    }

    function test_createChannel() public {
        vm.startPrank(publisher);
        token.mint(1000);

        vm.expectRevert();
        payment.createChannel(provider, token, 1, 500);

        token.approve(address(payment), 500);
        payment.createChannel(provider, token, 1, 500);

        assertEq(500, token.balanceOf(address(payment)));

        token.approve(address(payment), 500);

        vm.expectRevert(Payment.AlreadyExists.selector);
        payment.createChannel(provider, token, 1, 500);
    }

    function test_deposit() public {
        vm.startPrank(publisher);
        token.mint(1000);

        token.approve(address(payment), 500);
        payment.createChannel(provider, token, 1, 500);

        token.approve(address(payment), 500);
        payment.deposit(provider, token, 500);

        assertEq(token.balanceOf(address(payment)), 1000);
    }

    function test_withdraw() public {
        vm.startPrank(publisher);
        token.mint(500);
        token.approve(address(payment), 500);
        payment.createChannel(provider, token, 1, 500);

        vm.startPrank(provider);
        token.mint(100);

        payment.withdraw(publisher, token, 25, address(0));
        assertEq(token.balanceOf(provider), 125);
        payment.withdraw(publisher, token, 25, address(0));
        assertEq(token.balanceOf(provider), 150);

        vm.expectRevert(Payment.AmountRequired.selector);
        payment.withdrawUpTo(publisher, token, 25, address(0));
        vm.expectRevert(Payment.AmountRequired.selector);
        payment.withdrawUpTo(publisher, token, 50, address(0));

        payment.withdrawUpTo(publisher, token, 100, address(0));
        assertEq(token.balanceOf(provider), 200);

        vm.expectRevert(Payment.InsufficientFunds.selector);
        payment.withdrawUpTo(publisher, token, 501, address(1));

        payment.withdrawUpTo(publisher, token, 500, address(1));
        assertEq(token.balanceOf(provider), 200);
        assertEq(token.balanceOf(address(1)), 400);
    }

    function test_unlock() public {
        vm.startPrank(publisher);
        token.mint(500);
        token.approve(address(payment), 1000);

        vm.expectRevert(Payment.DoesNotExist.selector);
        payment.unlock(provider, token);

        payment.createChannel(provider, token, 20, 500);

        vm.expectRevert(Payment.ChannelLocked.selector);
        payment.withdrawUnlocked(provider, token);

        vm.startPrank(provider);
        payment.withdraw(publisher, token, 100, address(0));
        assertEq(100, token.balanceOf(provider));
        vm.startPrank(publisher);

        vm.expectRevert(Payment.ChannelLocked.selector);
        payment.withdrawUnlocked(provider, token);

        payment.unlock(provider, token);
        // advance the block timestamp
        vm.warp(block.timestamp + 20);

        payment.withdrawUnlocked(provider, token);
        assertEq(400, token.balanceOf(publisher));

        vm.expectRevert(Payment.AmountRequired.selector);
        payment.withdrawUnlocked(provider, token);

        vm.expectRevert(Payment.AlreadyExists.selector);
        payment.createChannel(provider, token, 20, 400);

        payment.closeChannel(provider, token);

        payment.createChannel(provider, token, 20, 400);
    }

    function test_unlock_withdraw() public {
        vm.startPrank(publisher);
        token.mint(500);
        token.approve(address(payment), 500);

        payment.createChannel(provider, token, 20, 500);

        payment.unlock(provider, token);
        vm.warp(10);

        vm.expectRevert(Payment.ChannelLocked.selector);
        payment.withdrawUnlocked(provider, token);

        vm.startPrank(provider);
        payment.withdrawUpTo(publisher, token, 200, address(0));
        assertEq(token.balanceOf(provider), 200);

        vm.startPrank(publisher);

        payment.withdrawUnlocked(provider, token);
        assertEq(token.balanceOf(publisher), 300);
    }
}
