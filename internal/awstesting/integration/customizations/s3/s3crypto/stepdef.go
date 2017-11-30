// +build integration

// Package s3crypto contains shared step definitions that are used across integration tests
package s3crypto

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/gucumber/gucumber"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awstesting/integration"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3crypto"
)

func init() {
	gucumber.When(`^I get all fixtures for "(.+?)" from "(.+?)"$`,
		func(cekAlg, bucket string) {
			prefix := "plaintext_test_case_"
			baseFolder := "crypto_tests/" + cekAlg
			s3Client := gucumber.World["client"].(*s3.S3)

			out, err := s3Client.ListObjectsRequest(&s3.ListObjectsInput{
				Bucket: aws.String(bucket),
				Prefix: aws.String(baseFolder + "/" + prefix),
			}).Send()
			if err != nil {
				gucumber.T.Errorf("expect no error, got %v", err)
			}

			plaintexts := make(map[string][]byte)
			for _, obj := range out.Contents {
				plaintextKey := obj.Key
				ptObj, err := s3Client.GetObjectRequest(&s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    plaintextKey,
				}).Send()
				if err != nil {
					gucumber.T.Errorf("expect no error, got %v", err)
				}
				caseKey := strings.TrimPrefix(*plaintextKey, baseFolder+"/"+prefix)
				plaintext, err := ioutil.ReadAll(ptObj.Body)
				if err != nil {
					gucumber.T.Errorf("expect no error, got %v", err)
				}

				plaintexts[caseKey] = plaintext
			}
			gucumber.World["baseFolder"] = baseFolder
			gucumber.World["bucket"] = bucket
			gucumber.World["plaintexts"] = plaintexts
		})

	gucumber.Then(`^I decrypt each fixture against "(.+?)" "(.+?)"$`, func(lang, version string) {
		plaintexts := gucumber.World["plaintexts"].(map[string][]byte)
		baseFolder := gucumber.World["baseFolder"].(string)
		bucket := gucumber.World["bucket"].(string)
		prefix := "ciphertext_test_case_"
		s3Client := gucumber.World["client"].(*s3.S3)
		s3CryptoClient := gucumber.World["decryptionClient"].(*s3crypto.DecryptionClient)
		language := "language_" + lang

		ciphertexts := make(map[string][]byte)
		for caseKey := range plaintexts {
			cipherKey := baseFolder + "/" + version + "/" + language + "/" + prefix + caseKey

			// To get metadata for encryption key
			ctObj, err := s3Client.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    &cipherKey,
			}).Send()
			if err != nil {
				continue
			}

			// We don't support wrap, so skip it
			if ctObj.Metadata["X-Amz-Wrap-Alg"] == "" || ctObj.Metadata["X-Amz-Wrap-Alg"] != "kms" {
				continue
			}

			ctObj, err = s3CryptoClient.GetObjectRequest(&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    &cipherKey,
			}).Send()
			if err != nil {
				gucumber.T.Errorf("expect no error, got %v", err)
			}

			ciphertext, err := ioutil.ReadAll(ctObj.Body)
			if err != nil {
				gucumber.T.Errorf("expect no error, got %v", err)
			}
			ciphertexts[caseKey] = ciphertext
		}
		gucumber.World["decrypted"] = ciphertexts
	})

	gucumber.And(`^I compare the decrypted ciphertext to the plaintext$`, func() {
		plaintexts := gucumber.World["plaintexts"].(map[string][]byte)
		ciphertexts := gucumber.World["decrypted"].(map[string][]byte)
		for caseKey, ciphertext := range ciphertexts {
			if e, a := len(plaintexts[caseKey]), len(ciphertext); e != a {
				gucumber.T.Errorf("expect %v, got %v", e, a)
			}
			if e, a := plaintexts[caseKey], ciphertext; !bytes.Equal(e, a) {
				gucumber.T.Errorf("expect %v, got %v", e, a)
			}
		}
	})

	gucumber.Then(`^I encrypt each fixture with "(.+?)" "(.+?)" "(.+?)" and "(.+?)"$`, func(kek, v1, v2, cek string) {
		var handler s3crypto.CipherDataGenerator
		var builder s3crypto.ContentCipherBuilder
		switch kek {
		case "kms":
			arn, err := getAliasInformation(v1, v2)
			if err != nil {
				gucumber.T.Errorf("expect nil, got %v", nil)
			}

			b64Arn := base64.StdEncoding.EncodeToString([]byte(arn))
			if err != nil {
				gucumber.T.Errorf("expect nil, got %v", nil)
			}
			gucumber.World["Masterkey"] = b64Arn

			cfg := integration.Config()
			cfg.Region = v2

			handler = s3crypto.NewKMSKeyGenerator(kms.New(cfg), arn)
			if err != nil {
				gucumber.T.Errorf("expect nil, got %v", nil)
			}
		default:
			gucumber.T.Skip()
		}

		switch cek {
		case "aes_gcm":
			builder = s3crypto.AESGCMContentCipherBuilder(handler)
		case "aes_cbc":
			builder = s3crypto.AESCBCContentCipherBuilder(handler, s3crypto.AESCBCPadder)
		default:
			gucumber.T.Skip()
		}

		cfg := integration.Config()
		cfg.Region = "us-west-2"

		c := s3crypto.NewEncryptionClient(cfg, builder, func(c *s3crypto.EncryptionClient) {
		})
		gucumber.World["encryptionClient"] = c
		gucumber.World["cek"] = cek
	})

	gucumber.And(`^upload "(.+?)" data with folder "(.+?)"$`, func(language, folder string) {
		c := gucumber.World["encryptionClient"].(*s3crypto.EncryptionClient)
		cek := gucumber.World["cek"].(string)
		bucket := gucumber.World["bucket"].(string)
		plaintexts := gucumber.World["plaintexts"].(map[string][]byte)
		key := gucumber.World["Masterkey"].(string)
		for caseKey, plaintext := range plaintexts {
			input := &s3.PutObjectInput{
				Bucket: &bucket,
				Key:    aws.String("crypto_tests/" + cek + "/" + folder + "/language_" + language + "/ciphertext_test_case_" + caseKey),
				Body:   bytes.NewReader(plaintext),
				Metadata: map[string]string{
					"Masterkey": key,
				},
			}

			_, err := c.PutObject(input)
			if err != nil {
				gucumber.T.Errorf("expect nil, got %v", nil)
			}
		}
	})
}

func getAliasInformation(alias, region string) (string, error) {
	arn := ""

	cfg := integration.Config()
	cfg.Region = region

	svc := kms.New(cfg)

	truncated := true
	var marker *string
	for truncated {
		out, err := svc.ListAliasesRequest(&kms.ListAliasesInput{
			Marker: marker,
		}).Send()
		if err != nil {
			return arn, err
		}
		for _, aliasEntry := range out.Aliases {
			if *aliasEntry.AliasName == "alias/"+alias {
				return *aliasEntry.AliasArn, nil
			}
		}
		truncated = *out.Truncated
		marker = out.NextMarker
	}

	return "", errors.New("The alias " + alias + " does not exist in your account. Please add the proper alias to a key")
}
