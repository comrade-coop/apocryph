// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.13;

import {Test, console2} from "forge-std/Test.sol";
import {SimplePayment} from "../src/SimplePayment.sol";
import {MockToken} from "../src/MockToken.sol";
import {IERC20} from "openzeppelin-contracts/contracts/token/ERC20/IERC20.sol";
import {IERC20Errors} from "openzeppelin-contracts/contracts/interfaces/draft-IERC6093.sol";
import {Ownable} from "../lib/openzeppelin-contracts/contracts/access/Ownable.sol";

contract SimplePaymentTest is Test {
    address aapp;
    address withdrawWallet;
    address payer;
    SimplePayment public payment;
    MockToken public token;

    function setUp() public {
        aapp = vm.createWallet("aapp").addr;
        withdrawWallet = vm.createWallet("withdrawWallet").addr;
        payer = vm.createWallet("payer").addr;
        token = new MockToken();
//         vm.startPrank(aapp);
        payment = new SimplePayment(token);
        payment.transferOwnership(aapp);
    }

    function test_withdraw() public {
        vm.startPrank(payer);
        token.mint(500);
        token.approve(address(payment), 500);

        vm.startPrank(payer);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, payer));
        payment.withdraw(payer, 1, withdrawWallet, 400);

        vm.startPrank(withdrawWallet);
        vm.expectRevert(abi.encodeWithSelector(Ownable.OwnableUnauthorizedAccount.selector, withdrawWallet));
        payment.withdraw(payer, 1, withdrawWallet, 400);
        
        assertEq(payment.totalPaid(payer, 1), 0);
        assertEq(payment.totalPaid(payer, 2), 0);
        
        vm.startPrank(aapp);
        vm.expectEmit(address(payment));
        emit SimplePayment.Withdraw(payer, 1, 400);
        vm.expectEmit(address(token));
        emit IERC20.Transfer(payer, withdrawWallet, 400);
        payment.withdraw(payer, 1, withdrawWallet, 400);
        
        assertEq(payment.totalPaid(payer, 1), 400);
        
        vm.startPrank(aapp);
        vm.expectEmit(address(payment));
        emit SimplePayment.Withdraw(payer, 1, 50);
        vm.expectEmit(address(token));
        emit IERC20.Transfer(payer, withdrawWallet, 50);
        payment.withdraw(payer, 1, withdrawWallet, 50);
        
        assertEq(payment.totalPaid(payer, 1), 450);
        assertEq(payment.totalPaid(payer, 2), 0);
    }
}
