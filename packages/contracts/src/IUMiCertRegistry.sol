// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title IUMiCertRegistry
 * @author IU-MiCert Team
 * @notice This contract stores term Verkle roots for IU Vietnam's academic credentials.
 * Students can use these published roots to verify their academic journey receipts.
 * Verifiers can check that a student's receipt contains valid blockchain-anchored data.
 */
contract IUMiCertRegistry is Ownable {
    
    struct TermRoot {
        string termId;           // e.g., "Semester_1_2023"
        uint256 totalStudents;   // Number of students in this term
        uint256 publishedAt;     // When this root was published
        bool exists;             // Whether this term root exists
    }

    // Mapping from Verkle root hash to term information
    mapping(bytes32 => TermRoot) public termRoots;
    
    // Array to enumerate all published roots
    bytes32[] public publishedRoots;
    
    // Events
    event TermRootPublished(
        bytes32 indexed verkleRoot,
        string indexed termId,
        uint256 totalStudents,
        uint256 timestamp
    );

    constructor(address initialOwner) Ownable(initialOwner) {}

    /**
     * @notice Publishes a term's Verkle root commitment
     * @dev Only IU Vietnam (owner) can publish term roots
     * @param _verkleRoot The Verkle tree root for this academic term
     * @param _termId Human-readable term identifier
     * @param _totalStudents Number of students included in this term
     */
    function publishTermRoot(
        bytes32 _verkleRoot, 
        string memory _termId,
        uint256 _totalStudents
    ) external onlyOwner {
        require(_verkleRoot != bytes32(0), "Invalid Verkle root");
        require(bytes(_termId).length > 0, "Term ID required");
        require(_totalStudents > 0, "Invalid student count");
        require(!termRoots[_verkleRoot].exists, "Root already published");

        termRoots[_verkleRoot] = TermRoot({
            termId: _termId,
            totalStudents: _totalStudents,
            publishedAt: block.timestamp,
            exists: true
        });

        publishedRoots.push(_verkleRoot);

        emit TermRootPublished(_verkleRoot, _termId, _totalStudents, block.timestamp);
    }

    /**
     * @notice Verifies if a receipt's blockchain anchor is valid
     * @dev Students include blockchain_anchor in their receipts for verification
     * @param _blockchainAnchor The Verkle root from student's receipt
     * @return isValid Whether this anchor exists on blockchain
     * @return termId The term this anchor belongs to
     * @return publishedAt When this term was published
     */
    function verifyReceiptAnchor(bytes32 _blockchainAnchor) 
        external 
        view 
        returns (bool isValid, string memory termId, uint256 publishedAt) 
    {
        TermRoot memory termRoot = termRoots[_blockchainAnchor];
        return (termRoot.exists, termRoot.termId, termRoot.publishedAt);
    }

    /**
     * @notice Gets information about a published term root
     */
    function getTermRootInfo(bytes32 _verkleRoot) 
        external 
        view 
        returns (
            string memory termId,
            uint256 totalStudents, 
            uint256 publishedAt,
            bool exists
        ) 
    {
        TermRoot memory termRoot = termRoots[_verkleRoot];
        return (termRoot.termId, termRoot.totalStudents, termRoot.publishedAt, termRoot.exists);
    }

    /**
     * @notice Returns total number of published term roots
     */
    function getPublishedRootsCount() external view returns (uint256) {
        return publishedRoots.length;
    }

    /**
     * @notice Gets a published root by index
     */
    function getPublishedRoot(uint256 _index) external view returns (bytes32) {
        require(_index < publishedRoots.length, "Index out of bounds");
        return publishedRoots[_index];
    }
}