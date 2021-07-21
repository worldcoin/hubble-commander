package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
)

func signatureProofToCalldata(proof *models.SignatureProof) *rollup.TypesSignatureProof {
	result := &rollup.TypesSignatureProof{
		States:          make([]rollup.TypesUserState, 0, len(proof.States)),
		StateWitnesses:  make([][][32]byte, 0, len(proof.States)),
		Pubkeys:         make([][4]*big.Int, 0, len(proof.PublicKeys)),
		PubkeyWitnesses: make([][][32]byte, 0, len(proof.PublicKeys)),
	}
	for i := range proof.States {
		stateProof := stateMerkleProofToCalldata(&proof.States[i])
		result.States = append(result.States, stateProof.State)
		result.StateWitnesses = append(result.StateWitnesses, stateProof.Witness)

		result.Pubkeys = append(result.Pubkeys, proof.PublicKeys[i].PublicKey.BigInts())
		result.PubkeyWitnesses = append(result.PubkeyWitnesses, proof.PublicKeys[i].Witness.Bytes())
	}
	return result
}
