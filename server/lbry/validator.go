package lbry

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"time"

	"github.com/lbryio/commentron/commentapi/lbry"

	"github.com/karlseguin/ccache"

	"github.com/lbryio/commentron/config"

	"github.com/lbryio/commentron/helper"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/jsonrpc"
	"github.com/lbryio/lbry.go/v2/schema/keys"

	"github.com/btcsuite/btcd/btcec"
)

// ValidateSignatures determines if signatures should be validated or not ( not used yet)
var ValidateSignatures bool
var cqPublicKeyCache = ccache.New(ccache.Configure().GetsPerPromote(1).MaxSize(100000))

func getPublicKeyForChannel(channelClaimID string) ([]byte, error) {
	channel, err := SDK.GetClaim(channelClaimID)
	if err != nil {
		pk, err := retrievePKFromCQForChannel(channelClaimID)
		if err != nil {
			return nil, err
		}
		return pk, nil
	}
	return channel.Value.GetChannel().GetPublicKey(), nil
}

type cqResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error"`
	Data    []struct {
		Certificate string `json:"certificate"`
	} `json:"data"`
}

type certificate struct {
	Version   int    `json:"version"`
	KeyType   int    `json:"keyType"`
	PublicKey []byte `json:"publicKey"`
}

func retrievePKFromCQForChannel(channelClaimID string) ([]byte, error) {
	item, err := cqPublicKeyCache.Fetch(channelClaimID, 24*time.Hour, func() (interface{}, error) {
		pk, err := getPublicKeyFromCQForChannel(channelClaimID)
		if err != nil {
			return nil, err
		}
		return pk, err
	})
	if err != nil {
		return nil, err
	}

	v := item.Value()
	pk, ok := v.([]byte)
	if ok {
		return pk, nil
	}

	return nil, errors.Err("could not cast result to byte array from cache")
}

func getPublicKeyFromCQForChannel(channelClaimID string) ([]byte, error) {
	c := http.Client{Timeout: 1 * time.Second}
	sql := fmt.Sprintf(`SELECT certificate FROM claim WHERE claim_id = "%s"`, channelClaimID)
	response, err := c.Get(fmt.Sprintf("https://chainquery.lbry.com/api/sql?query=%s", url.QueryEscape(sql)))
	if err != nil {
		return nil, errors.Err(err)
	}
	defer helper.CloseBody(response.Body)
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Err(err)
	}

	var certResp cqResponse
	err = json.Unmarshal(b, &certResp)
	if err != nil {
		return nil, errors.Err(err)
	}

	if len(certResp.Data) > 0 {
		cert := &certificate{}
		err := json.Unmarshal([]byte(certResp.Data[0].Certificate), cert)
		if err != nil {
			return nil, errors.Err(err)
		}
		return cert.PublicKey, nil
	}
	return nil, errors.Err("no certificate found from CQ")
}

// ValidateSignature validates the signature was signed by the channel reference.
func ValidateSignature(channelClaimID, signature, signingTS, data string) error {
	if config.IsTestMode {
		return nil
	}
	pk, err := getPublicKeyForChannel(channelClaimID)
	if err != nil {
		return err
	}
	return validateSignature(channelClaimID, signature, signingTS, data, pk)
}

// ValidateSignatureFromClaim validates the signature was signed by the channel reference.
func ValidateSignatureFromClaim(channel *jsonrpc.Claim, signature, signingTS, data string) error {
	if channel == nil {
		return errors.Err("no channel to validate")
	}
	if channel.Value.GetChannel() == nil {
		return errors.Err("no channel for public key")
	}
	pk := channel.Value.GetChannel().GetPublicKey()
	return validateSignature(channel.ClaimID, signature, signingTS, data, pk)

}

// encodePrivateKey encodes an ECDSA private key to PEM format.
func encodePrivateKey(key *btcec.PrivateKey) ([]byte, error) {
	derPrivKey, err := keys.PrivateKeyToDER(key)
	if err != nil {
		return nil, err
	}

	keyBlock := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: derPrivKey,
	}

	return pem.EncodeToMemory(keyBlock), nil
}

func validateSignature(channelClaimID, signature, signingTS, data string, pubkey []byte) error {
	publicKey, err := getPublicKeyFromBytes(pubkey)
	if err != nil {
		return errors.Err(err)
	}
	injest := sha256.Sum256(
		helper.CreateDigest(
			[]byte(signingTS),
			lbry.UnhelixifyAndReverse(channelClaimID),
			[]byte(data),
		))
	sig, err := hex.DecodeString(signature)
	if err != nil {
		return errors.Err(err)
	}
	signatureBytes := [64]byte{}
	for i, b := range sig {
		signatureBytes[i] = b
	}
	sigValid := isSignatureValid(signatureBytes, publicKey, injest[:])
	if !sigValid {
		return errors.Err("could not validate the signature")
	}
	return nil
}

func isSignatureValid(signature [64]byte, publicKey *btcec.PublicKey, injest []byte) bool {

	R := &big.Int{}
	S := &big.Int{}
	R.SetBytes(signature[:32])
	S.SetBytes(signature[32:])
	return ecdsa.Verify(publicKey.ToECDSA(), injest, R, S)
}

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

func getPublicKeyFromBytes(pubKeyBytes []byte) (*btcec.PublicKey, error) {
	PKInfo := publicKeyInfo{}
	_, err := asn1.Unmarshal(pubKeyBytes, &PKInfo)
	if err != nil {
		return nil, errors.Err(err)
	}
	pubkeyBytes1 := PKInfo.PublicKey.Bytes
	return btcec.ParsePubKey(pubkeyBytes1, btcec.S256())
}
