// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Script, console2} from "forge-std/Script.sol";
import {IUMiCertRegistry} from "../src/IUMiCertRegistry.sol";

/**
 * @title Deploy IU-MiCert Contracts
 * @notice This script deploys the essential contracts for the IU-MiCert system.
 */
contract Deploy is Script {
    function run()
        public
        returns (
            IUMiCertRegistry
        )
    {
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        address initialOwner = vm.addr(deployerPrivateKey);

        console2.log("Deploying contracts with owner:", initialOwner);

        IUMiCertRegistry registry = new IUMiCertRegistry(initialOwner);
        console2.log("IUMiCertRegistry deployed to:", address(registry));


        vm.stopBroadcast();

        return (registry);
    }
}