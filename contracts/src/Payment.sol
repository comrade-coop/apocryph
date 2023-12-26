// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.22;

import {IERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";

using SafeERC20 for IERC20;

contract Payment {
    error AlreadyExists();
    error DoesNotExist();
    error AmountRequired();
    error ChannelLocked();
    error InsufficientFunds();

    event UnlockTimerStarted(
        address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 unlockedAt
    );
    event ChannelCreated(address indexed publisher, address indexed provider, bytes32 indexed podId);
    event Deposited(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 depositAmount);
    event Unlocked(address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 unlockedAmount);
    event Withdrawn(
        address indexed publisher, address indexed provider, bytes32 indexed podId, uint256 withdrawnAmount
    );
    event ChannelClosed(address indexed publisher, address indexed provider, bytes32 indexed podId);

    struct Channel {
        uint256 investedByPublisher;
        uint256 withdrawnByProvider;
        uint256 unlockTime; // minimum time in seconds needed to unlock the funds
        uint256 unlockedAt; // time @ unlock + unlockTime
    }

    // publisher => provider => token => PodID => funds
    mapping(address => mapping(address => mapping(bytes32 => Channel))) public channels;

    IERC20 public token;

    constructor(IERC20 _token) {
        token = _token;
    }

    // called by publisher to create a new payment channel; must approve a withdraw by this contract's address
    function createChannel(address provider, bytes32 podId, uint256 unlockTime, uint256 initialAmount) public {
        if (initialAmount == 0) revert AmountRequired();
        address publisher = msg.sender;
        Channel storage channel = channels[publisher][provider][podId];
        if (channel.investedByPublisher != 0) revert AlreadyExists();
        assert(channel.withdrawnByProvider == 0);
        channel.investedByPublisher = initialAmount;
        channel.unlockTime = unlockTime;

        emit Deposited(publisher, provider, podId, initialAmount);

        token.safeTransferFrom(msg.sender, address(this), initialAmount);
    }

    // add more funds to the payment channel
    function deposit(address provider, bytes32 podId, uint256 amount) public {
        if (amount == 0) revert AmountRequired();
        address publisher = msg.sender;
        Channel storage channel = channels[publisher][provider][podId];

        channel.investedByPublisher = channel.investedByPublisher + amount;
        channel.unlockedAt = 0;

        emit Deposited(publisher, provider, podId, amount);

        token.safeTransferFrom(msg.sender, address(this), amount);
    }

    // initiate the process of unlocking the funds stored in the contract
    function unlock(address provider, bytes32 podId) public {
        address publisher = msg.sender;
        Channel storage channel = channels[publisher][provider][podId];
        if (channel.investedByPublisher == 0) revert DoesNotExist();

        uint256 newUnlockedAt = block.timestamp + channel.unlockTime;
        if (channel.unlockedAt == 0 || channel.unlockedAt < newUnlockedAt) {
            channel.unlockedAt = newUnlockedAt;
            emit UnlockTimerStarted(publisher, provider, podId, newUnlockedAt);
        }
    }

    // transfer the now-unlocked funds back to the publisher
    function withdrawUnlocked(address provider, bytes32 podId) public {
        address publisher = msg.sender;
        Channel storage channel = channels[publisher][provider][podId];
        if (channel.unlockedAt == 0 || block.timestamp < channel.unlockedAt) revert ChannelLocked();

        uint256 leftoverFunds = channel.investedByPublisher - channel.withdrawnByProvider;
        if (leftoverFunds == 0) revert AmountRequired();

        channel.investedByPublisher = channel.withdrawnByProvider;

        emit Unlocked(publisher, provider, podId, leftoverFunds);

        token.safeTransfer(publisher, leftoverFunds);
    }

    // withdrawUnlockedFunds and destroy all previous traces of the channel's existence
    function closeChannel(address provider, bytes32 podId) public {
        address publisher = msg.sender;
        Channel storage channel = channels[publisher][provider][podId];
        if (channel.unlockedAt == 0 || block.timestamp < channel.unlockedAt) revert ChannelLocked();

        uint256 leftoverFunds = channel.investedByPublisher - channel.withdrawnByProvider;
        delete channels[publisher][provider][podId];

        if (leftoverFunds != 0) emit Unlocked(publisher, provider, podId, leftoverFunds);
        emit ChannelClosed(publisher, provider, podId);

        if (leftoverFunds != 0) token.safeTransfer(publisher, leftoverFunds);
    }

    // allows the provider to withdraw as many tokens as would be needed to reach totalWithdrawAmount since the opening of the channel
    function withdrawUpTo(address publisher, bytes32 podId, uint256 totalWithdrawAmount, address transferAddress)
        public
    {
        if (transferAddress == address(0)) {
            transferAddress = msg.sender;
        }

        address provider = msg.sender;
        Channel storage channel = channels[publisher][provider][podId];
        if (totalWithdrawAmount > channel.investedByPublisher) revert InsufficientFunds();
        if (totalWithdrawAmount <= channel.withdrawnByProvider) revert AmountRequired();

        uint256 transferAmount = totalWithdrawAmount - channel.withdrawnByProvider;
        channel.withdrawnByProvider = totalWithdrawAmount;

        emit Withdrawn(publisher, provider, podId, transferAmount);

        if (channel.unlockedAt != 0) {
            channel.unlockedAt = block.timestamp;
        }

        token.safeTransfer(transferAddress, transferAmount);
    }

    // allows the provider to withdraw amount more tokens
    function withdraw(address publisher, bytes32 podId, uint256 amount, address transferAddress) public {
        withdrawUpTo(
            publisher, podId, channels[publisher][msg.sender][podId].withdrawnByProvider + amount, transferAddress
        );
    }

    // allows one to check the amount of as-of-yet unclaimed tokens
    function available(address publisher, address provider, bytes32 podId) public view returns (uint256) {
        Channel storage channel = channels[publisher][provider][podId];
        return channel.investedByPublisher - channel.withdrawnByProvider;
    }

    // allows one to check the amount of so-far claimed tokens
    function withdrawn(address publisher, address provider, bytes32 podId) public view returns (uint256) {
        Channel storage channel = channels[publisher][provider][podId];
        return channel.withdrawnByProvider;
    }
}
