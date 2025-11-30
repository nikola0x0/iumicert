# IUMiCertRegistry v2 Deployment Instructions

## Updated Contract Features
✅ Root versioning support (v1, v2, v3...)
✅ Revocation through root supersession
✅ Backward compatible verification
✅ Complete version history tracking

## Pre-Deployment Checklist

1. ✅ Contract compiled successfully
2. ✅ All 15 tests passed
3. ⚠️ Need Sepolia ETH in deployer wallet
4. ⚠️ Need PRIVATE_KEY in .env

## Deployment Command

```bash
# Make sure you're in packages/contracts directory
cd packages/contracts

# Deploy to Sepolia
forge script script/Deploy.s.sol:Deploy --rpc-url $SEPOLIA_RPC_URL --broadcast --verify -vvvv

# Or if using environment variable for RPC:
forge script script/Deploy.s.sol:Deploy --rpc-url sepolia --broadcast --verify -vvvv
```

## Expected Output

```
Deploying contracts with owner: 0xYourAddress...
IUMiCertRegistry deployed to: 0xNEW_CONTRACT_ADDRESS
```

## Post-Deployment: Update Contract Addresses

Update the following files with the new contract address:

### Backend (.env files):
1. `packages/issuer/.env`
   - Line 9: `IUMICERT_CONTRACT_ADDRESS=0xNEW_ADDRESS`

2. `packages/issuer/.env.production`
   - Line 11: `IUMICERT_CONTRACT_ADDRESS=0xNEW_ADDRESS`

3. `packages/issuer/.env.example`
   - Update example address

### Frontend (.env files):
4. `packages/issuer/web/iumicert-issuer/.env.local`
   - Line 2: `NEXT_PUBLIC_IUMICERT_CONTRACT_ADDRESS=0xNEW_ADDRESS`

5. `packages/issuer/web/iumicert-issuer/.env.example`
   - Line 2: Update example address

6. `packages/client/iumicert-client/.env` (if exists)
   - Update with new address

### Documentation:
7. `CLAUDE.md`
   - Update contract address in Smart Contracts section

8. `packages/issuer/README.md`
   - Update contract address

## Verification

After deployment, verify on Sepolia Etherscan:
https://sepolia.etherscan.io/address/0xNEW_ADDRESS

Check that:
✅ Contract is verified (code visible)
✅ Owner is set correctly
✅ No initial term roots published yet

## Testing New Deployment

```bash
# From packages/issuer directory
# Test publishing a term root with the new contract
./micert publish-roots Semester_1_2023

# Verify it worked
cast call 0xNEW_ADDRESS "latestVersion(string)(uint256)" "Semester_1_2023" --rpc-url sepolia
# Should return: 1 (version 1)

# Check term history
cast call 0xNEW_ADDRESS "getTermHistory(string)(uint256[],bytes32[])" "Semester_1_2023" --rpc-url sepolia
```

## Rollback Plan (If Issues)

The old contract at `0x4bE58F5EaFDa3b09BA87c2F5Eb17a23c37C0dD60` remains functional.
Simply revert the .env files to use the old address.

## Next Steps After Deployment

1. Update Go blockchain integration code (see GO_INTEGRATION.md)
2. Update frontend to use new contract functions
3. Test complete revocation workflow
4. Update thesis with new contract address
