// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.22;

import {IERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";

using SafeERC20 for IERC20;

uint16 constant PERMISSION_ADMIN = 1;
uint16 constant PERMISSION_MANAGE_AUTHORIZATIONS = 2;
uint16 constant PERMISSION_WITHDRAW = 4;
uint16 constant PERMISSION_NO_LIMIT = 8;

contract PaymentV2 {
    error AuthorizationLocked();
    error NotAuthorized();
    error InsufficientFunds();
    error AlreadyInitialized();

    event Deposit(bytes32 indexed channelId, address indexed payer, uint256 amount);
    event Authorize(bytes32 indexed channelId, address indexed recipient, uint16 permissions, uint256 reservedAmount, uint256 limit);
    event Unlock(bytes32 indexed channelId, address indexed recipient);
    event Withdraw(bytes32 indexed channelId, address indexed recipient, uint256 amount);

    struct Channel {
        uint128 available;
        uint128 reserved;
    }
    struct ChannelAuthorization {
        uint128 reservation;
        uint128 limit;

        uint240 unlocksAt;
        uint16 permissions;
    }

    // channelId => channel
    mapping(bytes32 => Channel) public channels;
    // channelId => recipient => authorization
    mapping(bytes32 => mapping(address => ChannelAuthorization)) public channelAuthorizations;

    IERC20 public token;
    uint256 public unlockTime;

    constructor(IERC20 _token, uint256 _unlockTime) {
        token = _token;
        unlockTime = _unlockTime;
    }

    function getChannelId(address creator, bytes32 channelDiscriminator) public pure returns (bytes32 channelId) {
        return keccak256(abi.encode(creator, channelDiscriminator));
    }

    function createAndAuthorize(bytes32 channelDiscriminator, uint256 initialAmount, address recipient, uint16 permissions, uint256 reservedAmount, uint256 limit) public {
        bytes32 channelId = create(channelDiscriminator, initialAmount);
        authorize(channelId, recipient, permissions, reservedAmount, limit);
    }

    // create a new channel
    function create(bytes32 channelDiscriminator, uint256 initialAmount) public returns (bytes32 channelId){
        channelId = getChannelId(msg.sender, channelDiscriminator);
        Channel storage channel = channels[channelId];
        ChannelAuthorization storage authorization = channelAuthorizations[channelId][msg.sender];

        if(channel.available != 0) {
            revert AlreadyInitialized();
        }
        channel.available = uint128(initialAmount);
        authorization.permissions = PERMISSION_ADMIN;

        emit Deposit(channelId, msg.sender, initialAmount);
        emit Authorize(channelId, msg.sender, PERMISSION_ADMIN, 0, 0);

        if (initialAmount > 0) {
            token.safeTransferFrom(msg.sender, address(this), initialAmount);
        }
    }

    // add more funds to a channel (does not check authorizations)
    function deposit(bytes32 channelId, uint256 amount) public {
        //if (amount == 0) revert AmountRequired();
        Channel storage channel = channels[channelId];

        channel.available = channel.available + uint128(amount);

        emit Deposit(channelId, msg.sender, amount);

        token.safeTransferFrom(msg.sender, address(this), amount);
    }

    // authorize a recipient for a channel
    function authorize(bytes32 channelId, address recipient, uint16 permissions, uint256 reservedAmount, uint256 extraLimit) public {
        Channel storage channel = channels[channelId];
        ChannelAuthorization storage ownAuthorization = channelAuthorizations[channelId][msg.sender];
        ChannelAuthorization storage authorization = channelAuthorizations[channelId][recipient];

        if (ownAuthorization.permissions & PERMISSION_ADMIN == PERMISSION_ADMIN) {
            // Allow
        } else if (ownAuthorization.permissions & PERMISSION_MANAGE_AUTHORIZATIONS == PERMISSION_MANAGE_AUTHORIZATIONS) {
            uint16 allowedPermissions = ownAuthorization.permissions & ~PERMISSION_MANAGE_AUTHORIZATIONS;
            if (permissions & allowedPermissions != permissions) {
                revert NotAuthorized(); // Cannot add additional permissions to sub-authorized recipients created
            }

            if (ownAuthorization.permissions & PERMISSION_NO_LIMIT == PERMISSION_NO_LIMIT) {
                // Allow regardless of limit
            } else {
                if (ownAuthorization.limit < extraLimit) {
                    revert InsufficientFunds();
                }
                ownAuthorization.limit -= uint128(extraLimit); // Reduce own limit, transferring it to the recipient
                emit Authorize(channelId, msg.sender, ownAuthorization.permissions, ownAuthorization.reservation, ownAuthorization.limit);
            }
        } else {
            revert NotAuthorized(); // Default deny
        }

        if (authorization.permissions != 0) {
            if (authorization.unlocksAt == 0 || block.timestamp < authorization.unlocksAt) {
                revert AuthorizationLocked();
            }
        }

        channel.reserved = channel.reserved + uint128(reservedAmount) - authorization.reservation;
        authorization.permissions = permissions;
        authorization.limit += uint128(extraLimit);
        authorization.reservation = uint128(reservedAmount);

        emit Authorize(channelId, recipient, permissions, reservedAmount, authorization.limit);
    }

    // start the unlocking period for a channel authorization
    function unlock(bytes32 channelId, address recipient) public {
        //Channel storage channel = channels[channelId];
        ChannelAuthorization storage ownAuthorization = channelAuthorizations[channelId][msg.sender];
        ChannelAuthorization storage authorization = channelAuthorizations[channelId][recipient];

         if (ownAuthorization.permissions & PERMISSION_ADMIN == PERMISSION_ADMIN) {
            // Allow
        } else if (ownAuthorization.permissions & PERMISSION_MANAGE_AUTHORIZATIONS == PERMISSION_MANAGE_AUTHORIZATIONS) {
            uint256 allowedSubPermissions = ownAuthorization.permissions & ~PERMISSION_MANAGE_AUTHORIZATIONS;
            if (authorization.permissions & allowedSubPermissions != authorization.permissions) {
                revert NotAuthorized(); // Cannot unlock recipients with more permissions than us
            }
            // Allow
        } else {
            revert NotAuthorized(); // Default deny
        }

        authorization.unlocksAt = uint240(block.timestamp + unlockTime);

        emit Unlock(channelId, recipient);
    }

    // withdraw money from a channel
    function withdraw(bytes32 channelId, address transferAddress, uint256 transferAmount) public {
        //if (amount == 0) revert AmountRequired();
        if (transferAddress == address(0)) transferAddress = msg.sender;

        Channel storage channel = channels[channelId];
        ChannelAuthorization storage ownAuthorization = channelAuthorizations[channelId][msg.sender];

        if (ownAuthorization.permissions & PERMISSION_ADMIN == PERMISSION_ADMIN) {
            // Allow
        } else if (ownAuthorization.permissions & PERMISSION_WITHDRAW == PERMISSION_WITHDRAW) {
            if (ownAuthorization.permissions & PERMISSION_NO_LIMIT == PERMISSION_NO_LIMIT) {
                // Allow regardless of limit
            } else {
                if (ownAuthorization.limit < transferAmount) {
                    revert InsufficientFunds();
                }
                ownAuthorization.limit -= uint128(transferAmount); // Reduce limit and allow
                emit Authorize(channelId, msg.sender, ownAuthorization.permissions, ownAuthorization.reservation, ownAuthorization.limit);
            }
        } else {
            revert NotAuthorized(); // Default deny
        }

        if (channel.available < transferAmount) {
            revert InsufficientFunds();
        }

        channel.available -= uint128(transferAmount);

        if (channel.available < channel.reserved) { // Reduce this recipient's reservation in order to keep the channel's total reservation fit the available amount; if the recipient is not happy with the new reservation amount, they should gracefully abandon work, as this means we are close to exhausting the funds of the channel
            uint256 reservationReduction = (channel.reserved - channel.available);

            if (ownAuthorization.reservation < reservationReduction) {
                revert InsufficientFunds();
            }

            ownAuthorization.reservation -= uint128(reservationReduction);
            channel.reserved -= uint128(reservationReduction);

            emit Authorize(channelId, msg.sender, ownAuthorization.permissions, ownAuthorization.reservation, ownAuthorization.limit);
        }

        emit Withdraw(channelId, msg.sender, transferAmount);

        token.safeTransfer(transferAddress, transferAmount);
    }

}
