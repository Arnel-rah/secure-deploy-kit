package oidcverify

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testOIDCProvider struct {
	server     *httptest.Server
	privateKey *rsa.PrivateKey
	kid        string
	issuer     string
}

func newTestOIDCProvider(t *testing.T) *testOIDCProvider {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating RSA key: %v", err)
	}

	p := &testOIDCProvider{privateKey: privateKey, kid: "test-key-1"}

	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"issuer":   p.issuer,
			"jwks_uri": p.issuer + "/jwks",
		})
	})
	mux.HandleFunc("/jwks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"keys": []map[string]string{
				{
					"kty": "RSA",
					"kid": p.kid,
					"n":   base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.N.Bytes()),
					"e":   base64.RawURLEncoding.EncodeToString(bigIntToBytes(privateKey.PublicKey.E)),
				},
			},
		})
	})

	p.server = httptest.NewServer(mux)
	p.issuer = p.server.URL
	return p
}

func bigIntToBytes(e int) []byte {
	if e == 0 {
		return []byte{0}
	}
	var b []byte
	for e > 0 {
		b = append([]byte{byte(e & 0xff)}, b...)
		e >>= 8
	}
	return b
}

func (p *testOIDCProvider) close() {
	p.server.Close()
}

func (p *testOIDCProvider) issueToken(t *testing.T, subject string, roles []string, expiresIn time.Duration) string {
	t.Helper()

	header := map[string]string{"alg": "RS256", "typ": "JWT", "kid": p.kid}
	headerJSON, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	claims := map[string]any{
		"sub":   subject,
		"iss":   p.issuer,
		"exp":   time.Now().Add(expiresIn).Unix(),
		"roles": roles,
	}
	claimsJSON, _ := json.Marshal(claims)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signedInput := headerB64 + "." + claimsB64
	hashed := sha256.Sum256([]byte(signedInput))

	sig, err := rsa.SignPKCS1v15(rand.Reader, p.privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		t.Fatalf("signing token: %v", err)
	}
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	return signedInput + "." + sigB64
}

func TestVerifier_ValidToken(t *testing.T) {
	provider := newTestOIDCProvider(t)
	defer provider.close()

	ctx := context.Background()
	verifier, err := NewVerifier(ctx, provider.issuer)
	if err != nil {
		t.Fatalf("NewVerifier: %v", err)
	}

	token := provider.issueToken(t, "user-123", []string{"admin"}, time.Hour)

	claims, err := verifier.Verify(ctx, token)
	if err != nil {
		t.Fatalf("Verify returned error for a valid token: %v", err)
	}
	if claims.Subject != "user-123" {
		t.Errorf("subject = %q, want %q", claims.Subject, "user-123")
	}
	if len(claims.Roles) != 1 || claims.Roles[0] != "admin" {
		t.Errorf("roles = %v, want [admin]", claims.Roles)
	}
}

func TestVerifier_ExpiredToken(t *testing.T) {
	provider := newTestOIDCProvider(t)
	defer provider.close()

	ctx := context.Background()
	verifier, err := NewVerifier(ctx, provider.issuer)
	if err != nil {
		t.Fatalf("NewVerifier: %v", err)
	}

	token := provider.issueToken(t, "user-123", nil, -time.Hour)

	_, err = verifier.Verify(ctx, token)
	if err != ErrExpiredToken {
		t.Errorf("Verify error = %v, want ErrExpiredToken", err)
	}
}

func TestVerifier_WrongIssuer(t *testing.T) {
	provider := newTestOIDCProvider(t)
	defer provider.close()

	ctx := context.Background()
	verifier, err := NewVerifier(ctx, provider.issuer)
	if err != nil {
		t.Fatalf("NewVerifier: %v", err)
	}

	header := map[string]string{"alg": "RS256", "typ": "JWT", "kid": provider.kid}
	headerJSON, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	claims := map[string]any{
		"sub": "attacker",
		"iss": "https://not-the-real-issuer.example.com",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	claimsJSON, _ := json.Marshal(claims)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signedInput := headerB64 + "." + claimsB64
	hashed := sha256.Sum256([]byte(signedInput))
	sig, _ := rsa.SignPKCS1v15(rand.Reader, provider.privateKey, crypto.SHA256, hashed[:])
	token := signedInput + "." + base64.RawURLEncoding.EncodeToString(sig)

	_, err = verifier.Verify(ctx, token)
	if err != ErrIssuerMismatch {
		t.Errorf("Verify error = %v, want ErrIssuerMismatch", err)
	}
}

func TestVerifier_TamperedSignature(t *testing.T) {
	provider := newTestOIDCProvider(t)
	defer provider.close()

	ctx := context.Background()
	verifier, err := NewVerifier(ctx, provider.issuer)
	if err != nil {
		t.Fatalf("NewVerifier: %v", err)
	}

	token := provider.issueToken(t, "user-123", nil, time.Hour)
	tampered := token[:len(token)-4] + "AAAA"

	_, err = verifier.Verify(ctx, tampered)
	if err != ErrInvalidSig {
		t.Errorf("Verify error = %v, want ErrInvalidSig", err)
	}
}

func TestVerifier_MalformedToken(t *testing.T) {
	provider := newTestOIDCProvider(t)
	defer provider.close()

	ctx := context.Background()
	verifier, err := NewVerifier(ctx, provider.issuer)
	if err != nil {
		t.Fatalf("NewVerifier: %v", err)
	}

	_, err = verifier.Verify(ctx, "not-a-jwt-at-all")
	if err != ErrMalformedToken {
		t.Errorf("Verify error = %v, want ErrMalformedToken", err)
	}
}
