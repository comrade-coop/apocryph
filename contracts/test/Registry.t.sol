// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

import {Test, console2} from "../lib/forge-std/src/Test.sol";
import {Registry} from "../src/Registry.sol";
import {MockToken} from "../src/MockToken.sol";
import {IERC20Errors} from "../lib/openzeppelin-contracts/contracts/interfaces/draft-IERC6093.sol";

contract RegistryTest is Test {
    Registry public registry;
    MockToken public token;
    address provider;
    address provider2;
    uint256 CPU_PRICE = 10;
    uint256 RAM_PRICE = 20;
    uint256 STORAGE_PRICE = 30;
    uint256 BANDWIDTH_EPRICE = 40;
    uint256 BANDWIDTH_INPRICE = 50;
    string Cpumodel = "Intel Xeon Platinum 8452Y Processor";
    string TeeType = "Secure Enclave";
    uint256 tableId = 1;
    uint256[5] prices;
    string[] regions;
    string[] multiAddresses;

    function setUp() public {
        provider = vm.createWallet("provider").addr;
        provider2 = vm.createWallet("provider2").addr;
        registry = new Registry();
        token = new MockToken();
        vm.startPrank(provider);
        prices[0] = CPU_PRICE;
        prices[1] = RAM_PRICE;
        prices[2] = STORAGE_PRICE;
        prices[3] = BANDWIDTH_EPRICE;
        prices[4] = BANDWIDTH_INPRICE;

        regions = new string[](2);
        regions[0] = "bul-east-1";
        regions[1] = "alg-west-2";

        multiAddresses = new string[](3);
        multiAddresses[0] = "/p2p/12D3KooWPcMp99mZkfdk8qfHrfjFzQbZcjxSf14egDQtq5zWFhWp"; // for ipfs libp2p
        multiAddresses[1] = "/ip4/127.0.0.1/tcp/4001/p2p/QmWjoiG9mSxprZ18eKam9SsMrxehNNRXeBz6ehYj3wy74y"; // libp2p request response
        multiAddresses[2] = "https://kubo.business/"; // dns

        registry.registerProvider("kubo-Cloud", regions, multiAddresses);

        registry.registerPricingTable(address(token), prices, Cpumodel, TeeType);
    }

    function test_register() public {
        uint256 CpuPrice = registry.getCpuPrice(address(token), tableId);
        uint256 RamPrice = registry.getRamPrice(address(token), tableId);
        uint256 StoragePrice = registry.getStoragePrice(address(token), tableId);
        uint256 BandiwidthEgressPrice = registry.getBandwidthEgressPrice(address(token), tableId);
        uint256 BandiwidthIngressPrice = registry.getBandwidthIngressPrice(address(token), tableId);

        string memory model = registry.getCpuModel(address(token), tableId);
        string memory teeType = registry.getTeeType(address(token), tableId);

        address[] memory providers = registry.getProviders(address(token), tableId);
        assertEq(providers.length, 1);
        assertEq(providers[0], provider);
        assertEq(CpuPrice, CPU_PRICE);
        assertEq(RamPrice, RAM_PRICE);
        assertEq(StoragePrice, STORAGE_PRICE);
        assertEq(BandiwidthEgressPrice, BANDWIDTH_EPRICE);
        assertEq(BandiwidthIngressPrice, BANDWIDTH_INPRICE);
        assertEq(model, Cpumodel);
        assertEq(teeType, TeeType);

        registry.registerPricingTable(address(token), prices, Cpumodel, TeeType);
        providers = registry.getProviders(address(token), tableId + 1);
        assertEq(providers.length, 1);
        registry.registerPricingTable(address(token), prices, Cpumodel, TeeType);
        assertEq(registry.pricingTableId(), 3);
    }

    function test_subscribe() public {
        vm.startPrank(provider2);
        vm.expectRevert(); // must be registered
        registry.subscribe(address(token), tableId);
        registry.registerProvider("buko-cloud", regions, multiAddresses);
        registry.subscribe(address(token), tableId);
        address[] memory providers = registry.getProviders(address(token), tableId);
        assertEq(providers.length, 2);
        assertEq(providers[0], provider);
        assertEq(providers[1], provider2);
        vm.expectRevert(); // already subscribed
        registry.subscribe(address(token), tableId);
		bool subscribed = registry.isSubscribed(address(token), tableId);
		assertEq(subscribed, true);
    }

	function test_unsubscribe() public {
        vm.startPrank(provider2);
        vm.expectRevert(); // must be registered
        registry.unsubscribe(address(token), tableId);
        vm.startPrank(provider);
        registry.unsubscribe(address(token), tableId);
        address[] memory providers = registry.getProviders(address(token), tableId);
        assertEq(providers.length, 0);
        vm.expectRevert(); // already unsubscribed
        registry.unsubscribe(address(token), tableId);
		bool subscribed = registry.isSubscribed(address(token), tableId);
		assertEq(subscribed, false);
	}
}
