// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {PaymentV2, PERMISSION_MANAGE_AUTHORIZATIONS, PERMISSION_WITHDRAW, PERMISSION_NO_LIMIT} from "../src/PaymentV2.sol";
import {MockToken} from "../src/MockToken.sol";
import {IERC20Errors} from "openzeppelin-contracts/contracts/interfaces/draft-IERC6093.sol";

contract PaymentV2Test is Test {
    PaymentV2 public paymentV2;
    MockToken public token;
    address recipient;
    address recipient2;
    address payer;
    bytes32 channelDiscriminator = bytes32(uint256(1234));
    bytes32 channelId;
    uint256 unlockTime = 10;

    function setUp() public {
        payer = vm.createWallet("payer").addr;
        recipient = vm.createWallet("recipient").addr;
        recipient2 = vm.createWallet("recipient2").addr;
        token = new MockToken();
        paymentV2 = new PaymentV2(token, unlockTime);
        channelId = paymentV2.getChannelId(payer, channelDiscriminator);
    }

    function test_create() public {
        vm.startPrank(payer);
        token.mint(1000);

        vm.expectRevert();
        paymentV2.create(channelDiscriminator, 500);

        token.approve(address(paymentV2), 500);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Deposit(channelId, payer, 500);
        paymentV2.create(channelDiscriminator, 500);

        assertEq(500, token.balanceOf(address(paymentV2)));

        token.approve(address(paymentV2), 500);

        vm.expectRevert(PaymentV2.AlreadyInitialized.selector);
        paymentV2.create(channelDiscriminator, 500);
    }

    function test_authorize() public {
        vm.startPrank(payer);
        token.mint(500);

        vm.expectRevert(PaymentV2.NotAuthorized.selector);
        paymentV2.authorize(channelId, recipient, PERMISSION_MANAGE_AUTHORIZATIONS, 123, 400);

        token.approve(address(paymentV2), 500);
        paymentV2.create(channelDiscriminator, 500);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient, PERMISSION_MANAGE_AUTHORIZATIONS | PERMISSION_WITHDRAW, 123, 400);
        paymentV2.authorize(channelId, recipient, PERMISSION_MANAGE_AUTHORIZATIONS | PERMISSION_WITHDRAW, 123, 400);
    }

    function test_authorize_recipient2() public {
        vm.startPrank(payer);
        token.mint(500);

        token.approve(address(paymentV2), 500);
        paymentV2.create(channelDiscriminator, 500);
        paymentV2.authorize(channelId, recipient, PERMISSION_MANAGE_AUTHORIZATIONS | PERMISSION_WITHDRAW, 123, 400);

        vm.startPrank(recipient);

        vm.expectRevert(PaymentV2.NotAuthorized.selector); // Can't pass MANAGE_AUTHORIZATIONS a second level
        paymentV2.authorize(channelId, recipient2, PERMISSION_MANAGE_AUTHORIZATIONS | PERMISSION_WITHDRAW, 100, 100);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient, PERMISSION_MANAGE_AUTHORIZATIONS | PERMISSION_WITHDRAW, 123, 300); // Own limit decreased
        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient2, PERMISSION_WITHDRAW, 100, 100);
        paymentV2.authorize(channelId, recipient2, PERMISSION_WITHDRAW, 100, 100);
    }

    function test_createAndAuthorize() public {
        vm.startPrank(payer);
        token.mint(500);
        token.approve(address(paymentV2), 500);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Deposit(channelId, payer, 500);
        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient, PERMISSION_MANAGE_AUTHORIZATIONS, 123, 400);
        paymentV2.createAndAuthorize(channelDiscriminator, 500, recipient, PERMISSION_MANAGE_AUTHORIZATIONS, 123, 400);
    }

    function test_withdraw() public {
        vm.startPrank(payer);
        token.mint(500);
        token.approve(address(paymentV2), 500);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient, PERMISSION_WITHDRAW, 123, 400);
        paymentV2.createAndAuthorize(channelDiscriminator, 500, recipient, PERMISSION_WITHDRAW, 123, 400);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Withdraw(channelId, payer, 1);
        paymentV2.withdraw(channelId, address(0), 1);
        assertEq(1, token.balanceOf(payer));
        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Withdraw(channelId, payer, 2); // payer, because we are the ones drawing the funds, even if someone else receieves them
        paymentV2.withdraw(channelId, recipient2, 2);

        vm.startPrank(recipient);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient, PERMISSION_WITHDRAW, 123, 200);
        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Withdraw(channelId, recipient, 200);
        paymentV2.withdraw(channelId, address(0), 200);


        vm.expectRevert(PaymentV2.InsufficientFunds.selector);
        paymentV2.withdraw(channelId, address(0), 201);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient, PERMISSION_WITHDRAW, 97, 0); // 20, because we are running out of funds
        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Withdraw(channelId, recipient, 200);
        paymentV2.withdraw(channelId, address(0), 200);
    }

    function test_withdraw_unlimited_recipient2() public {
        vm.startPrank(payer);
        token.mint(500);
        token.approve(address(paymentV2), 500);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient, PERMISSION_WITHDRAW | PERMISSION_NO_LIMIT, 123, 0);
        paymentV2.createAndAuthorize(channelDiscriminator, 500, recipient, PERMISSION_WITHDRAW | PERMISSION_NO_LIMIT, 123, 0);
        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient2, PERMISSION_WITHDRAW | PERMISSION_NO_LIMIT, 300, 0);
        paymentV2.authorize(channelId, recipient2, PERMISSION_WITHDRAW | PERMISSION_NO_LIMIT, 300, 0);

        vm.startPrank(recipient);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient, PERMISSION_WITHDRAW | PERMISSION_NO_LIMIT, 100, 0);
        paymentV2.withdraw(channelId, address(0), 100);

        vm.expectRevert(PaymentV2.InsufficientFunds.selector); // Funds are available, but are reserved to another channel
        paymentV2.withdraw(channelId, address(0), 180);

        vm.startPrank(recipient2);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient2, PERMISSION_WITHDRAW | PERMISSION_NO_LIMIT, 120, 0);
        paymentV2.withdraw(channelId, address(0), 180);
    }

    function test_unlock() public {
        vm.startPrank(payer);
        token.mint(500);
        token.approve(address(paymentV2), 500);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Deposit(channelId, payer, 500);
        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient, PERMISSION_MANAGE_AUTHORIZATIONS, 123, 400);
        paymentV2.createAndAuthorize(channelDiscriminator, 500, recipient, PERMISSION_MANAGE_AUTHORIZATIONS, 123, 400);

        vm.expectRevert(PaymentV2.AuthorizationLocked.selector);
        paymentV2.authorize(channelId, recipient, 0, 0, 0);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Unlock(channelId, recipient);
        paymentV2.unlock(channelId, recipient);

        vm.expectRevert(PaymentV2.AuthorizationLocked.selector);
        paymentV2.authorize(channelId, recipient, 0, 0, 0);

        vm.warp(block.timestamp + 9);

        vm.expectRevert(PaymentV2.AuthorizationLocked.selector);
        paymentV2.authorize(channelId, recipient, 0, 0, 0);

        vm.warp(block.timestamp + 10);

        vm.expectEmit(true, true, false, true);
        emit PaymentV2.Authorize(channelId, recipient, 0, 0, 400); // 400 cause using extraLimit // TODO
        paymentV2.authorize(channelId, recipient, 0, 0, 0);
    }
}
