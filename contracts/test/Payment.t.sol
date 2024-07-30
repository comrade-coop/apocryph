// SPDX-License-Identifier: GPL-3.0

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
    bytes32 podId;

    function setUp() public {
        publisher = vm.createWallet("publisher").addr;
        provider = vm.createWallet("provider").addr;
        token = new MockToken();
        payment = new Payment(token);
        podId = bytes32(0);
    }

    function test_createChannel() public {
        vm.startPrank(publisher);
        token.mint(1000);

        vm.expectRevert();
        payment.createChannel(provider, podId, 1, 500);

        token.approve(address(payment), 500);
        payment.createChannel(provider, podId, 1, 500);

        assertEq(500, token.balanceOf(address(payment)));

        token.approve(address(payment), 500);

        vm.expectRevert(Payment.AlreadyExists.selector);
        payment.createChannel(provider, podId, 1, 500);
    }

    function test_deposit() public {
        vm.startPrank(publisher);
        token.mint(1000);

        token.approve(address(payment), 500);
        vm.expectRevert(Payment.DoesNotExist.selector);
        payment.deposit(provider, podId, 500);

        payment.createChannel(provider, podId, 1, 500);

        token.approve(address(payment), 500);
        payment.deposit(provider, podId, 500);

        assertEq(token.balanceOf(address(payment)), 1000);
    }

    function test_withdraw() public {
        vm.startPrank(publisher);
        token.mint(500);
        token.approve(address(payment), 500);
        payment.createChannel(provider, podId, 1, 500);

        vm.startPrank(provider);
        token.mint(100);

        payment.withdraw(publisher, podId, 25, address(0));
        assertEq(token.balanceOf(provider), 125);
        assertEq(payment.withdrawn(publisher, provider, podId), 25);
        assertEq(payment.available(publisher, provider, podId), 475);
        payment.withdraw(publisher, podId, 25, address(0));
        assertEq(token.balanceOf(provider), 150);
        assertEq(payment.withdrawn(publisher, provider, podId), 50);
        assertEq(payment.available(publisher, provider, podId), 450);

        vm.expectRevert(Payment.AmountRequired.selector);
        payment.withdrawUpTo(publisher, podId, 25, address(0));
        vm.expectRevert(Payment.AmountRequired.selector);
        payment.withdrawUpTo(publisher, podId, 50, address(0));

        payment.withdrawUpTo(publisher, podId, 100, address(0));
        assertEq(token.balanceOf(provider), 200);
        assertEq(payment.withdrawn(publisher, provider, podId), 100);
        assertEq(payment.available(publisher, provider, podId), 400);

        vm.expectRevert(Payment.InsufficientFunds.selector);
        payment.withdrawUpTo(publisher, podId, 501, address(1));

        payment.withdrawUpTo(publisher, podId, 500, address(1));
        assertEq(token.balanceOf(provider), 200);
        assertEq(token.balanceOf(address(1)), 400);
        assertEq(payment.withdrawn(publisher, provider, podId), 500);
        assertEq(payment.available(publisher, provider, podId), 0);
    }

    function test_unlock() public {
        vm.startPrank(publisher);
        token.mint(500);
        token.approve(address(payment), 1000);

        vm.expectRevert(Payment.DoesNotExist.selector);
        payment.unlock(provider, podId);

        payment.createChannel(provider, podId, 20, 500);

        vm.expectRevert(Payment.ChannelLocked.selector);
        payment.withdrawUnlocked(provider, podId);

        vm.startPrank(provider);
        payment.withdraw(publisher, podId, 100, address(0));
        assertEq(100, token.balanceOf(provider));
        vm.startPrank(publisher);

        vm.expectRevert(Payment.ChannelLocked.selector);
        payment.withdrawUnlocked(provider, podId);

        payment.unlock(provider, podId);
        // advance the block timestamp
        vm.warp(block.timestamp + 20);

        payment.withdrawUnlocked(provider, podId);
        assertEq(400, token.balanceOf(publisher));

        vm.expectRevert(Payment.AmountRequired.selector);
        payment.withdrawUnlocked(provider, podId);

        vm.expectRevert(Payment.AlreadyExists.selector);
        payment.createChannel(provider, podId, 20, 400);

        payment.closeChannel(provider, podId);

        payment.createChannel(provider, podId, 20, 400);
    }

    function test_unlock_withdraw() public {
        vm.startPrank(publisher);
        token.mint(500);
        token.approve(address(payment), 500);

        payment.createChannel(provider, podId, 20, 500);

        payment.unlock(provider, podId);
        vm.warp(10);

        vm.expectRevert(Payment.ChannelLocked.selector);
        payment.withdrawUnlocked(provider, podId);

        vm.startPrank(provider);
        payment.withdrawUpTo(publisher, podId, 200, address(0));
        assertEq(token.balanceOf(provider), 200);

        vm.startPrank(publisher);

        payment.withdrawUnlocked(provider, podId);
        assertEq(token.balanceOf(publisher), 300);
    }

    function test_authorize() public {
        // Setup
        vm.startPrank(publisher);
        token.mint(1000);
        token.approve(address(payment), 1000);
        payment.createChannel(provider, podId, 1, 500);

        // Test authorizing an address
        address authorizedAddr = address(0x123);
        payment.authorize(authorizedAddr, provider, podId);

        // Verify authorization (we'll need to add a getter function in the contract for this)
        assertTrue(payment.isAuthorized(publisher, provider, podId, authorizedAddr));

        // Test authorizing for non-existent channel
        bytes32 nonExistentPodId = keccak256("non-existent");
        vm.expectRevert(Payment.DoesNotExist.selector);
        payment.authorize(authorizedAddr, provider, nonExistentPodId);

        vm.stopPrank();
    }

    function test_createSubChannel() public {
        // Setup
        vm.startPrank(publisher);
        // mint 1000 to the caller (publisher)
        token.mint(1000);
        token.approve(address(payment), 1000);
        payment.createChannel(provider, podId, 1, 500);
        bytes32 newPodId = bytes32(uint256(1));

        address authorizedAddr = address(0x123);
        payment.authorize(authorizedAddr, provider, podId);
        vm.stopPrank();

        // Test creating a sub-channel
        vm.startPrank(authorizedAddr);
        address newProvider = address(0x456);
        payment.createSubChannel(publisher, provider, podId, newProvider,
                                 newPodId, 200);

        // Verify sub-channel creation
        (uint256 investedAmount, uint256 withdrawnAmount, uint256 unlockTime,) =
            payment.channels(authorizedAddr, newProvider, newPodId);
        assertEq(investedAmount, 200);
        assertEq(withdrawnAmount, 0);
        assertEq(unlockTime, 1); // should be the same as the main channel's unlock time

        // Verify main channel balance reduction
        (uint256 mainChannelBalance,,,) = payment.channels(publisher, provider, podId);
        assertEq(mainChannelBalance, 300);

        // Test creating sub-channel with insufficient funds
        vm.expectRevert(Payment.InsufficientFunds.selector);
        payment.createSubChannel(publisher, provider, podId, newProvider,
                                 newPodId,400);

        // Test creating sub-channel from non-existent main channel
        bytes32 nonExistentPodId = keccak256("non-existent");
        vm.expectRevert(Payment.DoesNotExist.selector);
        payment.createSubChannel(publisher, provider, nonExistentPodId,
                                 newProvider, newPodId,100);

        // Test creating sub-channel without authorization
        vm.stopPrank();
        vm.startPrank(address(0x789)); // Non-authorized address
        vm.expectRevert(Payment.NotAuthorized.selector);
        payment.createSubChannel(publisher, provider, podId,
                                 newProvider,newPodId, 100);

        vm.stopPrank();
    }
}
