// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title IUMiCertRegistry
 * @author IU-MiCert Team
 * @notice Enhanced registry with term root versioning for revocation support.
 *
 * This contract manages Verkle tree roots for academic credentials with version tracking.
 * When credentials need to be revoked, a new version of the term root is published,
 * and the old version is marked as superseded. This integrates naturally with the
 * student workflow of downloading updated receipts after each term.
 *
 * Key Features:
 * - Multiple versions per term (v1, v2, v3...)
 * - Revocation through root replacement (not simple flagging)
 * - Backward compatible verification (old receipts remain verifiable)
 * - Efficient batch revocation at term boundaries
 */
contract IUMiCertRegistry is Ownable {

    struct TermRootVersion {
        bytes32 rootHash;           // The Verkle tree root commitment
        uint256 version;            // Version number (1, 2, 3...)
        uint256 totalStudents;      // Number of students in this version
        uint256 publishedAt;        // Timestamp when published
        bool isSuperseded;          // True if a newer version exists
        bytes32 supersededBy;       // Hash of the root that supersedes this one
        string supersessionReason;  // Why this version was superseded
    }

    // termId => version => TermRootVersion
    mapping(string => mapping(uint256 => TermRootVersion)) public termVersions;

    // termId => latest version number
    mapping(string => uint256) public latestVersion;

    // Fast lookup: rootHash => termId (for verification)
    mapping(bytes32 => string) public rootToTerm;

    // rootHash => version number
    mapping(bytes32 => uint256) public rootToVersion;

    // Track all published roots for enumeration
    bytes32[] public publishedRoots;

    // Track all term IDs
    string[] public publishedTerms;
    mapping(string => bool) private termExists;

    // Events
    event TermRootPublished(
        string indexed termId,
        bytes32 indexed verkleRoot,
        uint256 version,
        uint256 totalStudents,
        uint256 timestamp
    );

    event TermRootSuperseded(
        string indexed termId,
        uint256 oldVersion,
        bytes32 indexed oldRoot,
        uint256 newVersion,
        bytes32 indexed newRoot,
        string reason
    );

    constructor(address initialOwner) Ownable(initialOwner) {}

    /**
     * @notice Publishes a new term root (initial publication or new term)
     * @dev First version of a term is automatically version 1
     * @param _verkleRoot The Verkle tree root commitment
     * @param _termId Term identifier (e.g., "Semester_1_2023")
     * @param _totalStudents Number of students in this term
     */
    function publishTermRoot(
        bytes32 _verkleRoot,
        string memory _termId,
        uint256 _totalStudents
    ) external onlyOwner {
        require(_verkleRoot != bytes32(0), "Invalid Verkle root");
        require(bytes(_termId).length > 0, "Term ID required");
        require(_totalStudents > 0, "Invalid student count");
        require(keccak256(bytes(rootToTerm[_verkleRoot])) == keccak256(bytes("")), "Root already published");

        uint256 version = latestVersion[_termId] + 1;

        // Store version data
        termVersions[_termId][version] = TermRootVersion({
            rootHash: _verkleRoot,
            version: version,
            totalStudents: _totalStudents,
            publishedAt: block.timestamp,
            isSuperseded: false,
            supersededBy: bytes32(0),
            supersessionReason: ""
        });

        // Update latest version
        latestVersion[_termId] = version;

        // Update lookups
        rootToTerm[_verkleRoot] = _termId;
        rootToVersion[_verkleRoot] = version;
        publishedRoots.push(_verkleRoot);

        // Track term if first version
        if (version == 1) {
            publishedTerms.push(_termId);
            termExists[_termId] = true;
        }

        emit TermRootPublished(_termId, _verkleRoot, version, _totalStudents, block.timestamp);
    }

    /**
     * @notice Supersede existing term version and publish new version (for revocation)
     * @dev Used when credentials need to be revoked - rebuilds tree without revoked data
     * @param _termId The term being updated
     * @param _newVerkleRoot New Verkle root (with revoked credentials removed)
     * @param _newTotalStudents Updated student count (may change if students fully removed)
     * @param _reason Reason for supersession (e.g., "Credential revocation: IT089IU")
     */
    function supersedeTerm(
        string memory _termId,
        bytes32 _newVerkleRoot,
        uint256 _newTotalStudents,
        string memory _reason
    ) external onlyOwner {
        require(_newVerkleRoot != bytes32(0), "Invalid new root");
        require(bytes(_reason).length > 0, "Reason required");
        require(keccak256(bytes(rootToTerm[_newVerkleRoot])) == keccak256(bytes("")), "New root already published");

        uint256 currentVersion = latestVersion[_termId];
        require(currentVersion > 0, "Term not found");

        // Mark current version as superseded
        TermRootVersion storage currentRoot = termVersions[_termId][currentVersion];
        require(!currentRoot.isSuperseded, "Current version already superseded");

        currentRoot.isSuperseded = true;
        currentRoot.supersededBy = _newVerkleRoot;
        currentRoot.supersessionReason = _reason;

        // Publish new version
        uint256 newVersion = currentVersion + 1;

        termVersions[_termId][newVersion] = TermRootVersion({
            rootHash: _newVerkleRoot,
            version: newVersion,
            totalStudents: _newTotalStudents,
            publishedAt: block.timestamp,
            isSuperseded: false,
            supersededBy: bytes32(0),
            supersessionReason: ""
        });

        // Update latest version
        latestVersion[_termId] = newVersion;

        // Update lookups
        rootToTerm[_newVerkleRoot] = _termId;
        rootToVersion[_newVerkleRoot] = newVersion;
        publishedRoots.push(_newVerkleRoot);

        emit TermRootSuperseded(
            _termId,
            currentVersion,
            currentRoot.rootHash,
            newVersion,
            _newVerkleRoot,
            _reason
        );

        emit TermRootPublished(_termId, _newVerkleRoot, newVersion, _newTotalStudents, block.timestamp);
    }

    /**
     * @notice Get latest root for a term
     * @param _termId The term identifier
     * @return rootHash Latest Verkle root
     * @return version Latest version number
     * @return totalStudents Number of students in latest version
     * @return publishedAt When latest version was published
     */
    function getLatestRoot(string memory _termId)
        external
        view
        returns (
            bytes32 rootHash,
            uint256 version,
            uint256 totalStudents,
            uint256 publishedAt
        )
    {
        uint256 latest = latestVersion[_termId];
        require(latest > 0, "Term not found");

        TermRootVersion memory termRoot = termVersions[_termId][latest];
        return (termRoot.rootHash, termRoot.version, termRoot.totalStudents, termRoot.publishedAt);
    }

    /**
     * @notice Comprehensive root validity check for verifiers
     * @param _verkleRoot The root hash from a student's receipt
     * @return status 0=Invalid, 1=Valid&Current, 2=Valid&Outdated, 3=Superseded
     * @return termId Which term this root belongs to
     * @return version Version number of this root
     * @return latestRoot Current valid root for this term
     * @return message Human-readable status message
     */
    function checkRootStatus(bytes32 _verkleRoot)
        external
        view
        returns (
            uint8 status,
            string memory termId,
            uint256 version,
            bytes32 latestRoot,
            string memory message
        )
    {
        // Check if root exists
        termId = rootToTerm[_verkleRoot];
        if (keccak256(bytes(termId)) == keccak256(bytes(""))) {
            return (0, "", 0, bytes32(0), "Root not found in registry");
        }

        version = rootToVersion[_verkleRoot];
        uint256 currentVersion = latestVersion[termId];
        TermRootVersion memory rootData = termVersions[termId][version];
        TermRootVersion memory currentData = termVersions[termId][currentVersion];

        latestRoot = currentData.rootHash;

        // Check if this is the current version
        if (version == currentVersion) {
            return (1, termId, version, latestRoot, "Valid - Current version");
        }

        // Check if superseded
        if (rootData.isSuperseded) {
            string memory reason = rootData.supersessionReason;
            return (
                3,
                termId,
                version,
                latestRoot,
                string(abi.encodePacked("Superseded - ", reason))
            );
        }

        // Old version but not explicitly superseded
        return (
            2,
            termId,
            version,
            latestRoot,
            "Valid but outdated - Update recommended"
        );
    }

    /**
     * @notice Get complete version history for a term
     * @param _termId The term to query
     * @return versions Array of all version numbers
     * @return roots Array of corresponding root hashes
     */
    function getTermHistory(string memory _termId)
        external
        view
        returns (uint256[] memory versions, bytes32[] memory roots)
    {
        uint256 totalVersions = latestVersion[_termId];
        require(totalVersions > 0, "Term not found");

        versions = new uint256[](totalVersions);
        roots = new bytes32[](totalVersions);

        for (uint256 i = 1; i <= totalVersions; i++) {
            versions[i-1] = i;
            roots[i-1] = termVersions[_termId][i].rootHash;
        }

        return (versions, roots);
    }

    /**
     * @notice Get detailed information about a specific version
     * @param _termId Term identifier
     * @param _version Version number
     */
    function getVersionInfo(string memory _termId, uint256 _version)
        external
        view
        returns (
            bytes32 rootHash,
            uint256 totalStudents,
            uint256 publishedAt,
            bool isSuperseded,
            bytes32 supersededBy,
            string memory supersessionReason
        )
    {
        require(_version > 0 && _version <= latestVersion[_termId], "Invalid version");

        TermRootVersion memory v = termVersions[_termId][_version];
        return (
            v.rootHash,
            v.totalStudents,
            v.publishedAt,
            v.isSuperseded,
            v.supersededBy,
            v.supersessionReason
        );
    }

    /**
     * @notice Legacy compatibility - verify receipt anchor (backward compatible)
     * @dev Returns basic validity without version details
     * @param _blockchainAnchor The Verkle root from student's receipt
     * @return isValid Whether this anchor exists and is not superseded
     * @return termId The term this anchor belongs to
     * @return publishedAt When this version was published
     */
    function verifyReceiptAnchor(bytes32 _blockchainAnchor)
        external
        view
        returns (bool isValid, string memory termId, uint256 publishedAt)
    {
        termId = rootToTerm[_blockchainAnchor];
        if (keccak256(bytes(termId)) == keccak256(bytes(""))) {
            return (false, "", 0);
        }

        uint256 version = rootToVersion[_blockchainAnchor];
        TermRootVersion memory rootData = termVersions[termId][version];

        // Valid if exists and not superseded
        isValid = !rootData.isSuperseded;
        publishedAt = rootData.publishedAt;

        return (isValid, termId, publishedAt);
    }

    /**
     * @notice Get total number of published terms
     */
    function getPublishedTermsCount() external view returns (uint256) {
        return publishedTerms.length;
    }

    /**
     * @notice Get term ID by index
     */
    function getPublishedTerm(uint256 _index) external view returns (string memory) {
        require(_index < publishedTerms.length, "Index out of bounds");
        return publishedTerms[_index];
    }

    /**
     * @notice Returns total number of published roots (all versions)
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

    /**
     * @notice Check if a term has been published
     */
    function isTermPublished(string memory _termId) external view returns (bool) {
        return latestVersion[_termId] > 0;
    }
}
