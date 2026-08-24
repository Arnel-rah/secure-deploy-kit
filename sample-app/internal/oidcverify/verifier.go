package oidcverify

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

type discoveryDocument struct {
	Issuer  string `json:"issuer"`
	JWKSURI string `json:"jwks_uri"`
}

type jwk struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwkSet struct {
	Keys []jwk `json:"keys"`
}

type Verifier struct {
	issuer      string
	jwksURI     string
	httpClient  *http.Client
	cacheTTL    time.Duration
	mu          sync.RWMutex
	keys        map[string]*rsa.PublicKey
	lastFetched time.Time
}

func NewVerifier(ctx context.Context, issuer string) (*Verifier, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	discoveryURL := strings.TrimRight(issuer, "/") + "/.well-known/openid-configuration"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discoveryURL, nil)
	if err != nil {
		return nil, fmt.Errorf("building discovery request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching discovery document: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discovery document returned status %d", resp.StatusCode)
	}

	var doc discoveryDocument
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, fmt.Errorf("decoding discovery document: %w", err)
	}

	v := &Verifier{
		issuer:     doc.Issuer,
		jwksURI:    doc.JWKSURI,
		httpClient: client,
		cacheTTL:   10 * time.Minute,
		keys:       make(map[string]*rsa.PublicKey),
	}

	if err := v.refreshKeys(ctx); err != nil {
		return nil, fmt.Errorf("initial JWKS fetch: %w", err)
	}

	return v, nil
}

func (v *Verifier) refreshKeys(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.jwksURI, nil)
	if err != nil {
		return err
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS endpoint returned status %d", resp.StatusCode)
	}

	var set jwkSet
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return err
	}

	keys := make(map[string]*rsa.PublicKey, len(set.Keys))
	for _, k := range set.Keys {
		if k.Kty != "RSA" {
			continue
		}
		pub, err := jwkToRSAPublicKey(k)
		if err != nil {
			continue
		}
		keys[k.Kid] = pub
	}

	v.mu.Lock()
	v.keys = keys
	v.lastFetched = time.Now()
	v.mu.Unlock()

	return nil
}

func jwkToRSAPublicKey(k jwk) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(k.N)
	if err != nil {
		return nil, fmt.Errorf("decoding modulus: %w", err)
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(k.E)
	if err != nil {
		return nil, fmt.Errorf("decoding exponent: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{N: n, E: int(e.Int64())}, nil
}

func (v *Verifier) getKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	v.mu.RLock()
	key, ok := v.keys[kid]
	stale := time.Since(v.lastFetched) > v.cacheTTL
	v.mu.RUnlock()

	if ok && !stale {
		return key, nil
	}

	if err := v.refreshKeys(ctx); err != nil {
		if ok {
			return key, nil
		}
		return nil, err
	}

	v.mu.RLock()
	key, ok = v.keys[kid]
	v.mu.RUnlock()
	if !ok {
		return nil, errors.New("signing key not found after refresh")
	}
	return key, nil
}

type Claims struct {
	Subject string   `json:"sub"`
	Issuer  string   `json:"iss"`
	Expiry  int64    `json:"exp"`
	Roles   []string `json:"roles,omitempty"`
	Raw     map[string]any
}

var (
	ErrMalformedToken = errors.New("malformed token")
	ErrInvalidSig      = errors.New("invalid signature")
	ErrExpiredToken    = errors.New("token expired")
	ErrIssuerMismatch  = errors.New("issuer mismatch")
)


func (v *Verifier) Verify(ctx context.Context, rawToken string) (*Claims, error) {
	parts := strings.Split(rawToken, ".")
	if len(parts) != 3 {
		return nil, ErrMalformedToken
	}
	headerB64, payloadB64, sigB64 := parts[0], parts[1], parts[2]

	headerJSON, err := base64.RawURLEncoding.DecodeString(headerB64)
	if err != nil {
		return nil, fmt.Errorf("%w: header: %v", ErrMalformedToken, err)
	}
	var header struct {
		Alg string `json:"alg"`
		Kid string `json:"kid"`
	}
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return nil, fmt.Errorf("%w: header: %v", ErrMalformedToken, err)
	}
	if header.Alg != "RS256" {
		return nil, fmt.Errorf("%w: unsupported alg %q", ErrMalformedToken, header.Alg)
	}

	key, err := v.getKey(ctx, header.Kid)
	if err != nil {
		return nil, fmt.Errorf("resolving signing key: %w", err)
	}

	signature, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return nil, fmt.Errorf("%w: signature: %v", ErrMalformedToken, err)
	}

	signedInput := headerB64 + "." + payloadB64
	hashed := sha256.Sum256([]byte(signedInput))
	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signature); err != nil {
		return nil, ErrInvalidSig
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return nil, fmt.Errorf("%w: payload: %v", ErrMalformedToken, err)
	}

	var raw map[string]any
	if err := json.Unmarshal(payloadJSON, &raw); err != nil {
		return nil, fmt.Errorf("%w: payload: %v", ErrMalformedToken, err)
	}

	claims := &Claims{Raw: raw}
	if err := json.Unmarshal(payloadJSON, claims); err != nil {
		return nil, fmt.Errorf("%w: claims: %v", ErrMalformedToken, err)
	}

	if claims.Expiry == 0 || time.Now().After(time.Unix(claims.Expiry, 0)) {
		return nil, ErrExpiredToken
	}
	if claims.Issuer != v.issuer {
		return nil, ErrIssuerMismatch
	}

	return claims, nil
}


func publicKeyToPEM(pub *rsa.PublicKey) (string, error) {
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}
	block := &pem.Block{Type: "PUBLIC KEY", Bytes: der}
	return string(pem.EncodeToMemory(block)), nil
}
