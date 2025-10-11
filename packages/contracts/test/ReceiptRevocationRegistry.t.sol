// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {ReceiptRevocationRegistry} from "../src/ReceiptRevocationRegistry.sol";

contract ReceiptRevocationRegistryTest is Test {
    ReceiptRevocationRegistry public registry;
    address public owner;
    address public notOwner;
    
    bytes32 public constant RECEIPT_ID_1 = keccak256("ITITIU00001_journey");
    bytes32 public constant RECEIPT_ID_2 = keccak256("ITITIU00002_journey");
    string public constant REASON = "Academic misconduct";

    function setUp() public {
        owner = makeAddr("owner");
        notOwner = makeAddr("notOwner");
        
        vm.prank(owner);
        registry = new ReceiptRevocationRegistry(owner);
    }

    function testInitialState() public {
        assertEq(registry.owner(), owner);
        assertEq(registry.getTotalRevocations(), 0);
        
        (bool isRevoked, uint256 revokedAt, string memory reason) = registry.isReceiptRevoked(RECEIPT_ID_1);
        assertFalse(isRevoked);
        assertEq(revokedAt, 0);
        assertEq(reason, "");
    }

    function testRevokeReceipt() public {
        vm.prank(owner);
        registry.revokeReceipt(RECEIPT_ID_1, REASON);
        
        (bool isRevoked, uint256 revokedAt, string memory reason) = registry.isReceiptRevoked(RECEIPT_ID_1);
        assertTrue(isRevoked);
        assertGt(revokedAt, 0);
        assertEq(reason, REASON);
        assertEq(registry.getTotalRevocations(), 1);
    }

    function testRevokeReceiptEmitsEvent() public {
        vm.expectEmit(true, false, false, false);
        emit ReceiptRevocationRegistry.ReceiptRevoked(RECEIPT_ID_1, REASON, block.timestamp);
        
        vm.prank(owner);
        registry.revokeReceipt(RECEIPT_ID_1, REASON);
    }

    function testCannotRevokeSameReceiptTwice() public {
        vm.prank(owner);
        registry.revokeReceipt(RECEIPT_ID_1, REASON);
        
        vm.expectRevert("Receipt already revoked");
        vm.prank(owner);
        registry.revokeReceipt(RECEIPT_ID_1, REASON);
    }

    function testOnlyOwnerCanRevoke() public {
        vm.expectRevert();
        vm.prank(notOwner);
        registry.revokeReceipt(RECEIPT_ID_1, REASON);
    }

    function testCannotRevokeWithoutReason() public {
        vm.expectRevert("Reason required");
        vm.prank(owner);
        registry.revokeReceipt(RECEIPT_ID_1, "");
    }

    function testMultipleRevocations() public {
        vm.prank(owner);
        registry.revokeReceipt(RECEIPT_ID_1, REASON);
        
        vm.prank(owner);
        registry.revokeReceipt(RECEIPT_ID_2, "Grade change");
        
        assertEq(registry.getTotalRevocations(), 2);
        assertEq(registry.getRevokedReceiptId(0), RECEIPT_ID_1);
        assertEq(registry.getRevokedReceiptId(1), RECEIPT_ID_2);
    }
}