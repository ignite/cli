package cometbls

import (
	"encoding/hex"
	"errors"
	"testing"
)

const (
	testChainID               = "union-devnet-1337"
	testTrustedValidatorsHash = "20DDFE7A0F75C65D876316091ECCD494A54A2BB324C872015F73E528D53CB9C4"
	testProofHex              = "03CF56142A1E03D2445A82100FEAF70C1CD95A731ED85792AFFF5792EC0BDD2108991BB56F9043A269F88903DE616A9AB99A3C5AB778E566744B060456C5616C" +
		"06BCE7F1930421768C2CBD79F88D08EC3A52D7C9A867064E973064385E9C945E02951190DD7CE1662546733DD540188C96E608CA750FEF36B39E2577833634C7" +
		"0AE6F1A6D00DC6C21446AAF285EF35D944E8782B131300574F9A889C7E708A2325E9A78013BBE869D38B19C602DAF69644C77D177E99ED76398BCEE13C61FDBF" +
		"2E178A5BA028A36033E54D1D9A0071E82E04079A5305347EBAC6D66F6EBFA48B1DA1BF9DC5A51EFA292E1DC7B85D26F18422EB386C48CA75434039764448BB96" +
		"268DDC2CF683DDCA4BD83DF21C5631CF784375EEBE77EABC2DE77886BF1D48392C9C52E063B4A7131EAB9ABBA12A9F26888BC37366D41AC7D4BAC0BF6755ACB0" +
		"09BF9F36F380B6D0EEAABF066503A1B6E01DCC965D968D7694E01B1755E6BDD21C7A80B41682748F9B7151714BE34AA79AAD48BBB2A84525F6CDF812658C6E4F"
)

func mustDecodeHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("invalid hex: %v", err)
	}
	return b
}

func mustHash32(t *testing.T, s string) [32]byte {
	t.Helper()
	b := mustDecodeHex(t, s)
	if len(b) != 32 {
		t.Fatalf("expected 32 bytes, got %d", len(b))
	}
	var out [32]byte
	copy(out[:], b)
	return out
}

func validHeader(t *testing.T) LightHeader {
	t.Helper()
	return LightHeader{
		Height:             3405691582,
		TimeSeconds:        1732205251,
		TimeNanos:          998131342,
		ValidatorsHash:     mustHash32(t, testTrustedValidatorsHash),
		NextValidatorsHash: mustHash32(t, testTrustedValidatorsHash),
		AppHash:            mustHash32(t, "EE7E3E58F98AC95D63CE93B270981DF3EE54CA367F8D521ED1F444717595CD36"),
	}
}

// Verifies the full pipeline through the native wrapper: returns "" on
// success, a non-empty error message otherwise.
func TestXVerifyZKPHappyPath(t *testing.T) {
	tvh := mustHash32(t, testTrustedValidatorsHash)
	header := validHeader(t)
	proof := mustDecodeHex(t, testProofHex)

	got := X_verifyZKP(testChainID, tvh[:], EncodeLightHeader(header), proof)
	if got != "" {
		t.Fatalf("expected empty success string, got %q", got)
	}
}

func TestXVerifyZKPReturnsErrorString(t *testing.T) {
	tvh := mustHash32(t, testTrustedValidatorsHash)
	header := validHeader(t)
	header.TimeSeconds = 1732205252 // tamper
	proof := mustDecodeHex(t, testProofHex)

	got := X_verifyZKP(testChainID, tvh[:], EncodeLightHeader(header), proof)
	if got == "" {
		t.Fatalf("expected non-empty error string for tampered proof")
	}
	// The string should be the Error() of ErrInvalidProof (not a panic trace).
	if got != ErrInvalidProof.Error() {
		t.Fatalf("expected %q, got %q", ErrInvalidProof.Error(), got)
	}
}

func TestXVerifyZKPMalformedInputs(t *testing.T) {
	tvh := mustHash32(t, testTrustedValidatorsHash)
	header := validHeader(t)
	proof := mustDecodeHex(t, testProofHex)

	t.Run("short trusted validators hash", func(t *testing.T) {
		got := X_verifyZKP(testChainID, tvh[:31], EncodeLightHeader(header), proof)
		if got == "" {
			t.Fatalf("expected error for short tvh")
		}
	})
	t.Run("short header", func(t *testing.T) {
		got := X_verifyZKP(testChainID, tvh[:], EncodeLightHeader(header)[:100], proof)
		if got != ErrInvalidHeaderLen.Error() {
			t.Fatalf("expected %q, got %q", ErrInvalidHeaderLen.Error(), got)
		}
	})
	t.Run("short proof", func(t *testing.T) {
		got := X_verifyZKP(testChainID, tvh[:], EncodeLightHeader(header), proof[:383])
		if got != ErrInvalidRawProof.Error() {
			t.Fatalf("expected %q, got %q", ErrInvalidRawProof.Error(), got)
		}
	})
	t.Run("chain id too long", func(t *testing.T) {
		got := X_verifyZKP("abcdefghijklmnopqrstuvwxyz0123456", tvh[:], EncodeLightHeader(header), proof)
		if got != ErrInvalidChainIDLen.Error() {
			t.Fatalf("expected %q, got %q", ErrInvalidChainIDLen.Error(), got)
		}
	})
}

// Direct-call tests for the exported Go API (used from Go consumers, not just
// via the gno native binding).

func TestVerifyZKPHappyPath(t *testing.T) {
	tvh := mustHash32(t, testTrustedValidatorsHash)
	header := validHeader(t)
	proof := mustDecodeHex(t, testProofHex)

	if err := VerifyZKP(testChainID, tvh[:], EncodeLightHeader(header), proof); err != nil {
		t.Fatalf("VerifyZKP: %v", err)
	}
}

func TestVerifyZKPTamperedBlock(t *testing.T) {
	tvh := mustHash32(t, testTrustedValidatorsHash)
	header := validHeader(t)
	header.TimeSeconds = 1732205252
	proof := mustDecodeHex(t, testProofHex)

	err := VerifyZKP(testChainID, tvh[:], EncodeLightHeader(header), proof)
	if !errors.Is(err, ErrInvalidProof) {
		t.Fatalf("expected ErrInvalidProof, got %v", err)
	}
}

func TestParseZKPLengthCheck(t *testing.T) {
	if _, err := ParseZKP(make([]byte, ExpectedProofSize-1)); !errors.Is(err, ErrInvalidRawProof) {
		t.Fatalf("expected ErrInvalidRawProof for short input, got %v", err)
	}
	if _, err := ParseZKP(make([]byte, ExpectedProofSize+1)); !errors.Is(err, ErrInvalidRawProof) {
		t.Fatalf("expected ErrInvalidRawProof for long input, got %v", err)
	}
	// A zero-filled 384-byte buffer parses all five points as the identity,
	// which is a valid affine encoding in gnark-crypto; the function should
	// succeed.
	if _, err := ParseZKP(make([]byte, ExpectedProofSize)); err != nil {
		t.Fatalf("expected zero-identity buffer to parse, got %v", err)
	}
}

func TestParseZKPGenuineVector(t *testing.T) {
	proof := mustDecodeHex(t, testProofHex)
	zkp, err := ParseZKP(proof)
	if err != nil {
		t.Fatalf("ParseZKP: %v", err)
	}
	// Sanity: all five G1 coords should be non-zero for the test vector.
	if zkp.Proof.A.X.IsZero() || zkp.Proof.A.Y.IsZero() {
		t.Fatalf("Proof.A parsed as zero")
	}
	if zkp.ProofCommitment.X.IsZero() {
		t.Fatalf("ProofCommitment parsed as zero")
	}
	if zkp.ProofCommitmentPoK.X.IsZero() {
		t.Fatalf("ProofCommitmentPoK parsed as zero")
	}
}

func TestEncodeDecodeLightHeaderRoundTrip(t *testing.T) {
	h := LightHeader{
		Height:      1,
		TimeSeconds: 2,
		TimeNanos:   3,
	}
	for i := range byte(32) {
		h.ValidatorsHash[i] = i + 1
		h.NextValidatorsHash[i] = i + 2
		h.AppHash[i] = i + 3
	}
	buf := EncodeLightHeader(h)
	if len(buf) != EncodedLightHeaderSize {
		t.Fatalf("encoded length %d, want %d", len(buf), EncodedLightHeaderSize)
	}
	got, err := DecodeLightHeader(buf)
	if err != nil {
		t.Fatalf("DecodeLightHeader: %v", err)
	}
	if got != h {
		t.Fatalf("round-trip mismatch:\n got  %+v\n want %+v", got, h)
	}
}

func TestDecodeLightHeaderRejectsNegatives(t *testing.T) {
	header := validHeader(t)
	buf := EncodeLightHeader(header)

	t.Run("negative height", func(t *testing.T) {
		// Overwrite height bytes with a negative i64 (0xFF...FF).
		bad := append([]byte(nil), buf...)
		for i := range 8 {
			bad[i] = 0xFF
		}
		if _, err := DecodeLightHeader(bad); !errors.Is(err, ErrInvalidHeight) {
			t.Fatalf("expected ErrInvalidHeight, got %v", err)
		}
	})
	t.Run("negative timestamp seconds", func(t *testing.T) {
		bad := append([]byte(nil), buf...)
		for i := 8; i < 16; i++ {
			bad[i] = 0xFF
		}
		if _, err := DecodeLightHeader(bad); !errors.Is(err, ErrInvalidTimestamp) {
			t.Fatalf("expected ErrInvalidTimestamp, got %v", err)
		}
	})
	t.Run("negative timestamp nanos", func(t *testing.T) {
		bad := append([]byte(nil), buf...)
		for i := 16; i < 20; i++ {
			bad[i] = 0xFF
		}
		if _, err := DecodeLightHeader(bad); !errors.Is(err, ErrInvalidTimestamp) {
			t.Fatalf("expected ErrInvalidTimestamp, got %v", err)
		}
	})
	t.Run("wrong length", func(t *testing.T) {
		if _, err := DecodeLightHeader(buf[:50]); !errors.Is(err, ErrInvalidHeaderLen) {
			t.Fatalf("expected ErrInvalidHeaderLen, got %v", err)
		}
	})
}

func TestPublicInputsDeterministic(t *testing.T) {
	tvh := mustHash32(t, testTrustedValidatorsHash)
	header := validHeader(t)
	proof := mustDecodeHex(t, testProofHex)
	zkp, err := ParseZKP(proof)
	if err != nil {
		t.Fatalf("ParseZKP: %v", err)
	}

	inputs1, err := PublicInputs(testChainID, tvh, header, zkp)
	if err != nil {
		t.Fatalf("PublicInputs: %v", err)
	}
	inputs2, err := PublicInputs(testChainID, tvh, header, zkp)
	if err != nil {
		t.Fatalf("PublicInputs second call: %v", err)
	}
	if inputs1 != inputs2 {
		t.Fatalf("PublicInputs not deterministic: %v vs %v", inputs1, inputs2)
	}
	// Changing the header must change inputsHash.
	header.Height++
	inputs3, err := PublicInputs(testChainID, tvh, header, zkp)
	if err != nil {
		t.Fatalf("PublicInputs after mutation: %v", err)
	}
	if inputs3[0] == inputs1[0] {
		t.Fatalf("inputsHash unchanged after header mutation")
	}
	// commitmentHash depends only on the proof commitment, not the header,
	// so it must be stable.
	if inputs3[1] != inputs1[1] {
		t.Fatalf("commitmentHash changed with header mutation — should only depend on proof commitment")
	}
}
