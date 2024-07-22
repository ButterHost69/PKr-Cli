package encrypt

import (
	// "bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

// import "path/filepath"

func AESGenerakeKey(length int)([]byte, error) {
	// keep length 16, 24, 32 -> 128, 192, 256 respectively
	key := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}

	return key, nil
}

func AESGenerateIV()([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	return iv, nil
}

func AESEncrypt(source_filepath string, destination_filepath string ,key []byte, IV []byte) (error) {
	inputFile, err := os.Open(source_filepath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destination_filepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	stream := cipher.NewCTR(block, IV)

	writer := &cipher.StreamWriter{S: stream, W: outputFile}

	if _, err := io.Copy(writer, inputFile); err != nil {
		return err
	}

	return nil
}

// func AESDecrypt(){}
// Needs Work... Later Me will do it 
func decryptFile(filename string, key []byte, iv []byte) error {
	inputFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(filename + ".dec")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	stream := cipher.NewCTR(block, iv)

	reader := &cipher.StreamReader{S: stream, R: inputFile}

	if _, err := io.Copy(outputFile, reader); err != nil {
		return err
	}

	return nil
}
