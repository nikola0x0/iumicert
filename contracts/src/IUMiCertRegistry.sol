// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/**
 * @title IUMiCertRegistry
 * @dev Smart contract for IU-MiCert hybrid credential system
 * Stores term-level Verkle tree roots with temporal verification
 */
contract IUMiCertRegistry {
    
    // Term tree structure stored on-chain
    struct TermTree {
        bytes32 verkleRoot;      // Root of the term's Verkle tree
        uint256 publishedAt;     // Block timestamp when published
        uint32 version;          // Version counter for updates
        bool active;             // Whether this term is still active
    }
    
    // Mapping from term ID hash to term tree data
    mapping(bytes32 => TermTree) public termTrees;
    
    // Mapping from institution address to authorized status
    mapping(address => bool) public authorizedInstitutions;
    
    // Institution that deployed this contract
    address public deployer;
    
    // Events
    event TermTreePublished(
        bytes32 indexed termIdHash,
        bytes32 indexed verkleRoot,
        uint256 publishedAt,
        uint32 version
    );
    
    event TermTreeUpdated(
        bytes32 indexed termIdHash,
        bytes32 indexed oldRoot,
        bytes32 indexed newRoot,
        uint32 newVersion
    );
    
    event InstitutionAuthorized(address indexed institution);
    event InstitutionRevoked(address indexed institution);
    
    // Modifiers
    modifier onlyAuthorized() {
        require(
            authorizedInstitutions[msg.sender] || msg.sender == deployer,
            "Not authorized to publish terms"
        );
        _;
    }
    
    modifier onlyDeployer() {
        require(msg.sender == deployer, "Only deployer can manage institutions");
        _;
    }
    
    constructor() {
        deployer = msg.sender;
        authorizedInstitutions[msg.sender] = true;
    }
    
    /**
     * @dev Publish a new term tree or update an existing one
     * @param termIdHash Hash of the term identifier (e.g., H256("Fall_2024"))
     * @param verkleRoot Root commitment of the term's Verkle tree
     */
    function publishTermTree(
        bytes32 termIdHash,
        bytes32 verkleRoot
    ) external onlyAuthorized {
        require(verkleRoot != bytes32(0), "Invalid verkle root");
        
        TermTree storage termTree = termTrees[termIdHash];
        
        if (termTree.publishedAt == 0) {
            // New term
            termTree.verkleRoot = verkleRoot;
            termTree.publishedAt = block.timestamp;
            termTree.version = 1;
            termTree.active = true;
            
            emit TermTreePublished(termIdHash, verkleRoot, block.timestamp, 1);
        } else {
            // Update existing term
            require(termTree.active, "Term is no longer active");
            
            bytes32 oldRoot = termTree.verkleRoot;
            termTree.verkleRoot = verkleRoot;
            termTree.version++;
            
            emit TermTreeUpdated(termIdHash, oldRoot, verkleRoot, termTree.version);
        }
    }
    
    /**
     * @dev Deactivate a term (prevent further updates)
     * @param termIdHash Hash of the term identifier
     */
    function deactivateTerm(bytes32 termIdHash) external onlyAuthorized {
        TermTree storage termTree = termTrees[termIdHash];
        require(termTree.publishedAt != 0, "Term does not exist");
        require(termTree.active, "Term already inactive");
        
        termTree.active = false;
    }
    
    /**
     * @dev Get term tree information
     * @param termIdHash Hash of the term identifier
     * @return verkleRoot Root of the Verkle tree
     * @return publishedAt Timestamp when published
     * @return version Current version
     * @return active Whether the term is active
     */
    function getTermTree(bytes32 termIdHash) 
        external 
        view 
        returns (
            bytes32 verkleRoot,
            uint256 publishedAt,
            uint32 version,
            bool active
        ) 
    {
        TermTree storage termTree = termTrees[termIdHash];
        require(termTree.publishedAt != 0, "Term does not exist");
        
        return (
            termTree.verkleRoot,
            termTree.publishedAt,
            termTree.version,
            termTree.active
        );
    }
    
    /**
     * @dev Verify that a credential was issued before the term publication
     * This is critical for anti-forgery protection
     * @param termIdHash Hash of the term identifier
     * @param issuedTimestamp Unix timestamp when credential was issued
     * @return valid Whether the timing is valid
     */
    function verifyCredentialTiming(
        bytes32 termIdHash,
        uint256 issuedTimestamp
    ) external view returns (bool valid) {
        TermTree storage termTree = termTrees[termIdHash];
        require(termTree.publishedAt != 0, "Term does not exist");
        
        // Credential must be issued before or at the time of term publication
        return issuedTimestamp <= termTree.publishedAt;
    }
    
    /**
     * @dev Authorize an institution to publish terms
     * @param institution Address of the institution
     */
    function authorizeInstitution(address institution) external onlyDeployer {
        require(institution != address(0), "Invalid institution address");
        require(!authorizedInstitutions[institution], "Already authorized");
        
        authorizedInstitutions[institution] = true;
        emit InstitutionAuthorized(institution);
    }
    
    /**
     * @dev Revoke institution authorization
     * @param institution Address of the institution
     */
    function revokeInstitution(address institution) external onlyDeployer {
        require(authorizedInstitutions[institution], "Not authorized");
        require(institution != deployer, "Cannot revoke deployer");
        
        authorizedInstitutions[institution] = false;
        emit InstitutionRevoked(institution);
    }
    
    /**
     * @dev Batch publish multiple terms (gas optimization)
     * @param termIdHashes Array of term identifier hashes
     * @param verkleRoots Array of corresponding Verkle roots
     */
    function batchPublishTermTrees(
        bytes32[] calldata termIdHashes,
        bytes32[] calldata verkleRoots
    ) external onlyAuthorized {
        require(
            termIdHashes.length == verkleRoots.length,
            "Array length mismatch"
        );
        require(termIdHashes.length <= 50, "Too many terms in batch");
        
        for (uint256 i = 0; i < termIdHashes.length; i++) {
            publishTermTree(termIdHashes[i], verkleRoots[i]);
        }
    }
    
    /**
     * @dev Get multiple term trees (gas-efficient batch query)
     * @param termIdHashes Array of term identifier hashes
     * @return verkleRoots Array of Verkle roots
     * @return publishedAts Array of publication timestamps
     * @return versions Array of versions
     * @return actives Array of active statuses
     */
    function getBatchTermTrees(bytes32[] calldata termIdHashes)
        external
        view
        returns (
            bytes32[] memory verkleRoots,
            uint256[] memory publishedAts,
            uint32[] memory versions,
            bool[] memory actives
        )
    {
        uint256 length = termIdHashes.length;
        verkleRoots = new bytes32[](length);
        publishedAts = new uint256[](length);
        versions = new uint32[](length);
        actives = new bool[](length);
        
        for (uint256 i = 0; i < length; i++) {
            TermTree storage termTree = termTrees[termIdHashes[i]];
            verkleRoots[i] = termTree.verkleRoot;
            publishedAts[i] = termTree.publishedAt;
            versions[i] = termTree.version;
            actives[i] = termTree.active;
        }
    }
}

/**
 * @title IUMiCertVerifier
 * @dev Verification utilities for off-chain proof verification
 * This contract provides helper functions for verifiers
 */
contract IUMiCertVerifier {
    
    // Reference to the main registry contract
    IUMiCertRegistry public immutable registry;
    
    constructor(address _registry) {
        require(_registry != address(0), "Invalid registry address");
        registry = IUMiCertRegistry(_registry);
    }
    
    /**
     * @dev Verify a complete credential receipt
     * @param termIdHash Hash of the term identifier
     * @param studentDIDHash Hash of the student's DID
     * @param studentTermRoot Expected student term root from Verkle proof
     * @param maxIssuedTimestamp Latest issuance timestamp among revealed courses
     * @return valid Whether the credential receipt is valid
     * @return verkleRoot The term's Verkle root for off-chain verification
     */
    function verifyCredentialReceipt(
        bytes32 termIdHash,
        bytes32 studentDIDHash,
        bytes32 studentTermRoot,
        uint256 maxIssuedTimestamp
    ) external view returns (bool valid, bytes32 verkleRoot) {
        // Get term information
        (
            bytes32 termVerkleRoot,
            uint256 publishedAt,
            uint32 version,
            bool active
        ) = registry.getTermTree(termIdHash);
        
        // Verify timing: all courses must be issued before term publication
        if (maxIssuedTimestamp > publishedAt) {
            return (false, bytes32(0));
        }
        
        // Return the Verkle root for off-chain Verkle proof verification
        // The calling application must verify:
        // 1. Verkle proof: studentDIDHash → studentTermRoot under termVerkleRoot
        // 2. Merkle proofs: individual courses under studentTermRoot
        return (true, termVerkleRoot);
    }
    
    /**
     * @dev Hash term ID string to bytes32
     * @param termId String term identifier (e.g., "Fall_2024")
     * @return Hash of the term ID
     */
    function hashTermId(string calldata termId) external pure returns (bytes32) {
        return keccak256(abi.encodePacked(termId));
    }
    
    /**
     * @dev Hash student DID to bytes32
     * @param studentDID String student DID (e.g., "did:example:student001")
     * @return Hash of the student DID
     */
    function hashStudentDID(string calldata studentDID) external pure returns (bytes32) {
        return keccak256(abi.encodePacked(studentDID));
    }
    
    /**
     * @dev Batch verify timing for multiple credentials
     * @param termIdHashes Array of term identifier hashes
     * @param issuedTimestamps Array of issuance timestamps
     * @return valid Array of validity results
     */
    function batchVerifyTiming(
        bytes32[] calldata termIdHashes,
        uint256[] calldata issuedTimestamps
    ) external view returns (bool[] memory valid) {
        require(
            termIdHashes.length == issuedTimestamps.length,
            "Array length mismatch"
        );
        
        valid = new bool[](termIdHashes.length);
        
        for (uint256 i = 0; i < termIdHashes.length; i++) {
            valid[i] = registry.verifyCredentialTiming(
                termIdHashes[i],
                issuedTimestamps[i]
            );
        }
    }
}