// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {IERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/IERC20.sol";

contract Payment {
    //events
    event NewPriceUpdate(address client, address provider, address token, uint256 price);
    event NewPrice(address client, address provider, address token, uint256 price);

    struct Channel {
        uint256 total;
        uint256 owedAmount;
        uint256 deadline;
        uint256 price; // execution price
        uint256 suggestedPrice; // used by the provider to update the price
        uint256 minAdvanceDuration; // minimum time in seconds allowed to update the deadline
    }
    // client => provider => token => funds

    mapping(address => mapping(address => mapping(address => Channel))) public channels;

    modifier notZero(uint256 amount) {
        require(amount > 0, "can't accpet zero as a value");
        _;
    }

    modifier notExpired(address client, address provider, address token) {
        Channel memory funds = channels[client][provider][token];
        require(block.timestamp < funds.deadline, "Channel Expired");
        _;
    }

    modifier allowedTransfer(address token, address owner, uint256 amount) {
        uint256 allowance = IERC20(token).allowance(owner, address(this));
        require(allowance == amount, "allowance does not match specified amount");
        _;
    }

    // caller of this function must approve the amount to be withdrawn by this contract address
    function createChannel(
        address provider,
        address token,
        uint256 amount,
        uint256 deadline,
        uint256 minAdvanceDuration,
        uint256 price
    ) public notZero(amount) allowedTransfer(token, msg.sender, amount) {
        require(deadline > block.timestamp, "Deadline Expired");
        Channel memory funds = channels[msg.sender][provider][token];
        require(funds.total == 0, "Channel already created");
        funds = Channel(amount, 0, deadline, price, 0, minAdvanceDuration);

        channels[msg.sender][provider][token] = funds;
        IERC20(token).transferFrom(msg.sender, address(this), funds.total);
    }

    // add more funds to the payment channel, also must be proceeded with client approval
    function lockFunds(address provider, address token, uint256 amount)
        public
        notZero(amount)
        notExpired(msg.sender, provider, token)
    {
        Channel memory funds = channels[msg.sender][provider][token];
        channels[msg.sender][provider][token].total = funds.total + amount;
        IERC20(token).transferFrom(msg.sender, address(this), amount);
    }

    // check if not an empty channel
    // check whether the deadline is reached or not
    // transfer what is left to the provider that is not yet claimed
    // transfer the available funds to the client
    function unclockFunds(address token, address provider) public {
        Channel memory funds = channels[msg.sender][provider][token];
        require(block.timestamp >= funds.deadline, "Deadline not reached yet");
        if (funds.owedAmount > 0) {
            IERC20(token).transfer(provider, funds.owedAmount);
        }
        require(funds.total > 0, "Empty Channel");
        IERC20(token).transfer(msg.sender, funds.total);
        channels[msg.sender][provider][token].total = 0;
        channels[msg.sender][provider][token].owedAmount = 0;
    }

    // allows the provider to withdraw the owedAmount
    function withdraw(address token, address client) public {
        Channel memory funds = channels[client][msg.sender][token];
        require(funds.owedAmount > 0, "Zero Ownership");
        IERC20(token).transfer(msg.sender, funds.owedAmount);
    }

    // this will start the dispute period
    function updatePrice(address client, address token, uint256 price) public notZero(price) {
        channels[client][msg.sender][token].suggestedPrice = price;
        emit NewPriceUpdate(client, msg.sender, token, price);
    }
    // Accepts a new price for a specified token in a payment channel
    // and trigger an event for the updated price.

    function acceptNewPrice(address provider, address token) public {
        Channel memory funds = channels[msg.sender][provider][token];
        channels[msg.sender][provider][token].price = funds.suggestedPrice;
        channels[msg.sender][provider][token].suggestedPrice = 0;
        emit NewPrice(msg.sender, msg.sender, token, funds.price);
    }
    // updates the deadline
    // check if the newDeadline is bigger than the minimum amount agreed on

    function updateDeadline(address provider, address token, uint256 newDeadline) public {
        Channel memory funds = channels[msg.sender][provider][token];
        require(newDeadline > funds.minAdvanceDuration, "New Deadline too short");
        require(newDeadline > funds.deadline, "New Deadline is less than current deadline");
        funds.deadline = newDeadline;
        channels[msg.sender][msg.sender][token] = funds;
    }

    // verify the metric uploaded
    // TODO calculate the amount that should be owed, based on the uploaded metric after verifying it
    function uploadMetrics(address client, address token, uint256 units) public {
        Channel memory funds = channels[client][msg.sender][token];
        // multiply units of execution (derived from the metric msg) with the price per execution
        uint256 amount =  units * funds.price;
        require(amount + funds.owedAmount <= funds.total, "OwedAmount is bigger than channel's available funds");
        channels[client][msg.sender][token].owedAmount = amount + funds.owedAmount;
    }

    // TODO
    // function updateMinDeadline
}
