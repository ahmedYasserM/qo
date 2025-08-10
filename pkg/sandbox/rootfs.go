package sandbox

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	_ "embed"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/ahmedYasserM/qo/pkg/logger"
)

//go:embed rootfs.tar.gz
var embeddedRootfs []byte

const target = "/tmp"

var Rootfs string = filepath.Join(target, "rootfs")

// PathExists checks if a file or directory exists.
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ExtractRootfs extracts the tar-archived rootfs folder in /tmp
func ExtractRootfs() error {
	if pathExists(Rootfs) {
		if err := os.RemoveAll(Rootfs); err != nil {
			return err
		}
	}

	gzReader, err := gzip.NewReader(io.NopCloser(bytes.NewReader(embeddedRootfs)))
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // done
		}
		if err != nil {
			return err
		}

		destPath := filepath.Join(target, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return err
			}

			outFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()

		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return err
			}

			if err := os.Symlink(header.Linkname, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func StartSandBox() error {
	if len(os.Args) == 1 && os.Args[0] == "init" {
		if err := syscall.Chroot(Rootfs); err != nil {
			return err
		}

		if err := os.Chdir("/"); err != nil {
			return err
		}

		if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
			return err
		}

		logger.Info("You are now inside the isolated enviornemnt.")

		cmd := exec.Command("/bin/bash")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err := syscall.Unmount("/proc", 0); err != nil {
			return err
		}

		return err
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.Args = []string{"init"}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	err := cmd.Run()

	return err
}
