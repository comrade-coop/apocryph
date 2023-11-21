// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {ERC20} from "../lib/openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";

contract Registry {
    // for searching pricing tables by prices offchain
    event NewPricingTable(
        uint256 indexed Id,
        uint256 indexed CpuPrice,
        uint256 indexed RamPrice,
        uint256 StoragePrice,
        uint256 BandwidthEgressPrice,
        uint256 BandwidthIngressPrice
    );
    event Subscribed(uint256 indexed id, address indexed provider);
    event UnSubscribed(uint256 indexed id, address indexed provider);

    mapping(address => mapping(uint256 => PricingTable)) public pricingTables; // IERC20Token => PricingTableId => PricingTable
    mapping(address => Provider) providers; // providerAddress => providerInfo

    struct Provider {
        string name;
        string[] multiAddresses;
        string[] regions;
    }
    // TODO string[] attestationDetails;

    uint256 public pricingTableId;

    enum ResourceKind {
        CPU,
        RAM,
        STORAGE,
        BANDWIDTH_EGRESS,
        BANDWIDTH_INGRESS
    }

    uint8 constant NUM_RESOURCES = 5;

    struct PricingTable {
        mapping(ResourceKind => uint256) resources; // RessourceKind => price
        mapping(address => bool) subscriptions; // providerAddress => providerInfo
        address[] providers;
        string Cpumodel; // intel, amd, arm, ...etc
        string TeeType; // Secure Enclaves, CVM, ...etc
		// string GpuModel;
    }

    modifier registered() {
        Provider storage provider = providers[msg.sender];
        require(bytes(provider.name).length > 0, "Provider Must be registered");
        _;
    }

    modifier subscription(address token, uint256 tableId, bool condition) {
        PricingTable storage pricingTable = pricingTables[token][tableId];
        string memory errorMsg;
        if (condition) {
            errorMsg = "Provider Already UnSubscribed";
        } else {
            errorMsg = "Provider Already Subscribed";
        }
        require(pricingTable.subscriptions[msg.sender] == condition, errorMsg);
        _;
    }

    function registerProvider(string memory name, string[] memory regions, string[] memory multiAddresses) public {
        require(regions.length > 0, "Regions array cannot be empty");
        require(multiAddresses.length > 0, "MultiAddresses array cannot be empty");
        require(bytes(name).length > 0, "Name cannot be empty");

        // TODO verify provider attestation

        Provider storage provider = providers[msg.sender];
        provider.name = name;
        provider.multiAddresses = multiAddresses;
        provider.regions = regions;
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

        PricingTable storage pricingTable = pricingTables[token][pricingTableId + 1];
        pricingTable.resources[ResourceKind.CPU] = CpuPrice;
        pricingTable.resources[ResourceKind.RAM] = RamPrice;
        pricingTable.resources[ResourceKind.STORAGE] = StoragePrice;
        pricingTable.resources[ResourceKind.BANDWIDTH_EGRESS] = BandwidthEgressPrice;
        pricingTable.resources[ResourceKind.BANDWIDTH_INGRESS] = BandwidthIngressPrice;
        pricingTable.Cpumodel = cpumodel;
        pricingTable.TeeType = teeType;

        pricingTable.providers.push(msg.sender);
        pricingTable.subscriptions[msg.sender] = true;

        pricingTableId++;

        emit NewPricingTable(pricingTableId, CpuPrice, RamPrice, StoragePrice, BandwidthEgressPrice, BandwidthIngressPrice);
    }

    function subscribe(address token, uint256 tableId)
        public
        subscription(token, tableId, false)
        registered
        subscription(token, tableId, false)
    {
        PricingTable storage pricingTable = pricingTables[token][tableId];
        pricingTable.providers.push(msg.sender);
        pricingTable.subscriptions[msg.sender] = true;
        emit Subscribed(tableId, msg.sender);
    }

    function unsubscribe(address token, uint256 tableId) public registered subscription(token, tableId, true) {
        PricingTable storage pricingTable = pricingTables[token][tableId];
        for (uint256 index = 0; index < pricingTable.providers.length; index++) {
            if (pricingTable.providers[index] == msg.sender) {
                pricingTable.providers[index] = pricingTable.providers[pricingTable.providers.length - 1];
                pricingTable.providers.pop();
                emit UnSubscribed(tableId, msg.sender);
            }
        }
        pricingTable.subscriptions[msg.sender] = false;
    }

    function updateRegion(string[] memory regions) public registered {
        require(regions.length > 0, "Region cannot be empty");
        Provider storage provider = providers[msg.sender];
        provider.regions = regions;
    }

    function updateMultiAddr(string[] memory addresses) public registered {
        require(addresses.length > 0, "Addresses array cannot be empty");
        Provider storage provider = providers[msg.sender];
        provider.multiAddresses = addresses;
    }

    function getCpuPrice(address token, uint256 id) external view returns (uint256 price) {
        return pricingTables[token][id].resources[ResourceKind.CPU];
    }

    function getRamPrice(address token, uint256 id) external view returns (uint256 price) {
        return pricingTables[token][id].resources[ResourceKind.RAM];
    }

    function getStoragePrice(address token, uint256 id) external view returns (uint256 price) {
        return pricingTables[token][id].resources[ResourceKind.STORAGE];
    }

    function getBandwidthEgressPrice(address token, uint256 id) external view returns (uint256 price) {
        return pricingTables[token][id].resources[ResourceKind.BANDWIDTH_EGRESS];
    }

    function getBandwidthIngressPrice(address token, uint256 id) external view returns (uint256 price) {
        return pricingTables[token][id].resources[ResourceKind.BANDWIDTH_INGRESS];
    }

    function getCpuModel(address token, uint256 id) external view returns (string memory model) {
        return pricingTables[token][id].Cpumodel;
    }

    function getTeeType(address token, uint256 id) external view returns (string memory tee) {
        return pricingTables[token][id].TeeType;
    }

    function getProviders(address token, uint256 id) external view returns (address[] memory) {
        return pricingTables[token][id].providers;
    }

    function isSubscribed(address token, uint256 id) external view returns (bool) {
        return pricingTables[token][id].subscriptions[msg.sender];
    }
}
