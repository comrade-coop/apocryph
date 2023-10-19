// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {IERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/IERC20.sol";

contract Payment {
    //events
    event NewPriceUpdate(address client, address provider, address token, uint256 price);
    event NewPrice(address client, address provider, address token, uint256 price);
    event PaymentChannel(address client, address provider, uint256 podID, address token);

    uint256 PodID = 0;

    struct Channel {
        uint256 total;
        uint256 owedAmount;
        uint256 deadline;
        uint256 price; // execution price
        uint256 suggestedPrice; // used by the provider to update the price
        uint256 minAdvanceDuration; // minimum time in seconds allowed to update the deadline
    }
    // client => provider => podID => token => funds

    mapping(address => mapping(address => mapping(uint256 => mapping(address => Channel)))) public channels;

    modifier notZero(uint256 amount) {
        require(amount > 0, "can't accpet zero as a value");
        _;
    }

    modifier notExpired(address client, address provider, uint256 podID, address token) {
        Channel memory funds = channels[client][provider][podID][token];
        require(block.timestamp < funds.deadline, "Channel Expired");
        _;
    }

    modifier allowedTransfer(address token, address owner, uint256 amount) {
        uint256 allowance = IERC20(token).allowance(owner, address(this));
        require(allowance == amount, "allowance does not match specified amount");
        _;
    }

    function incrementId() public {
        PodID = PodID + 1;
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
        Channel memory funds = channels[msg.sender][provider][PodID + 1][token];
        funds = Channel(amount, 0, deadline, price, 0, minAdvanceDuration);

        channels[msg.sender][provider][PodID + 1][token] = funds;
        IERC20(token).transferFrom(msg.sender, address(this), funds.total);
        incrementId();
    }

    // add more funds to the payment channel, also must be proceeded with client approval
    function lockFunds(address provider, uint256 podID, address token, uint256 amount)
        public
        notZero(amount)
        notExpired(msg.sender, provider, podID, token)
    {
        Channel memory funds = channels[msg.sender][provider][podID][token];
        channels[msg.sender][provider][podID][token].total = funds.total + amount;
        IERC20(token).transferFrom(msg.sender, address(this), amount);
    }

    // check if not an empty channel
    // check whether the deadline is reached or not
    // transfer what is left to the provider that is not yet claimed
    // transfer the available funds to the client
    function unclockFunds(address token, address provider, uint256 podID) public {
        Channel memory funds = channels[msg.sender][provider][podID][token];
        require(block.timestamp >= funds.deadline, "Deadline not reached yet");
        if (funds.owedAmount > 0) {
            IERC20(token).transfer(provider, funds.owedAmount);
        }
        require(funds.total > 0, "Empty Channel");
        IERC20(token).transfer(msg.sender, funds.total);
        channels[msg.sender][provider][podID][token].total = 0;
        channels[msg.sender][provider][podID][token].owedAmount = 0;
    }

    // allows the provider to withdraw the owedAmount
    function withdraw(address client, uint256 podID, address token) public {
        Channel memory funds = channels[client][msg.sender][podID][token];
        require(funds.owedAmount > 0, "Zero Ownership");
        IERC20(token).transfer(msg.sender, funds.owedAmount);
        channels[client][msg.sender][podID][token].owedAmount = 0;
    }

    // this will start the dispute period
    function updatePrice(address client, uint256 podID, address token, uint256 price)
        public
        notZero(price)
        notExpired(client, msg.sender, podID, token)
    {
        channels[client][msg.sender][podID][token].suggestedPrice = price;
        emit NewPriceUpdate(client, msg.sender, token, price);
    }
    // Accepts a new price for a specified token in a payment channel
    // and trigger an event for the updated price.

    function acceptNewPrice(address provider, uint256 podID, address token)
        public
        notExpired(msg.sender, provider, podID, token)
    {
        Channel memory funds = channels[msg.sender][provider][podID][token];
        channels[msg.sender][provider][podID][token].price = funds.suggestedPrice;
        channels[msg.sender][provider][podID][token].suggestedPrice = 0;
        emit NewPrice(msg.sender, msg.sender, token, funds.price);
    }
    // updates the deadline
    // check if the newDeadline is bigger than the minimum amount agreed on

    function updateDeadline(address provider, uint256 podID, address token, uint256 newDeadline) public {
        Channel memory funds = channels[msg.sender][provider][podID][token];
        require(newDeadline > funds.minAdvanceDuration, "New Deadline too short");
        require(newDeadline > funds.deadline, "New Deadline is less than current deadline");
        funds.deadline = newDeadline;
        channels[msg.sender][msg.sender][podID][token] = funds;
    }

    // verify the metric uploaded
    // TODO calculate the amount that should be owed, based on the uploaded metric after verifying it
    function uploadMetrics(address client, uint256 podID, address token, uint256 units) public {
        Channel memory funds = channels[client][msg.sender][podID][token];
        // multiply units of execution (derived from the metric msg) with the price per execution
        uint256 amount = units * funds.price;
        require(amount + funds.owedAmount <= funds.total, "OwedAmount is bigger than channel's available funds");
        channels[client][msg.sender][podID][token].owedAmount = amount + funds.owedAmount;
    }

    // TODO
    // function updateMinDeadline
}
