package archive

import (
	"archive/tar"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func newStreamEncryptWriter(w io.Writer, key []byte) (io.Writer, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	nonce := make([]byte, aes.BlockSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, nil, err
	}
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, nil, err
	}

	stream := cipher.NewCTR(block, nonce)
	writer := cipher.StreamWriter{
		S: stream,
		W: w,
	}
	return writer, nonce, nil
}

// Encrypt data with AES-GCM
func encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Add file to tar archive
func addFileToArchive(tw *tar.Writer, filename string, fileContent []byte) error {
	header := &tar.Header{
		Name: filename,
		Mode: 0600,
		Size: int64(len(fileContent)),
	}

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	_, err := tw.Write(fileContent)

	return err
}

func CreateEncryptedTarArchive(sourceDir, outputFile, unlockDate, password, key string) error {
	archiveFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer archiveFile.Close()

	salt := make([]byte, 16)
	rand.Read(salt)

	// Create encryptio writer
	encWriter, nonce, err := newStreamEncryptWriter(archiveFile, DeriveKey(password, salt))
	if err != nil {
		return err
	}

	// Write to the start of archiveFile for decryption later
	if _, err = archiveFile.Write(salt); err != nil {
		return err
	}

	// Write nonce after salt
	if _, err = archiveFile.Write(nonce); err != nil {
		return err
	}

	tw := tar.NewWriter(encWriter)
	defer tw.Close()

	// Add encrypted unlock time file
	encryptedFileContent, err := encrypt([]byte(unlockDate), DeriveKey(key, salt))
	if err != nil {
		return err
	}

	if err = addFileToArchive(tw, ".ut", encryptedFileContent); err != nil {
		return err
	}

	err = filepath.Walk(sourceDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a tar header
		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		// Update the header name to maintain the directory structure
		relPath, err := filepath.Rel(filepath.Dir(sourceDir), path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// Write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// If not a regular file, skip writing content
		if !info.Mode().IsRegular() {
			return nil
		}

		// Open file for reading
		fileToCopy, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileToCopy.Close()

		// Copy file content into the tw writer
		_, err = io.Copy(tw, fileToCopy)

		return err
	})

	return err
}
