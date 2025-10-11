// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {Ownable} from "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title ReceiptRevocationRegistry
 * @author IU-MiCert Team
 * @notice A simple on-chain registry for revoking academic journey receipts.
 * The owner of the contract (IU Vietnam) can revoke a receipt by its unique ID.
 * Verifiers can query this contract to check if a receipt is still valid.
 */
contract ReceiptRevocationRegistry is Ownable {
    
    struct RevocationRecord {
        bool isRevoked;
        uint256 revokedAt;
        string reason;
    }

    // Mapping from a unique receipt identifier to its revocation status.
    mapping(bytes32 => RevocationRecord) public revokedReceipts;
    
    // Array to track all revoked receipts
    bytes32[] public revokedReceiptIds;

    event ReceiptRevoked(bytes32 indexed receiptId, string reason, uint256 timestamp);

    /**
     * @notice Initializes the contract, setting the deployer as the initial owner.
     */
    constructor(address initialOwner) Ownable(initialOwner) {}

    /**
     * @notice Revokes an academic journey receipt.
     * @dev Only the owner (IU Vietnam) can call this function.
     * A receipt can only be revoked once.
     * @param _receiptId The unique identifier of the receipt to revoke.
     * @param _reason Human-readable reason for revocation.
     */
    function revokeReceipt(bytes32 _receiptId, string memory _reason) public onlyOwner {
        require(!revokedReceipts[_receiptId].isRevoked, "Receipt already revoked");
        require(bytes(_reason).length > 0, "Reason required");
        
        revokedReceipts[_receiptId] = RevocationRecord({
            isRevoked: true,
            revokedAt: block.timestamp,
            reason: _reason
        });
        
        revokedReceiptIds.push(_receiptId);

        emit ReceiptRevoked(_receiptId, _reason, block.timestamp);
    }

    /**
     * @notice Checks if a receipt is revoked
     * @param _receiptId The unique identifier of the receipt to check.
     * @return isRevoked Whether the receipt is revoked
     * @return revokedAt When it was revoked (0 if not revoked)  
     * @return reason The revocation reason
     */
    function isReceiptRevoked(bytes32 _receiptId) 
        external 
        view 
        returns (bool isRevoked, uint256 revokedAt, string memory reason) 
    {
        RevocationRecord memory record = revokedReceipts[_receiptId];
        return (record.isRevoked, record.revokedAt, record.reason);
    }

    /**
     * @notice Gets total number of revoked receipts
     */
    function getTotalRevocations() external view returns (uint256) {
        return revokedReceiptIds.length;
    }

    /**
     * @notice Gets a revoked receipt ID by index
     */
    function getRevokedReceiptId(uint256 _index) external view returns (bytes32) {
        require(_index < revokedReceiptIds.length, "Index out of bounds");
        return revokedReceiptIds[_index];
    }
}