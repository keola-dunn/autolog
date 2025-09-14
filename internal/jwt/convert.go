package jwt

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/lestrrat-go/jwx/jwk"
)

func ConvertPublicKeyPEMToJWK(KeyId string, key *rsa.PublicKey) (JWK, error) {
	n := base64.RawURLEncoding.EncodeToString((*key.N).Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes())

	return JWK{
		Kty: "RSA",
		KId: KeyId,
		Use: "sig",
		N:   n,
		E:   e,
	}, nil

}

func ConvertJWKToPEM(ctx context.Context, j JWK) (*rsa.PublicKey, error) {
	d, err := j.ToBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get jwk bytes: %w", err)
	}

	set, err := jwk.Parse([]byte(d))
	if err != nil {
		return nil, fmt.Errorf("failed to parse jwk: %w", err)
	}

	for it := set.Iterate(context.Background()); it.Next(context.Background()); {
		// this should be a single key
		// i don't love this library, so using this based on examples to not have to
		// re-work the inner JWT-PEM conversions myself.
		pair := it.Pair()
		key := pair.Value.(jwk.Key)

		var rawkey interface{} // This is the raw key, like *rsa.PrivateKey or *ecdsa.PrivateKey
		if err := key.Raw(&rawkey); err != nil {
			return nil, fmt.Errorf("failed to convert the key to expected: %w", err)
		}

		// We know this is an RSA Key
		rsa, ok := rawkey.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("failed to cast key as PublicKey: %w", err)
		}
		return rsa, nil
	}

	return nil, fmt.Errorf("unexpected error")
}
