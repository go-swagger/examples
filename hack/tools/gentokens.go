// SPDX-FileCopyrightText: Copyright 2015-2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type tokenSpec struct {
	filename string
	claims   map[string]any
}

// runGenTokens generates an RSA keypair and JWT tokens for the composed-auth example.
//
// It produces:
//   - keys/apiKey.prv / keys/apiKey.pem — RSA 4096-bit private/public key
//   - tokens/token-bearer-inventory-manager.jwt
//   - tokens/token-apikey-reseller.jwt
//   - tokens/token-apikey-customer.jwt
//
// All tokens are RS256-signed with the generated private key.
func runGenTokens() error {
	outDir := "composed-auth"
	if len(os.Args) > 2 {
		outDir = os.Args[2]
	}

	keysDir := filepath.Join(outDir, "keys")
	tokensDir := filepath.Join(outDir, "tokens")

	fmt.Printf("Generating RSA keypair and JWT tokens in %s/\n", outDir)

	// Generate RSA keypair
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("generating RSA key: %w", err)
	}

	if err := os.MkdirAll(keysDir, 0o755); err != nil {
		return fmt.Errorf("creating keys directory: %w", err)
	}

	if err := os.MkdirAll(tokensDir, 0o755); err != nil {
		return fmt.Errorf("creating tokens directory: %w", err)
	}

	// Write private key
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := writePEM(filepath.Join(keysDir, "apiKey.prv"), "RSA PRIVATE KEY", privDER); err != nil {
		return err
	}

	// Write public key
	pubDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshaling public key: %w", err)
	}

	if err := writePEM(filepath.Join(keysDir, "apiKey.pem"), "PUBLIC KEY", pubDER); err != nil {
		return err
	}

	fmt.Println("  keys/apiKey.prv — RSA 4096-bit private key")
	fmt.Println("  keys/apiKey.pem — RSA public key")

	// Generate tokens
	tokens := []tokenSpec{
		{
			filename: "token-bearer-inventory-manager.jwt",
			claims: map[string]any{
				"jti":   "fred",
				"iss":   "example.com",
				"roles": []string{"inventoryManager"},
			},
		},
		{
			filename: "token-apikey-reseller.jwt",
			claims: map[string]any{
				"jti":   "fred",
				"iss":   "example.com",
				"roles": []string{"reseller"},
			},
		},
		{
			filename: "token-apikey-customer.jwt",
			claims: map[string]any{
				"jti":   "ivan",
				"iss":   "example.com",
				"roles": []string{"customer"},
			},
		},
	}

	for _, tok := range tokens {
		jwt, err := signRS256(privateKey, tok.claims)
		if err != nil {
			return fmt.Errorf("signing %s: %w", tok.filename, err)
		}

		path := filepath.Join(tokensDir, tok.filename)
		if err := os.WriteFile(path, []byte(jwt), 0o644); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}

		fmt.Printf("  tokens/%s\n", tok.filename)
	}

	return nil
}

// signRS256 produces a compact RS256 JWT from the given claims.
func signRS256(key *rsa.PrivateKey, claims map[string]any) (string, error) {
	header := map[string]string{"alg": "RS256", "typ": "JWT"}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	enc := base64.RawURLEncoding
	signingInput := enc.EncodeToString(headerJSON) + "." + enc.EncodeToString(claimsJSON)

	hash := sha256.Sum256([]byte(signingInput))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash[:])
	if err != nil {
		return "", err
	}

	return signingInput + "." + enc.EncodeToString(sig), nil
}
