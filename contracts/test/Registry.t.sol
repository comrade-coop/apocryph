// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.22;

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
    string cid = "QmWjoiG9mSxprZ18eKam9SsMrxehNNRXeBz6ehYj3wy74y";
    uint256 tableId = 1;
    uint256[5] prices;

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
    }

    function test_register() public {
        registry.registerProvider(cid);
        vm.expectEmit(true, true, false, true, address(registry));
        // The event we expect
        emit Registry.NewPricingTable(
            address(token),
            tableId,
            CPU_PRICE,
            RAM_PRICE,
            STORAGE_PRICE,
            BANDWIDTH_EPRICE,
            BANDWIDTH_INPRICE,
            Cpumodel,
            TeeType
        );

        vm.expectEmit(true, true, false, false, address(registry));
        emit Registry.Subscribed(tableId, provider);

        registry.registerPricingTable(address(token), prices, Cpumodel, TeeType);
    }

    function test_subscribe() public {
        registry.registerProvider(cid);
        registry.registerPricingTable(address(token), prices, Cpumodel, TeeType);
        vm.startPrank(provider2);
        registry.registerProvider(cid);
        vm.expectEmit(true, true, false, false, address(registry));
        emit Registry.Subscribed(tableId, provider2);
        registry.subscribe(tableId);
        bool subscribed = registry.isSubscribed(tableId);
        assertEq(subscribed, true);
    }

    function test_unsubscribe() public {
        vm.expectRevert(); // must be registered
        registry.subscribe(tableId);
        registry.registerProvider(cid);
        registry.subscribe(tableId);
        bool subscribed = registry.isSubscribed(tableId);
        assertEq(subscribed, true);
        registry.unsubscribe(tableId);
        subscribed = registry.isSubscribed(tableId);
        assertEq(subscribed, false);
    }
}
