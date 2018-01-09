package s3crypto

import (
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/kmsiface"
)

const (
	// KMSWrap is a constant used during decryption to build a KMS key handler.
	KMSWrap = "kms"
)

// kmsKeyHandler will make calls to KMS to get the masterkey
type kmsKeyHandler struct {
	kms   kmsiface.KMSAPI
	cmkID string

	CipherData
}

// NewKMSKeyGenerator builds a new KMS key provider using the customer key ID and material
// description.
//
// Example:
//  cfg, err := external.LoadDefaultAWSConfig()
//	cmkID := "arn to key"
//	matdesc := s3crypto.MaterialDescription{}
//	handler := s3crypto.NewKMSKeyGenerator(kms.New(cfg), cmkID)
func NewKMSKeyGenerator(kmsClient kmsiface.KMSAPI, cmkID string) CipherDataGenerator {
	return NewKMSKeyGeneratorWithMatDesc(kmsClient, cmkID, MaterialDescription{})
}

// NewKMSKeyGeneratorWithMatDesc builds a new KMS key provider using the customer key ID and material
// description.
//
// Example:
//  cfg, err := external.LoadDefaultAWSConfig()
//	cmkID := "arn to key"
//	matdesc := s3crypto.MaterialDescription{}
//	handler, err := s3crypto.NewKMSKeyGeneratorWithMatDesc(kms.New(cfg), cmkID, matdesc)
func NewKMSKeyGeneratorWithMatDesc(kmsClient kmsiface.KMSAPI, cmkID string, matdesc MaterialDescription) CipherDataGenerator {
	if matdesc == nil {
		matdesc = MaterialDescription{}
	}
	matdesc["kms_cmk_id"] = cmkID

	// These values are read only making them thread safe
	kp := &kmsKeyHandler{
		kms:   kmsClient,
		cmkID: cmkID,
	}
	// These values are read only making them thread safe
	kp.CipherData.WrapAlgorithm = KMSWrap
	kp.CipherData.MaterialDescription = matdesc
	return kp
}

// decryptHandler initializes a KMS keyprovider with a material description. This
// is used with Decrypting kms content, due to the cmkID being in the material description.
func (kp kmsKeyHandler) decryptHandler(env Envelope) (CipherDataDecrypter, error) {
	m := MaterialDescription{}
	err := m.decodeDescription([]byte(env.MatDesc))
	if err != nil {
		return nil, err
	}

	cmkID, ok := m["kms_cmk_id"]
	if !ok {
		return nil, awserr.New("MissingCMKIDError", "Material description is missing CMK ID", nil)
	}

	kp.CipherData.MaterialDescription = m
	kp.cmkID = cmkID
	kp.WrapAlgorithm = KMSWrap
	return &kp, nil
}

// DecryptKey makes a call to KMS to decrypt the key.
func (kp *kmsKeyHandler) DecryptKey(key []byte) ([]byte, error) {
	req := kp.kms.DecryptRequest(&kms.DecryptInput{
		EncryptionContext: map[string]string(kp.CipherData.MaterialDescription),
		CiphertextBlob:    key,
		GrantTokens:       []string{},
	})
	resp, err := req.Send()
	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil
}

// GenerateCipherData makes a call to KMS to generate a data key, Upon making
// the call, it also sets the encrypted key.
func (kp *kmsKeyHandler) GenerateCipherData(keySize, ivSize int) (CipherData, error) {
	req := kp.kms.GenerateDataKeyRequest(&kms.GenerateDataKeyInput{
		EncryptionContext: kp.CipherData.MaterialDescription,
		KeyId:             &kp.cmkID,
		KeySpec:           kms.DataKeySpecAes256,
	})
	resp, err := req.Send()
	if err != nil {
		return CipherData{}, err
	}

	iv := generateBytes(ivSize)
	cd := CipherData{
		Key:                 resp.Plaintext,
		IV:                  iv,
		WrapAlgorithm:       KMSWrap,
		MaterialDescription: kp.CipherData.MaterialDescription,
		EncryptedKey:        resp.CiphertextBlob,
	}
	return cd, nil
}
