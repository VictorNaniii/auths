// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ssh

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"
	"time"

	"golang.org/x/crypto/ssh/testdata"
)

func TestParseCert(t *testing.T) {
	authKeyBytes := bytes.TrimSuffix(testdata.SSHCertificates["rsa"], []byte(" host.example.com\n"))

	key, _, _, rest, err := ParseAuthorizedKey(authKeyBytes)
	if err != nil {
		t.Fatalf("ParseAuthorizedKey: %v", err)
	}
	if len(rest) > 0 {
		t.Errorf("rest: got %q, want empty", rest)
	}

	if _, ok := key.(*Certificate); !ok {
		t.Fatalf("got %v (%T), want *Certificate", key, key)
	}

	marshaled := MarshalAuthorizedKey(key)
	// Before comparison, remove the trailing newline that
	// MarshalAuthorizedKey adds.
	marshaled = marshaled[:len(marshaled)-1]
	if !bytes.Equal(authKeyBytes, marshaled) {
		t.Errorf("marshaled certificate does not match original: got %q, want %q", marshaled, authKeyBytes)
	}
}

// Cert generated by ssh-keygen OpenSSH_6.8p1 OS X 10.10.3
// % ssh-keygen -s ca -I testcert -O source-address=192.168.1.0/24 -O force-command=/bin/sleep user.pub
// user.pub key: ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDACh1rt2DXfV3hk6fszSQcQ/rueMId0kVD9U7nl8cfEnFxqOCrNT92g4laQIGl2mn8lsGZfTLg8ksHq3gkvgO3oo/0wHy4v32JeBOHTsN5AL4gfHNEhWeWb50ev47hnTsRIt9P4dxogeUo/hTu7j9+s9lLpEQXCvq6xocXQt0j8MV9qZBBXFLXVT3cWIkSqOdwt/5ZBg+1GSrc7WfCXVWgTk4a20uPMuJPxU4RQwZW6X3+O8Pqo8C3cW0OzZRFP6gUYUKUsTI5WntlS+LAxgw1mZNsozFGdbiOPRnEryE3SRldh9vjDR3tin1fGpA5P7+CEB/bqaXtG3V+F2OkqaMN
// Critical Options:
//
//	force-command /bin/sleep
//	source-address 192.168.1.0/24
//
// Extensions:
//
//	permit-X11-forwarding
//	permit-agent-forwarding
//	permit-port-forwarding
//	permit-pty
//	permit-user-rc
const exampleSSHCertWithOptions = `ssh-rsa-cert-v01@openssh.com AAAAHHNzaC1yc2EtY2VydC12MDFAb3BlbnNzaC5jb20AAAAgDyysCJY0XrO1n03EeRRoITnTPdjENFmWDs9X58PP3VUAAAADAQABAAABAQDACh1rt2DXfV3hk6fszSQcQ/rueMId0kVD9U7nl8cfEnFxqOCrNT92g4laQIGl2mn8lsGZfTLg8ksHq3gkvgO3oo/0wHy4v32JeBOHTsN5AL4gfHNEhWeWb50ev47hnTsRIt9P4dxogeUo/hTu7j9+s9lLpEQXCvq6xocXQt0j8MV9qZBBXFLXVT3cWIkSqOdwt/5ZBg+1GSrc7WfCXVWgTk4a20uPMuJPxU4RQwZW6X3+O8Pqo8C3cW0OzZRFP6gUYUKUsTI5WntlS+LAxgw1mZNsozFGdbiOPRnEryE3SRldh9vjDR3tin1fGpA5P7+CEB/bqaXtG3V+F2OkqaMNAAAAAAAAAAAAAAABAAAACHRlc3RjZXJ0AAAAAAAAAAAAAAAA//////////8AAABLAAAADWZvcmNlLWNvbW1hbmQAAAAOAAAACi9iaW4vc2xlZXAAAAAOc291cmNlLWFkZHJlc3MAAAASAAAADjE5Mi4xNjguMS4wLzI0AAAAggAAABVwZXJtaXQtWDExLWZvcndhcmRpbmcAAAAAAAAAF3Blcm1pdC1hZ2VudC1mb3J3YXJkaW5nAAAAAAAAABZwZXJtaXQtcG9ydC1mb3J3YXJkaW5nAAAAAAAAAApwZXJtaXQtcHR5AAAAAAAAAA5wZXJtaXQtdXNlci1yYwAAAAAAAAAAAAABFwAAAAdzc2gtcnNhAAAAAwEAAQAAAQEAwU+c5ui5A8+J/CFpjW8wCa52bEODA808WWQDCSuTG/eMXNf59v9Y8Pk0F1E9dGCosSNyVcB/hacUrc6He+i97+HJCyKavBsE6GDxrjRyxYqAlfcOXi/IVmaUGiO8OQ39d4GHrjToInKvExSUeleQyH4Y4/e27T/pILAqPFL3fyrvMLT5qU9QyIt6zIpa7GBP5+urouNavMprV3zsfIqNBbWypinOQAw823a5wN+zwXnhZrgQiHZ/USG09Y6k98y1dTVz8YHlQVR4D3lpTAsKDKJ5hCH9WU4fdf+lU8OyNGaJ/vz0XNqxcToe1l4numLTnaoSuH89pHryjqurB7lJKwAAAQ8AAAAHc3NoLXJzYQAAAQCaHvUIoPL1zWUHIXLvu96/HU1s/i4CAW2IIEuGgxCUCiFj6vyTyYtgxQxcmbfZf6eaITlS6XJZa7Qq4iaFZh75C1DXTX8labXhRSD4E2t//AIP9MC1rtQC5xo6FmbQ+BoKcDskr+mNACcbRSxs3IL3bwCfWDnIw2WbVox9ZdcthJKk4UoCW4ix4QwdHw7zlddlz++fGEEVhmTbll1SUkycGApPFBsAYRTMupUJcYPIeReBI/m8XfkoMk99bV8ZJQTAd7OekHY2/48Ff53jLmyDjP7kNw1F8OaPtkFs6dGJXta4krmaekPy87j+35In5hFj7yoOqvSbmYUkeX70/GGQ`

func TestParseCertWithOptions(t *testing.T) {
	opts := map[string]string{
		"source-address": "192.168.1.0/24",
		"force-command":  "/bin/sleep",
	}
	exts := map[string]string{
		"permit-X11-forwarding":   "",
		"permit-agent-forwarding": "",
		"permit-port-forwarding":  "",
		"permit-pty":              "",
		"permit-user-rc":          "",
	}
	authKeyBytes := []byte(exampleSSHCertWithOptions)

	key, _, _, rest, err := ParseAuthorizedKey(authKeyBytes)
	if err != nil {
		t.Fatalf("ParseAuthorizedKey: %v", err)
	}
	if len(rest) > 0 {
		t.Errorf("rest: got %q, want empty", rest)
	}
	cert, ok := key.(*Certificate)
	if !ok {
		t.Fatalf("got %v (%T), want *Certificate", key, key)
	}
	if !reflect.DeepEqual(cert.CriticalOptions, opts) {
		t.Errorf("unexpected critical options - got %v, want %v", cert.CriticalOptions, opts)
	}
	if !reflect.DeepEqual(cert.Extensions, exts) {
		t.Errorf("unexpected Extensions - got %v, want %v", cert.Extensions, exts)
	}
	marshaled := MarshalAuthorizedKey(key)
	// Before comparison, remove the trailing newline that
	// MarshalAuthorizedKey adds.
	marshaled = marshaled[:len(marshaled)-1]
	if !bytes.Equal(authKeyBytes, marshaled) {
		t.Errorf("marshaled certificate does not match original: got %q, want %q", marshaled, authKeyBytes)
	}
}

func TestValidateCert(t *testing.T) {
	key, _, _, _, err := ParseAuthorizedKey(testdata.SSHCertificates["rsa-user-testcertificate"])
	if err != nil {
		t.Fatalf("ParseAuthorizedKey: %v", err)
	}
	validCert, ok := key.(*Certificate)
	if !ok {
		t.Fatalf("got %v (%T), want *Certificate", key, key)
	}
	checker := CertChecker{}
	checker.IsUserAuthority = func(k PublicKey) bool {
		return bytes.Equal(k.Marshal(), validCert.SignatureKey.Marshal())
	}

	if err := checker.CheckCert("testcertificate", validCert); err != nil {
		t.Errorf("Unable to validate certificate: %v", err)
	}
	invalidCert := &Certificate{
		Key:          testPublicKeys["rsa"],
		SignatureKey: testPublicKeys["ecdsa"],
		ValidBefore:  CertTimeInfinity,
		Signature:    &Signature{},
	}
	if err := checker.CheckCert("testcertificate", invalidCert); err == nil {
		t.Error("Invalid cert signature passed validation")
	}
}

func TestValidateCertTime(t *testing.T) {
	cert := Certificate{
		ValidPrincipals: []string{"user"},
		Key:             testPublicKeys["rsa"],
		ValidAfter:      50,
		ValidBefore:     100,
	}

	cert.SignCert(rand.Reader, testSigners["ecdsa"])

	for ts, ok := range map[int64]bool{
		25:  false,
		50:  true,
		99:  true,
		100: false,
		125: false,
	} {
		checker := CertChecker{
			Clock: func() time.Time { return time.Unix(ts, 0) },
		}
		checker.IsUserAuthority = func(k PublicKey) bool {
			return bytes.Equal(k.Marshal(),
				testPublicKeys["ecdsa"].Marshal())
		}

		if v := checker.CheckCert("user", &cert); (v == nil) != ok {
			t.Errorf("Authenticate(%d): %v", ts, v)
		}
	}
}

// TODO(hanwen): tests for
//
// host keys:
// * fallbacks

func TestHostKeyCert(t *testing.T) {
	cert := &Certificate{
		ValidPrincipals: []string{"hostname", "hostname.domain", "otherhost"},
		Key:             testPublicKeys["rsa"],
		ValidBefore:     CertTimeInfinity,
		CertType:        HostCert,
	}
	cert.SignCert(rand.Reader, testSigners["ecdsa"])

	checker := &CertChecker{
		IsHostAuthority: func(p PublicKey, addr string) bool {
			return addr == "hostname:22" && bytes.Equal(testPublicKeys["ecdsa"].Marshal(), p.Marshal())
		},
	}

	certSigner, err := NewCertSigner(cert, testSigners["rsa"])
	if err != nil {
		t.Errorf("NewCertSigner: %v", err)
	}

	for _, test := range []struct {
		addr                    string
		succeed                 bool
		certSignerAlgorithms    []string // Empty means no algorithm restrictions.
		clientHostKeyAlgorithms []string
	}{
		{addr: "hostname:22", succeed: true},
		{
			addr:                    "hostname:22",
			succeed:                 true,
			certSignerAlgorithms:    []string{KeyAlgoRSASHA256, KeyAlgoRSASHA512},
			clientHostKeyAlgorithms: []string{CertAlgoRSASHA512v01},
		},
		{
			addr:                    "hostname:22",
			succeed:                 false,
			certSignerAlgorithms:    []string{KeyAlgoRSASHA256, KeyAlgoRSASHA512},
			clientHostKeyAlgorithms: []string{CertAlgoRSAv01},
		},
		{
			addr:                    "hostname:22",
			succeed:                 false,
			certSignerAlgorithms:    []string{KeyAlgoRSASHA256, KeyAlgoRSASHA512},
			clientHostKeyAlgorithms: []string{KeyAlgoRSASHA512}, // Not a certificate algorithm.
		},
		{addr: "otherhost:22", succeed: false}, // The certificate is valid for 'otherhost' as hostname, but we only recognize the authority of the signer for the address 'hostname:22'
		{addr: "lasthost:22", succeed: false},
	} {
		c1, c2, err := netPipe()
		if err != nil {
			t.Fatalf("netPipe: %v", err)
		}
		defer c1.Close()
		defer c2.Close()

		errc := make(chan error)

		go func() {
			conf := ServerConfig{
				NoClientAuth: true,
			}
			if len(test.certSignerAlgorithms) > 0 {
				mas, err := NewSignerWithAlgorithms(certSigner.(AlgorithmSigner), test.certSignerAlgorithms)
				if err != nil {
					errc <- err
					return
				}
				conf.AddHostKey(mas)
			} else {
				conf.AddHostKey(certSigner)
			}
			_, _, _, err := NewServerConn(c1, &conf)
			errc <- err
		}()

		config := &ClientConfig{
			User:              "user",
			HostKeyCallback:   checker.CheckHostKey,
			HostKeyAlgorithms: test.clientHostKeyAlgorithms,
		}
		_, _, _, err = NewClientConn(c2, test.addr, config)

		if (err == nil) != test.succeed {
			t.Errorf("NewClientConn(%q): %v", test.addr, err)
		}

		err = <-errc
		if (err == nil) != test.succeed {
			t.Errorf("NewServerConn(%q): %v", test.addr, err)
		}
	}
}

type legacyRSASigner struct {
	Signer
}

func (s *legacyRSASigner) Sign(rand io.Reader, data []byte) (*Signature, error) {
	v, ok := s.Signer.(AlgorithmSigner)
	if !ok {
		return nil, fmt.Errorf("invalid signer")
	}
	return v.SignWithAlgorithm(rand, data, KeyAlgoRSA)
}

func TestCertTypes(t *testing.T) {
	algorithmSigner, ok := testSigners["rsa"].(AlgorithmSigner)
	if !ok {
		t.Fatal("rsa test signer does not implement the AlgorithmSigner interface")
	}
	multiAlgoSignerSHA256, err := NewSignerWithAlgorithms(algorithmSigner, []string{KeyAlgoRSASHA256})
	if err != nil {
		t.Fatalf("unable to create multi algorithm signer SHA256: %v", err)
	}
	// Algorithms are in order of preference, we expect rsa-sha2-512 to be used.
	multiAlgoSignerSHA512, err := NewSignerWithAlgorithms(algorithmSigner, []string{KeyAlgoRSASHA512, KeyAlgoRSASHA256})
	if err != nil {
		t.Fatalf("unable to create multi algorithm signer SHA512: %v", err)
	}

	var testVars = []struct {
		name   string
		signer Signer
		algo   string
	}{
		{CertAlgoECDSA256v01, testSigners["ecdsap256"], ""},
		{CertAlgoECDSA384v01, testSigners["ecdsap384"], ""},
		{CertAlgoECDSA521v01, testSigners["ecdsap521"], ""},
		{CertAlgoED25519v01, testSigners["ed25519"], ""},
		{CertAlgoRSAv01, testSigners["rsa"], KeyAlgoRSASHA256},
		{"legacyRSASigner", &legacyRSASigner{testSigners["rsa"]}, KeyAlgoRSA},
		{"multiAlgoRSASignerSHA256", multiAlgoSignerSHA256, KeyAlgoRSASHA256},
		{"multiAlgoRSASignerSHA512", multiAlgoSignerSHA512, KeyAlgoRSASHA512},
		{InsecureCertAlgoDSAv01, testSigners["dsa"], ""},
	}

	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("error generating host key: %v", err)
	}

	signer, err := NewSignerFromKey(k)
	if err != nil {
		t.Fatalf("error generating signer for ssh listener: %v", err)
	}

	conf := &ServerConfig{
		PublicKeyCallback: func(c ConnMetadata, k PublicKey) (*Permissions, error) {
			return new(Permissions), nil
		},
	}
	conf.AddHostKey(signer)

	for _, m := range testVars {
		t.Run(m.name, func(t *testing.T) {

			c1, c2, err := netPipe()
			if err != nil {
				t.Fatalf("netPipe: %v", err)
			}
			defer c1.Close()
			defer c2.Close()

			go NewServerConn(c1, conf)

			priv := m.signer
			if err != nil {
				t.Fatalf("error generating ssh pubkey: %v", err)
			}

			cert := &Certificate{
				CertType: UserCert,
				Key:      priv.PublicKey(),
			}
			cert.SignCert(rand.Reader, priv)

			certSigner, err := NewCertSigner(cert, priv)
			if err != nil {
				t.Fatalf("error generating cert signer: %v", err)
			}

			if m.algo != "" && cert.Signature.Format != m.algo {
				t.Errorf("expected %q signature format, got %q", m.algo, cert.Signature.Format)
			}

			config := &ClientConfig{
				User:            "user",
				HostKeyCallback: func(h string, r net.Addr, k PublicKey) error { return nil },
				Auth:            []AuthMethod{PublicKeys(certSigner)},
			}

			_, _, _, err = NewClientConn(c2, "", config)
			if err != nil {
				t.Fatalf("error connecting: %v", err)
			}
		})
	}
}

func TestCertSignWithMultiAlgorithmSigner(t *testing.T) {
	type testcase struct {
		sigAlgo    string
		algorithms []string
	}
	cases := []testcase{
		{
			sigAlgo:    KeyAlgoRSA,
			algorithms: []string{KeyAlgoRSA, KeyAlgoRSASHA512},
		},
		{
			sigAlgo:    KeyAlgoRSASHA256,
			algorithms: []string{KeyAlgoRSASHA256, KeyAlgoRSA, KeyAlgoRSASHA512},
		},
		{
			sigAlgo:    KeyAlgoRSASHA512,
			algorithms: []string{KeyAlgoRSASHA512, KeyAlgoRSASHA256},
		},
	}

	cert := &Certificate{
		Key:         testPublicKeys["rsa"],
		ValidBefore: CertTimeInfinity,
		CertType:    UserCert,
	}

	for _, c := range cases {
		t.Run(c.sigAlgo, func(t *testing.T) {
			signer, err := NewSignerWithAlgorithms(testSigners["rsa"].(AlgorithmSigner), c.algorithms)
			if err != nil {
				t.Fatalf("NewSignerWithAlgorithms error: %v", err)
			}
			if err := cert.SignCert(rand.Reader, signer); err != nil {
				t.Fatalf("SignCert error: %v", err)
			}
			if cert.Signature.Format != c.sigAlgo {
				t.Fatalf("got signature format %q, want %q", cert.Signature.Format, c.sigAlgo)
			}
		})
	}
}

func TestCertSignWithCertificate(t *testing.T) {
	cert := &Certificate{
		Key:         testPublicKeys["rsa"],
		ValidBefore: CertTimeInfinity,
		CertType:    UserCert,
	}
	if err := cert.SignCert(rand.Reader, testSigners["ecdsa"]); err != nil {
		t.Fatalf("SignCert: %v", err)
	}
	signer, err := NewSignerWithAlgorithms(testSigners["rsa"].(AlgorithmSigner), []string{KeyAlgoRSASHA256})
	if err != nil {
		t.Fatal(err)
	}
	certSigner, err := NewCertSigner(cert, signer)
	if err != nil {
		t.Fatalf("NewCertSigner: %v", err)
	}

	cert1 := &Certificate{
		Key:         testPublicKeys["ecdsa"],
		ValidBefore: CertTimeInfinity,
		CertType:    UserCert,
	}

	if err := cert1.SignCert(rand.Reader, certSigner); err == nil {
		t.Fatal("successfully signed a certificate using another certificate, it is expected to fail")
	}
}
