// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/IUMiCertRegistry.sol";

contract IUMiCertRegistryTest is Test {
    IUMiCertRegistry public registry;
    address public owner;
    address public attacker;

    bytes32 public root1 = keccak256("root_semester_1_2023_v1");
    bytes32 public root2 = keccak256("root_semester_1_2023_v2");
    bytes32 public root3 = keccak256("root_semester_2_2023_v1");

    function setUp() public {
        owner = address(this);
        attacker = address(0x1234);
        registry = new IUMiCertRegistry(owner);
    }

    function testPublishTermRoot() public {
        registry.publishTermRoot(root1, "Semester_1_2023", 100);

        (bytes32 latestRoot, uint256 version, uint256 totalStudents, uint256 publishedAt) =
            registry.getLatestRoot("Semester_1_2023");

        assertEq(latestRoot, root1);
        assertEq(version, 1);
        assertEq(totalStudents, 100);
        assertTrue(publishedAt > 0);
    }

    function testPublishMultipleTerms() public {
        registry.publishTermRoot(root1, "Semester_1_2023", 100);
        registry.publishTermRoot(root3, "Semester_2_2023", 120);

        (bytes32 latestRoot1, uint256 version1,,) = registry.getLatestRoot("Semester_1_2023");
        (bytes32 latestRoot2, uint256 version2,,) = registry.getLatestRoot("Semester_2_2023");

        assertEq(latestRoot1, root1);
        assertEq(version1, 1);
        assertEq(latestRoot2, root3);
        assertEq(version2, 1);
    }

    function testSupersedeTerm() public {
        // Publish initial version
        registry.publishTermRoot(root1, "Semester_1_2023", 100);

        // Supersede with new version (simulating revocation)
        registry.supersedeTerm(
            "Semester_1_2023",
            root2,
            99,
            "Revoked: Student ID 00001 - fraudulent credential"
        );

        // Check latest version
        (bytes32 latestRoot, uint256 version, uint256 totalStudents,) =
            registry.getLatestRoot("Semester_1_2023");

        assertEq(latestRoot, root2);
        assertEq(version, 2);
        assertEq(totalStudents, 99);

        // Check old version is marked as superseded
        (
            bytes32 oldRootHash,
            ,
            ,
            bool isSuperseded,
            bytes32 supersededBy,
            string memory reason
        ) = registry.getVersionInfo("Semester_1_2023", 1);

        assertEq(oldRootHash, root1);
        assertTrue(isSuperseded);
        assertEq(supersededBy, root2);
        assertEq(reason, "Revoked: Student ID 00001 - fraudulent credential");
    }

    function testCheckRootStatus_Current() public {
        registry.publishTermRoot(root1, "Semester_1_2023", 100);

        (
            uint8 status,
            string memory termId,
            uint256 version,
            bytes32 latestRoot,
            string memory message
        ) = registry.checkRootStatus(root1);

        assertEq(status, 1); // Valid & Current
        assertEq(termId, "Semester_1_2023");
        assertEq(version, 1);
        assertEq(latestRoot, root1);
        assertEq(message, "Valid - Current version");
    }

    function testCheckRootStatus_Superseded() public {
        registry.publishTermRoot(root1, "Semester_1_2023", 100);
        registry.supersedeTerm("Semester_1_2023", root2, 99, "Credential revocation");

        (
            uint8 status,
            string memory termId,
            uint256 version,
            bytes32 latestRoot,
            string memory message
        ) = registry.checkRootStatus(root1);

        assertEq(status, 3); // Superseded
        assertEq(termId, "Semester_1_2023");
        assertEq(version, 1);
        assertEq(latestRoot, root2);
        assertTrue(bytes(message).length > 0);
    }

    function testCheckRootStatus_Invalid() public {
        bytes32 nonExistentRoot = keccak256("nonexistent");

        (
            uint8 status,
            string memory termId,
            uint256 version,
            bytes32 latestRoot,
            string memory message
        ) = registry.checkRootStatus(nonExistentRoot);

        assertEq(status, 0); // Invalid
        assertEq(termId, "");
        assertEq(version, 0);
        assertEq(latestRoot, bytes32(0));
        assertEq(message, "Root not found in registry");
    }

    function testGetTermHistory() public {
        registry.publishTermRoot(root1, "Semester_1_2023", 100);
        registry.supersedeTerm("Semester_1_2023", root2, 99, "First revocation");

        (uint256[] memory versions, bytes32[] memory roots) =
            registry.getTermHistory("Semester_1_2023");

        assertEq(versions.length, 2);
        assertEq(roots.length, 2);
        assertEq(versions[0], 1);
        assertEq(versions[1], 2);
        assertEq(roots[0], root1);
        assertEq(roots[1], root2);
    }

    function testVerifyReceiptAnchor_Valid() public {
        registry.publishTermRoot(root1, "Semester_1_2023", 100);

        (bool isValid, string memory termId, uint256 publishedAt) =
            registry.verifyReceiptAnchor(root1);

        assertTrue(isValid);
        assertEq(termId, "Semester_1_2023");
        assertTrue(publishedAt > 0);
    }

    function testVerifyReceiptAnchor_Superseded() public {
        registry.publishTermRoot(root1, "Semester_1_2023", 100);
        registry.supersedeTerm("Semester_1_2023", root2, 99, "Revocation");

        (bool isValid, string memory termId, uint256 publishedAt) =
            registry.verifyReceiptAnchor(root1);

        assertFalse(isValid); // Not valid because superseded
        assertEq(termId, "Semester_1_2023");
        assertTrue(publishedAt > 0);
    }

    function test_RevertWhen_PublishDuplicateRoot() public {
        registry.publishTermRoot(root1, "Semester_1_2023", 100);

        vm.expectRevert("Root already published");
        registry.publishTermRoot(root1, "Semester_2_2023", 120);
    }

    function test_RevertWhen_SupersedeNonExistentTerm() public {
        vm.expectRevert("Term not found");
        registry.supersedeTerm("NonExistentTerm", root1, 100, "Test");
    }

    function test_RevertWhen_UnauthorizedPublish() public {
        vm.prank(attacker);
        vm.expectRevert();
        registry.publishTermRoot(root1, "Semester_1_2023", 100);
    }

    function test_RevertWhen_UnauthorizedSupersede() public {
        registry.publishTermRoot(root1, "Semester_1_2023", 100);

        vm.prank(attacker);
        vm.expectRevert();
        registry.supersedeTerm("Semester_1_2023", root2, 99, "Unauthorized");
    }

    function testMultipleVersions() public {
        // Publish version 1
        registry.publishTermRoot(root1, "Semester_1_2023", 100);

        // Supersede to version 2
        bytes32 root2a = keccak256("root_v2");
        registry.supersedeTerm("Semester_1_2023", root2a, 99, "First revocation");

        // Supersede to version 3
        bytes32 root3a = keccak256("root_v3");
        registry.supersedeTerm("Semester_1_2023", root3a, 98, "Second revocation");

        // Check latest
        (bytes32 latestRoot, uint256 version,,) = registry.getLatestRoot("Semester_1_2023");
        assertEq(latestRoot, root3a);
        assertEq(version, 3);

        // Check all versions exist
        (uint256[] memory versions, bytes32[] memory roots) =
            registry.getTermHistory("Semester_1_2023");

        assertEq(versions.length, 3);
        assertEq(roots[0], root1);
        assertEq(roots[1], root2a);
        assertEq(roots[2], root3a);
    }

    function testEventsEmitted() public {
        // Test TermRootPublished event
        vm.expectEmit(true, true, false, true);
        emit IUMiCertRegistry.TermRootPublished("Semester_1_2023", root1, 1, 100, block.timestamp);
        registry.publishTermRoot(root1, "Semester_1_2023", 100);

        // Test TermRootSuperseded event
        vm.expectEmit(true, true, true, false);
        emit IUMiCertRegistry.TermRootSuperseded(
            "Semester_1_2023",
            1,
            root1,
            2,
            root2,
            "Test revocation"
        );
        registry.supersedeTerm("Semester_1_2023", root2, 99, "Test revocation");
    }
}
