// SPDX-License-Identifier: GPL-3.0

pragma solidity ^0.8.22;

import {ERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";

contract Registry {
    // for searching pricing tables by prices offchain
    event NewPricingTable(
        address indexed token,
        uint256 indexed Id,
        uint256 CpuPrice,
        uint256 RamPrice,
        uint256 StoragePrice,
        uint256 BandwidthEgressPrice,
        uint256 BandwidthIngressPrice,
        string Cpumodel,
        string TeeType
    );
    event Subscribed(uint256 indexed id, address indexed provider);

    mapping(address => string) public providers; // providerAddress => provider profile cid
    mapping(address => mapping(uint256 => bool)) public subscription; // provider => tableId => subscribed/unsubscribed

    uint256 public pricingTableId;

    enum ResourceKind {
        CPU,
        RAM,
        STORAGE,
        BANDWIDTH_EGRESS,
        BANDWIDTH_INGRESS
    }

    uint8 constant NUM_RESOURCES = 5;

    modifier registered() {
        string memory cid = providers[msg.sender];
        require(bytes(cid).length > 0, "Provider Must be registered");
        _;
    }

    // could also be used to update the cid
    function registerProvider(string calldata cid) public {
        require(bytes(cid).length > 0, "cid must not be empty");
        // TODO verify provider attestation
        providers[msg.sender] = cid;
    }

    function registerPricingTable(
        address token,
        uint256[NUM_RESOURCES] memory Prices,
        string memory cpumodel,
        string memory teeType
    ) public registered {
        uint256 CpuPrice = Prices[uint256(ResourceKind.CPU)];
        uint256 RamPrice = Prices[uint256(ResourceKind.RAM)];
        uint256 StoragePrice = Prices[uint256(ResourceKind.STORAGE)];
        uint256 BandwidthEgressPrice = Prices[uint256(ResourceKind.BANDWIDTH_EGRESS)];
        uint256 BandwidthIngressPrice = Prices[uint256(ResourceKind.BANDWIDTH_INGRESS)];

        pricingTableId++;

        emit NewPricingTable(
            token,
            pricingTableId,
            CpuPrice,
            RamPrice,
            StoragePrice,
            BandwidthEgressPrice,
            BandwidthIngressPrice,
            cpumodel,
            teeType
        );
        subscription[msg.sender][pricingTableId] = true;
        emit Subscribed(pricingTableId, msg.sender);
    }

    function subscribe(uint256 tableId) public registered {
        subscription[msg.sender][tableId] = true;
        emit Subscribed(tableId, msg.sender);
    }

    function unsubscribe(uint256 tableId) public registered {
        subscription[msg.sender][tableId] = false;
    }

    function isSubscribed(address provider, uint256 tableId) external view returns (bool) {
        return subscription[provider][tableId];
    }
}
