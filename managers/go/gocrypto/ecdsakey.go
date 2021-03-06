package gocrypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"

	"github.com/mariotoffia/goservice/interfaces/ifcrypto"
	"github.com/mariotoffia/goservice/utils/cryptoutils"
)

// ECDSAPrivateKey implements the `ifcrypto.KeyPair` interface for a `*rsa.PrivateKey`.
type ECDSAPrivateKey struct {
	KeyBase
	key    *ecdsa.PrivateKey
	public *ECDSAPublicKey
}

// NewECDSAPrivateKeyFromKey creates a new `ECDSAPrivateKey`
//
// The public key portion derives the same usage as the private key
func NewECDSAPrivateKeyFromKey(
	id string,
	key *ecdsa.PrivateKey,
	usage ...ifcrypto.KeyUsage,
) *ECDSAPrivateKey {

	return &ECDSAPrivateKey{
		KeyBase: KeyBase{
			id:      id,
			keyType: ifcrypto.KeyTypeRsa,
			keySize: key.Params().BitSize,
			usage:   usage,
			chiper:  []ifcrypto.Chipher{},
		},
		key:    key,
		public: NewECDSAPublicKeyFromKey(id, &key.PublicKey, usage...),
	}

}

// NewECDSAPrivateKeyFromPEM initializes a new `*ecdsa.PrivateKey` from the underlying _PEM_ block.
func NewECDSAPrivateKeyFromPEM(
	block pem.Block,
	id string,
	usage ...ifcrypto.KeyUsage,
) (*ECDSAPrivateKey, error) {

	if block.Type == "PRIVATE KEY" {
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)

		if err != nil {
			return nil, err
		}

		if ecdsakey, ok := key.(*ecdsa.PrivateKey); ok {

			return NewECDSAPrivateKeyFromKey(id, ecdsakey, usage...), nil

		}

		return nil, fmt.Errorf("not a *ecdsa.PrivateKey: %T", key)

	}

	if block.Type == "EC PRIVATE KEY" {

		key, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		return NewECDSAPrivateKeyFromKey(id, key, usage...), nil

	}

	return nil, fmt.Errorf("unsupported PEM block: %s", block.Type)

}

// NewECDSAPrivateKey generates a new `ECDSAPrivateKey` using the `rand.Reader` as entropy.
func NewECDSAPrivateKey(id string, bits int, usage ...ifcrypto.KeyUsage) (*ECDSAPrivateKey, error) {

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	return NewECDSAPrivateKeyFromKey(id, key, usage...), nil
}

// Sign implements the `crypto.Signer` _interface_. The _opts_
// arguments is not used.
func (r *ECDSAPrivateKey) Sign(
	rand io.Reader,
	digest []byte,
	opts crypto.SignerOpts,
) ([]byte, error) {

	return r.key.Sign(rand, digest, opts)

}

// GetPublic returns the public portion of the key
func (r *ECDSAPrivateKey) GetPublic() ifcrypto.PublicKey {
	return r.public
}

// PEMWrite will write the key onto _w_.
//
// If private key, and _public_ is `true`, it will in addition write the public portion as well.
func (r *ECDSAPrivateKey) PEMWrite(w io.Writer, public bool) error {

	return cryptoutils.ECDSAPrivateKeyToPEM(w, r.key, public)

}

// GetKey gets the underlying key, if any.
//
// Some keys are remote and not possible to fetch. In such situations the function returns a remote id,
// most often the same as GetID() returns.
func (r *ECDSAPrivateKey) GetKey() interface{} {
	return r.key
}

// IsSymmetric returns `true` if this is a `KeyTypeSymmetric`
//
// This is a convenience function instead of `GetKeyType`.
func (r *ECDSAPrivateKey) IsSymmetric() bool {
	return false
}

// IsPrivate returns `true` if this is a `KeyType` other than `KeyTypeSymmetric` and is a private key.
//
// If `KeyTypeSymmetric` it will return `true` since all symmetric keys are considered as private.
func (r *ECDSAPrivateKey) IsPrivate() bool {
	return true
}

// IsRemoteKey returns `true` if the key is not present in current process memory.
//
// Typically hardware units or remote services will not reveal their private key. In such case, this
// method returns `true`. If present in memory such as a `*rsa.PrivateKey` it returns `false`.
func (r *ECDSAPrivateKey) IsRemoteKey() bool {
	return false
}

// ECDSAPublicKey implements the `ifcrypto.PublicKey` interface for `*rsa.PublicKey`
type ECDSAPublicKey struct {
	KeyBase
	key *ecdsa.PublicKey
}

// NewECDSAPublicKeyFromKey creates a instance based on a existing public key.
func NewECDSAPublicKeyFromKey(
	id string,
	key *ecdsa.PublicKey,
	usage ...ifcrypto.KeyUsage,
) *ECDSAPublicKey {

	return &ECDSAPublicKey{
		KeyBase: KeyBase{
			id:      id,
			keyType: ifcrypto.KeyTypeRsa,
			keySize: key.Params().BitSize,
			usage:   usage,
		},
		key: key,
	}

}

// NewECDSAPublicKeyFromPEM initializes a new `*ecdsa.PublicKey` from the underlying _PEM_ block.
func NewECDSAPublicKeyFromPEM(
	block pem.Block,
	id string,
	usage ...ifcrypto.KeyUsage,
) (*ECDSAPublicKey, error) {

	if block.Type == "PUBLIC KEY" || block.Type == "EC PUBLIC KEY" {

		key, err := x509.ParsePKIXPublicKey(block.Bytes)

		if err != nil {
			return nil, err
		}

		if ecdsakey, ok := key.(*ecdsa.PublicKey); ok {

			return NewECDSAPublicKeyFromKey(id, ecdsakey, usage...), nil

		}

		return nil, fmt.Errorf("not a *ecdsa.PublicKey: %T", key)

	}

	return nil, fmt.Errorf("unsupported PEM block: %s", block.Type)

}

// PEMWrite will write the key onto _w_.
//
// Since this is a public key, it will ignore the _public_ parameter.
func (r *ECDSAPublicKey) PEMWrite(w io.Writer, public bool) error {

	return cryptoutils.ECDSAPublicKeyToPEM(w, r.key)

}

// GetKey gets the underlying key, if any.
//
// Some keys are remote and not possible to fetch. In such situations the function returns a remote id,
// most often the same as GetID() returns.
func (r *ECDSAPublicKey) GetKey() interface{} {
	return r.key
}

// IsSymmetric returns `true` if this is a `KeyTypeSymmetric`
//
// This is a convenience function instead of `GetKeyType`.
func (r *ECDSAPublicKey) IsSymmetric() bool {
	return false
}

// IsPrivate returns `true` if this is a `KeyType` other than `KeyTypeSymmetric` and is a private key.
//
// If `KeyTypeSymmetric` it will return `true` since all symmetric keys are considered as private.
func (r *ECDSAPublicKey) IsPrivate() bool {
	return true
}

// IsRemoteKey returns `true` if the key is not present in current process memory.
//
// Typically hardware units or remote services will not reveal their private key. In such case, this
// method returns `true`. If present in memory such as a `*rsa.PrivateKey` it returns `false`.
func (r *ECDSAPublicKey) IsRemoteKey() bool {
	return false
}
