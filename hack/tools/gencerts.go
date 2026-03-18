// SPDX-FileCopyrightText: Copyright 2015-2026 go-swagger maintainers
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// runGenCerts generates self-signed TLS certificates for the todo-list-errors example.
//
// It produces:
//   - myCA.key / myCA.crt       — self-signed certificate authority
//   - mycert1.key / mycert1.crt — server certificate (CN=goswagger.local)
//   - myclient.key / myclient.crt — client certificate (CN=localhost)
//
// All certificates use ECDSA P-256 and are valid for 10 years.
func runGenCerts() error {
	outDir := "todo-list-errors"
	if len(os.Args) > 2 {
		outDir = os.Args[2]
	}

	fmt.Printf("Generating TLS certificates in %s/\n", outDir)

	// Generate CA
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("generating CA key: %w", err)
	}

	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "Go Swagger"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	caCertDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("creating CA certificate: %w", err)
	}

	caCert, err := x509.ParseCertificate(caCertDER)
	if err != nil {
		return fmt.Errorf("parsing CA certificate: %w", err)
	}

	if err := writeKeyPair(outDir, "myCA", caKey, caCertDER); err != nil {
		return err
	}

	// Generate server certificate
	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("generating server key: %w", err)
	}

	serverTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "goswagger.local"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"goswagger.local", "localhost", "www.example.com"},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1)},
	}

	serverCertDER, err := x509.CreateCertificate(rand.Reader, serverTemplate, caCert, &serverKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("creating server certificate: %w", err)
	}

	if err := writeKeyPair(outDir, "mycert1", serverKey, serverCertDER); err != nil {
		return err
	}

	// Generate client certificate
	clientKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("generating client key: %w", err)
	}

	clientTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(3),
		Subject: pkix.Name{
			CommonName:   "localhost",
			Country:      []string{"US"},
			Province:     []string{"California"},
			Locality:     []string{"San Francisco"},
			Organization: []string{"go-swagger"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	clientCertDER, err := x509.CreateCertificate(rand.Reader, clientTemplate, caCert, &clientKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("creating client certificate: %w", err)
	}

	if err := writeKeyPair(outDir, "myclient", clientKey, clientCertDER); err != nil {
		return err
	}

	fmt.Println("  myCA.key / myCA.crt       — certificate authority")
	fmt.Println("  mycert1.key / mycert1.crt  — server (CN=goswagger.local)")
	fmt.Println("  myclient.key / myclient.crt — client (CN=localhost)")

	return nil
}

func writeKeyPair(dir, name string, key *ecdsa.PrivateKey, certDER []byte) error {
	keyPath := filepath.Join(dir, name+".key")
	certPath := filepath.Join(dir, name+".crt")

	// Write private key
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return fmt.Errorf("marshaling %s key: %w", name, err)
	}

	if err := writePEM(keyPath, "EC PRIVATE KEY", keyDER); err != nil {
		return err
	}

	// Write certificate
	if err := writePEM(certPath, "CERTIFICATE", certDER); err != nil {
		return err
	}

	return nil
}

func writePEM(path, blockType string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating %s: %w", path, err)
	}
	defer f.Close()

	return pem.Encode(f, &pem.Block{Type: blockType, Bytes: data})
}
