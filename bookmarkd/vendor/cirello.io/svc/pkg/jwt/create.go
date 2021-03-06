// Copyright 2019 github.com/ucirello and https://cirello.io. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

// Package jwt manipulates JWT tokens.
package jwt // import "cirello.io/svc/pkg/jwt"

import (
	"crypto/x509"
	"time"

	"cirello.io/errors"
	jwt "github.com/dgrijalva/jwt-go"
)

// ServiceClaims define the set of claims used by cirello.io services.
type ServiceClaims struct {
	// Email is the actor who is logging in.
	Email string
	// Target defines to which service this token was created for.
	Target string
	// Trust defines the trust level so to give the application some context
	// on how it should handle low-trust logins.
	Trust string

	jwt.StandardClaims
}

// CreateFromCert a JWT whose content indicate a high-trust login.
func CreateFromCert(svcName string, caPEM []byte, cert *x509.Certificate, trustedHost bool) (string, error) {
	if len(cert.EmailAddresses) == 0 {
		return "", errors.E("certificate missing email")
	} else if len(cert.EmailAddresses) > 1 {
		return "", errors.E("multiple emails in the same certificate - cannot choose one")
	}

	trust := "medium"
	if trustedHost {
		trust = "high"
	}
	email := cert.EmailAddresses[0]
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		&ServiceClaims{
			Email:  email,
			Target: svcName,
			Trust:  trust,
		},
	)

	tokenString, err := token.SignedString(caPEM)
	return tokenString, errors.E(err, "cannot sign JWT")
}

// CreateFromEmail a JWT whose content indicate a low-trust login.
func CreateFromEmail(svcName string, caPEM []byte, email string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS512,
		&ServiceClaims{
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(expiration).Unix(),
			},
			Email:  email,
			Target: svcName,
			Trust:  "low",
		},
	)

	tokenString, err := token.SignedString(caPEM)
	return tokenString, errors.E(err, "cannot sign JWT")
}
