package jwt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrNotFound = errors.New("the requested resource was not found")
)

// JWK and JWKS are documented here by Auth0:
// https://auth0.com/docs/secure/tokens/json-web-tokens/json-web-key-sets

// JWKS is a JSON Web Key Set
type JWKS struct {
	Keys []JWK `json:"keys"`
}

func (j *JWKS) GetKey(keyId string) ([]byte, error) {
	for i := range j.Keys {
		if j.Keys[i].KId == keyId {

			d, err := json.Marshal(j.Keys[i])
			if err != nil {
				return nil, fmt.Errorf("failed to marshal key: %w", err)
			}
			return d, nil
		}
	}

	return nil, ErrNotFound
}

// JWK is a Json Web Key
type JWK struct {
	// The specific cryptographic algorithm used with the key. This is an optional parameter. By default, Auth0 includes the signing algorithm defined at the tenant level in the JSON Web Key Set (JWKS), which is then published. To allow for keys to be used with multiple algorithms rather than a single algorithm i.e. RS256, toggle off Include Signing Algorithms in JSON Web Key Set under Advanced Tenant settings in Dashboard. This removes the alg parameter and requires consumers of the JWKS to interpret the signing algorithms as needed.
	// ex. "RS256"
	Alg string `json:"alg"`

	// The family of cryptographic algorithms used with the key. Ex. "RSA"
	Kty string `json:"kty"`

	//How the key was meant to be used; "sig" represents the signature.
	Use string `json:"use"`

	//The x.509 certificate chain. The first entry in the array is the certificate to use for token verification; the other certificates can be used to verify this first certificate.
	X5c []string `json:"x5c"`

	// The modulus for the RSA public key.
	N string `json:"n"`

	// The exponent for the RSA public key.
	E string `json:"e"`

	// The unique identifier for the key.
	KId string `json:"kid"`

	// The thumbprint of the x.509 cert (SHA-1 thumbprint).
	X5t string `json:"x5t"`
}

func (j *JWK) ToBytes() ([]byte, error) {
	return json.Marshal(j)
}

func ConvertPublicKeyPEMToJWK(KeyId string, key []byte) (JWK, error) {
	k, err := jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		return JWK{}, fmt.Errorf("failed to parse public key from pem: %w", err)
	}

	n := base64.RawURLEncoding.EncodeToString((*k.N).Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(k.E)).Bytes())

	return JWK{
		Kty: "RSA",
		KId: KeyId,
		Use: "sig",
		N:   n,
		E:   e,
	}, nil

}
