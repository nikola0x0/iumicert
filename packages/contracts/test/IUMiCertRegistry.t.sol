// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Test} from "forge-std/Test.sol";
import {IUMiCertRegistry} from "../src/IUMiCertRegistry.sol";

contract IUMiCertRegistryTest is Test {
    IUMiCertRegistry public registry;
    address public owner;
    address public student;
    
    // Test data from actual IU-MiCert system
    bytes32 public constant VERKLE_ROOT_S1_2023 = 0x2b5a353f110276a123fa7fcf6a5a951e5944a0a30cef2df878534f7437007bc9;
    bytes32 public constant VERKLE_ROOT_S2_2023 = 0x57aca62e24eaac0892acaf7f1fcaff3b78a43602934f13a8d79af9c5a2692eb4;
    
    string public constant TERM_S1_2023 = "Semester_1_2023";
    string public constant TERM_S2_2023 = "Semester_2_2023";

    function setUp() public {
        owner = makeAddr("owner");
        student = makeAddr("student");
        
        vm.prank(owner);
        registry = new IUMiCertRegistry(owner);
    }

    function testInitialState() public {
        assertEq(registry.owner(), owner);
        assertEq(registry.getPublishedRootsCount(), 0);
        
        (,,,bool exists) = registry.getTermRootInfo(VERKLE_ROOT_S1_2023);
        assertFalse(exists);
    }

    function testPublishTermRoot() public {
        vm.prank(owner);
        registry.publishTermRoot(VERKLE_ROOT_S1_2023, TERM_S1_2023, 5);
        
        (string memory termId, uint256 totalStudents, uint256 publishedAt, bool exists) = 
            registry.getTermRootInfo(VERKLE_ROOT_S1_2023);
        
        assertTrue(exists);
        assertEq(termId, TERM_S1_2023);
        assertEq(totalStudents, 5);
        assertGt(publishedAt, 0);
        assertEq(registry.getPublishedRootsCount(), 1);
    }

    function testPublishTermRootEmitsEvent() public {
        vm.expectEmit(true, true, false, false);
        emit IUMiCertRegistry.TermRootPublished(VERKLE_ROOT_S1_2023, TERM_S1_2023, 5, block.timestamp);
        
        vm.prank(owner);
        registry.publishTermRoot(VERKLE_ROOT_S1_2023, TERM_S1_2023, 5);
    }

    function testVerifyReceiptAnchor() public {
        // Publish a term root
        vm.prank(owner);
        registry.publishTermRoot(VERKLE_ROOT_S1_2023, TERM_S1_2023, 5);
        
        // Verify the receipt anchor (this is what students do)
        (bool isValid, string memory termId, uint256 publishedAt) = 
            registry.verifyReceiptAnchor(VERKLE_ROOT_S1_2023);
        
        assertTrue(isValid);
        assertEq(termId, TERM_S1_2023);
        assertGt(publishedAt, 0);
    }

    function testVerifyInvalidReceiptAnchor() public {
        bytes32 invalidRoot = keccak256("invalid");
        
        (bool isValid, string memory termId, uint256 publishedAt) = 
            registry.verifyReceiptAnchor(invalidRoot);
        
        assertFalse(isValid);
        assertEq(termId, "");
        assertEq(publishedAt, 0);
    }

    function testCannotPublishSameRootTwice() public {
        vm.prank(owner);
        registry.publishTermRoot(VERKLE_ROOT_S1_2023, TERM_S1_2023, 5);
        
        vm.expectRevert("Root already published");
        vm.prank(owner);
        registry.publishTermRoot(VERKLE_ROOT_S1_2023, TERM_S1_2023, 5);
    }

    function testOnlyOwnerCanPublish() public {
        vm.expectRevert();
        vm.prank(student);
        registry.publishTermRoot(VERKLE_ROOT_S1_2023, TERM_S1_2023, 5);
    }

    function testPublishMultipleTerms() public {
        vm.prank(owner);
        registry.publishTermRoot(VERKLE_ROOT_S1_2023, TERM_S1_2023, 5);
        
        vm.prank(owner);
        registry.publishTermRoot(VERKLE_ROOT_S2_2023, TERM_S2_2023, 5);
        
        assertEq(registry.getPublishedRootsCount(), 2);
        assertEq(registry.getPublishedRoot(0), VERKLE_ROOT_S1_2023);
        assertEq(registry.getPublishedRoot(1), VERKLE_ROOT_S2_2023);
    }

    function testInvalidInputs() public {
        vm.expectRevert("Invalid Verkle root");
        vm.prank(owner);
        registry.publishTermRoot(bytes32(0), TERM_S1_2023, 5);
        
        vm.expectRevert("Term ID required");
        vm.prank(owner);
        registry.publishTermRoot(VERKLE_ROOT_S1_2023, "", 5);
        
        vm.expectRevert("Invalid student count");
        vm.prank(owner);
        registry.publishTermRoot(VERKLE_ROOT_S1_2023, TERM_S1_2023, 0);
    }
}