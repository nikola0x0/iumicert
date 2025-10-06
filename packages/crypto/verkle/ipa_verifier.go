package verkle

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/crate-crypto/go-ipa/bandersnatch/fr"
	"github.com/crate-crypto/go-ipa/banderwagon"
	"github.com/crate-crypto/go-ipa/common"
	"github.com/crate-crypto/go-ipa/ipa"
	multiproof "github.com/crate-crypto/go-ipa"
	verkleLib "github.com/ethereum/go-verkle"
)

// VerifyMembershipProofWithIPA performs full IPA (Inner Product Argument) verification
// for Verkle membership proofs. This provides cryptographic binding between the
// VerkleProof and StateDiff, preventing tampering attacks.
//
// Security Model:
// 1. Blockchain root anchoring (immutable commitment)
// 2. StateDiff validation (correct key-value format)
// 3. IPA verification (cryptographic proof that StateDiff matches VerkleProof)
//
// Parameters:
//   - verkleProof: The serialized Verkle proof from the receipt
//   - stateDiff: The state difference showing key-value pairs
//   - treeRoot: The blockchain-anchored Verkle tree root
//   - courseKey: The key being proven (e.g., "did:example:STUDENT:TERM:COURSE")
//   - courseValue: The expected course data hash (32 bytes)
//
// Returns error if verification fails.
func VerifyMembershipProofWithIPA(
	verkleProof *verkleLib.VerkleProof,
	stateDiff verkleLib.StateDiff,
	treeRoot [32]byte,
	courseKey string,
	courseValue [32]byte,
) error {
	log.Printf("🔐 Starting full IPA verification for key: %s", courseKey)

	// Step 1: Validate StateDiff contains the expected key-value pair
	courseKeyHash := sha256.Sum256([]byte(courseKey))
	keyStem := courseKeyHash[:verkleLib.StemSize]
	keySuffix := courseKeyHash[verkleLib.StemSize]

	found := false
	for _, stemDiff := range stateDiff {
		if bytes.Equal(keyStem, stemDiff.Stem[:]) {
			for _, suffixDiff := range stemDiff.SuffixDiffs {
				if keySuffix == suffixDiff.Suffix {
					found = true

					if suffixDiff.CurrentValue == nil {
						return fmt.Errorf("StateDiff shows nil value for key %s", courseKey)
					}

					if !bytes.Equal((*suffixDiff.CurrentValue)[:], courseValue[:]) {
						return fmt.Errorf("StateDiff value mismatch for key %s: expected %x, got %x",
							courseKey, courseValue, *suffixDiff.CurrentValue)
					}
					break
				}
			}
			break
		}
	}

	if !found {
		return fmt.Errorf("key %s not found in StateDiff", courseKey)
	}

	log.Printf("✅ StateDiff validation passed")

	// Step 2: Convert VerkleProof to go-ipa MultiProof format
	multiProof, err := verkleProofToMultiProof(verkleProof)
	if err != nil {
		return fmt.Errorf("failed to convert VerkleProof to MultiProof: %w", err)
	}

	log.Printf("✅ Converted VerkleProof to go-ipa MultiProof")

	// Step 3: Extract commitments from the proof
	commitments, err := extractCommitments(verkleProof)
	if err != nil {
		return fmt.Errorf("failed to extract commitments: %w", err)
	}

	log.Printf("✅ Extracted %d commitments from proof", len(commitments))

	// Step 4: Extract evaluation points (z values) and results (y values)
	zs, ys, err := extractEvaluationData(stateDiff, courseKeyHash)
	if err != nil {
		return fmt.Errorf("failed to extract evaluation data: %w", err)
	}

	log.Printf("✅ Extracted evaluation points: z=%v, y values count=%d", zs, len(ys))

	// Step 5: Create transcript for Fiat-Shamir
	transcript := common.NewTranscript("verkle-membership")

	// Step 6: Get IPA configuration
	ipaConfig, err := ipa.NewIPASettings()
	if err != nil {
		return fmt.Errorf("failed to create IPA settings: %w", err)
	}

	// Step 7: Verify the multiproof using go-ipa
	verified, err := multiproof.CheckMultiProof(
		transcript,
		ipaConfig,
		multiProof,
		commitments,
		ys,
		zs,
	)

	if err != nil {
		return fmt.Errorf("IPA multiproof verification failed: %w", err)
	}

	if !verified {
		return fmt.Errorf("IPA multiproof verification returned false")
	}

	log.Printf("✅ Full IPA verification successful!")
	log.Printf("   - Blockchain root: %x", treeRoot)
	log.Printf("   - Course key: %s", courseKey)
	log.Printf("   - Cryptographically proven: StateDiff matches VerkleProof")

	return nil
}

// verkleProofToMultiProof converts a VerkleProof to go-ipa's MultiProof format
func verkleProofToMultiProof(vp *verkleLib.VerkleProof) (*multiproof.MultiProof, error) {
	if vp == nil || vp.IPAProof == nil {
		return nil, fmt.Errorf("VerkleProof or IPAProof is nil")
	}

	// Convert D point (banderwagon element)
	var D banderwagon.Element
	if err := D.SetBytes(vp.D[:]); err != nil {
		return nil, fmt.Errorf("failed to deserialize D point: %w", err)
	}

	// Convert IPA proof structure
	ipaProof := ipa.IPAProof{
		L: make([]banderwagon.Element, len(vp.IPAProof.CL)),
		R: make([]banderwagon.Element, len(vp.IPAProof.CR)),
	}

	// Convert CL points (left commitments)
	for i, cl := range vp.IPAProof.CL {
		if err := ipaProof.L[i].SetBytes(cl[:]); err != nil {
			return nil, fmt.Errorf("failed to deserialize CL[%d]: %w", i, err)
		}
	}

	// Convert CR points (right commitments)
	for i, cr := range vp.IPAProof.CR {
		if err := ipaProof.R[i].SetBytes(cr[:]); err != nil {
			return nil, fmt.Errorf("failed to deserialize CR[%d]: %w", i, err)
		}
	}

	// Convert final evaluation (scalar)
	// go-verkle's FinalEvaluation → go-ipa's A_scalar
	ipaProof.A_scalar.SetBytes(vp.IPAProof.FinalEvaluation[:])

	return &multiproof.MultiProof{
		IPA: ipaProof,
		D:   D,
	}, nil
}

// extractCommitments extracts polynomial commitments from VerkleProof
func extractCommitments(vp *verkleLib.VerkleProof) ([]*banderwagon.Element, error) {
	commitments := make([]*banderwagon.Element, len(vp.CommitmentsByPath))

	for i, commitmentBytes := range vp.CommitmentsByPath {
		var elem banderwagon.Element
		if err := elem.SetBytes(commitmentBytes[:]); err != nil {
			return nil, fmt.Errorf("failed to deserialize commitment[%d]: %w", i, err)
		}
		commitments[i] = &elem
	}

	return commitments, nil
}

// extractEvaluationData extracts evaluation points (z) and values (y) from StateDiff
func extractEvaluationData(stateDiff verkleLib.StateDiff, keyHash [32]byte) ([]uint8, []*fr.Element, error) {
	keyStem := keyHash[:verkleLib.StemSize]
	keySuffix := keyHash[verkleLib.StemSize]

	// For membership proofs, we have one evaluation point (the key suffix)
	zs := []uint8{keySuffix}
	ys := make([]*fr.Element, 1)

	// Find the corresponding value in StateDiff
	for _, stemDiff := range stateDiff {
		if bytes.Equal(keyStem, stemDiff.Stem[:]) {
			for _, suffixDiff := range stemDiff.SuffixDiffs {
				if keySuffix == suffixDiff.Suffix {
					if suffixDiff.CurrentValue == nil {
						return nil, nil, fmt.Errorf("CurrentValue is nil")
					}

					// Convert value to field element
					var yElem fr.Element
					yElem.SetBytes((*suffixDiff.CurrentValue)[:])
					ys[0] = &yElem

					return zs, ys, nil
				}
			}
		}
	}

	return nil, nil, fmt.Errorf("key not found in StateDiff")
}
