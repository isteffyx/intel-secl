/*
 * Copyright (C) 2021 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"github.com/intel-secl/intel-secl/v4/pkg/lib/common/crypt"
	cLog "github.com/intel-secl/intel-secl/v4/pkg/lib/common/log"
	cMsg "github.com/intel-secl/intel-secl/v4/pkg/lib/common/log/message"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"unsafe"
)

var log = cLog.GetDefaultLogger()

func Encrypt(imagePath string, privateKeyLocation string, encryptedFileLocation string, wrappedKey []byte) error {
	log.Trace("pkg/wpm/util/encrypt.go:Encrypt() Entering")
	defer log.Trace("pkg/wpm/util/encrypt.go:Encrypt() Leaving")

	var encryptionHeader crypt.EncryptionHeader

	// reading image file
	image, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return errors.Wrap(err, "Error reading the image file")
	}

	key, err := UnwrapKey(wrappedKey, privateKeyLocation)
	if err != nil {
		return errors.Wrap(err, "Error while unwrapping the key")
	}
	// creating a new cipher block of 128 bits
	block, err := aes.NewCipher(key)
	if err != nil {
		return errors.Wrap(err, "Error initializing cipher")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return errors.Wrap(err, "Error creating a cipher block")
	}

	log.Infof("pkg/util/encrypt.go:Encrypt() %s", cMsg.EncKeyUsed)

	// assigning a 12 byte empty array to store random value
	iv := make([]byte, 12)
	// reading random value into the byte array
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return errors.Wrap(err, "Error creating random IV value")
	}

	copy(encryptionHeader.MagicText[:], crypt.EncryptionHeaderMagicText)
	copy(encryptionHeader.EncryptionAlgorithm[:], crypt.GCMEncryptionAlgorithm)
	copy(encryptionHeader.IV[:], iv)
	copy(encryptionHeader.Version[:], crypt.EncryptionHeaderVersion)
	encryptionHeader.OffsetInLittleEndian = uint32(unsafe.Sizeof(encryptionHeader))

	encryptionHeaderSlice := &bytes.Buffer{}
	err = binary.Write(encryptionHeaderSlice, binary.LittleEndian, encryptionHeader)
	if err != nil {
		return errors.Wrap(err, "Error while writing encryption header struc values in to buffer")
	}

	// The first 44 bytes of the encrypted file is the encryption header and
	// the rest is the data.
	encryptedDataWithHeader := gcm.Seal(encryptionHeaderSlice.Bytes(), iv, image, nil)
	err = ioutil.WriteFile(encryptedFileLocation, encryptedDataWithHeader, 0600)
	if err != nil {
		return errors.Wrap(err, "Error during writing the encrypted image to file")
	}

	log.Info("pkg/wpm/util/encrypt.go:Encrypt() Successfully encrypted image")
	return nil
}

func UnwrapKey(wrappedKey []byte, privateKeyLocation string) ([]byte, error) {
	log.Trace("pkg/wpm/util/encrypt.go:UnwrapKey() Entering")
	defer log.Trace("pkg/wpm/util/encrypt.go:UnwrapKey() Leaving")

	privateKey, err := ioutil.ReadFile(privateKeyLocation)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading private envelope key file")
	}

	privateKeyBlock, _ := pem.Decode(privateKey)
	var pri *rsa.PrivateKey
	pri, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "Error decoding private envelope key")
	}

	decryptedKey, err := rsa.DecryptOAEP(sha512.New384(), rand.Reader, pri, wrappedKey, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error while decrypting the key")
	}

	log.Info("pkg/wpm/util/encrypt.go:Encrypt() Successfully unwrapped key")
	return decryptedKey, nil
}
