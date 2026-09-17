package cometbls

// X_verifyZKP is the host-side native binding for the `verifyZKP` function
// declared in cometbls.gno. It returns an empty string on successful proof
// verification and a non-empty error message otherwise — gno cannot currently
// return a Go `error` across the native ABI, so the gno wrapper converts the
// string back into an error value.
//
// `headerEncoded` must be exactly EncodedLightHeaderSize bytes (116); see
// EncodeLightHeader in cometbls.gno for the canonical encoder.
func X_verifyZKP(chainID string, trustedValidatorsHash []byte, headerEncoded []byte, zkp []byte) string {
	if err := VerifyZKP(chainID, trustedValidatorsHash, headerEncoded, zkp); err != nil {
		return err.Error()
	}
	return ""
}
