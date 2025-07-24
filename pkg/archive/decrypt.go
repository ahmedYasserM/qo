package archive

import (
	"archive/tar"
	"crypto/aes"
	"crypto/cipher"
	"io"
	"os"
	"time"

	"github.com/ahmedYasserM/qo/pkg/logger"
)

func newStreamDecryptReader(r io.Reader, key []byte, nonce []byte) (io.Reader, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, nonce)
	reader := cipher.StreamReader{
		S: stream,
		R: r,
	}

	return reader, nil
}

// Decrypt data with AES-GCM
func decrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return nil, io.ErrUnexpectedEOF
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return aesgcm.Open(nil, nonce, ciphertext, nil)
}

// Check if we reach the unlock time or not
func checkUnlockTime(ut string) (bool, error) {
	layout := "2006-01-02 15:04"

	parsedUt, err := time.ParseInLocation(layout, ut, time.Local)
	if err != nil {
		return false, err
	}

	now := time.Now()

	return now.After(parsedUt) || now.Equal(parsedUt), nil
}

func DecryptTarArchive(encryptedFile, password string) error {
	file, err := os.Open(encryptedFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read salt (16 bytes)
	salt := make([]byte, 16)
	if _, err := io.ReadFull(file, salt); err != nil {
		return err
	}

	// Derive key
	key := DeriveKey(password, salt)

	// Read nonce
	nonce := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(file, nonce); err != nil {
		return err
	}

	// Create a stream decrypt reader from the encrypted file
	decryptReader, err := newStreamDecryptReader(file, key, nonce)
	if err != nil {
		return err
	}

	// Open the decrypted tar archive in memory
	tr := tar.NewReader(decryptReader)
	var ut []byte

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return err
		}

		if header.Name == ".ut" {
			encryptedUt, err := io.ReadAll(tr)
			if err != nil {
				return err
			}

			// Decrypt the unlock time
			if ut, err = decrypt(encryptedUt, DeriveKey("bar", salt)); err != nil {
				return err
			}

			break
		}
	}

	// if the current time >= the ulock time then canProceed wth the decryption
	canProceed, err := checkUnlockTime(string(ut))
	if err != nil {
		return err
	}

	if !canProceed {
		logger.Warn("Can not decrypt the archive before the unlock time.")
		return nil
	}

	logger.Info("Unlock time reached. Extracting archive...")

	// Start decryption again from the beginning, extract all files except .ut

	_, err = file.Seek(16+aes.BlockSize, io.SeekStart)
	if err != nil {
		return err
	}

	decryptReader, err = newStreamDecryptReader(file, key, nonce)
	if err != nil {
		return err
	}

	tr = tar.NewReader(decryptReader)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // End of archive
		}

		if err != nil {
			return err
		}

		if header.Name == ".ut" {
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(header.Name, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			toFile, err := os.OpenFile(header.Name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			defer toFile.Close()
			if err != nil {
				return err
			}

			if _, err = io.Copy(toFile, tr); err != nil {
				toFile.Close()
				return err
			}

		}
	}

	return nil
}
