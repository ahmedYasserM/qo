package sandbox

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/ahmedYasserM/qo/pkg/logger"
)

//go:embed rootfs.tar.gz
var embeddedRootfs []byte

const target = "/tmp"

var (
	Rootfs      string = filepath.Join(target, "rootfs")
	defaultUser string = "ahmed"
)

// PathExists checks if a file or directory exists.
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ExtractRootfs extracts the tar-archived rootfs folder in /tmp
func ExtractRootfs() error {
	if pathExists(Rootfs) {
		_ = syscall.Unmount(filepath.Join(Rootfs, "proc"), syscall.MNT_FORCE) // force unmount of /proc to handle possible previous exits using external kill signal

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

func dropToUser(username string) error {
	passwdBytes, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return err
	}
	var uid, gid int
	var homeDir string
	for _, line := range strings.Split(string(passwdBytes), "\n") {
		if strings.HasPrefix(line, username+":") {
			parts := strings.Split(line, ":")
			uid, _ = strconv.Atoi(parts[2])
			gid, _ = strconv.Atoi(parts[3])
			homeDir = parts[5]
			break
		}
	}
	if uid == 0 && username != "root" {
		return fmt.Errorf("user %s not found in chroot /etc/passwd", username)
	}
	if err := syscall.Setgid(gid); err != nil {
		return err
	}
	if err := syscall.Setuid(uid); err != nil {
		return err
	}

	// Set environment variables
	os.Setenv("HOME", homeDir)
	os.Setenv("USER", username)
	os.Setenv("LOGNAME", username)

	return nil
}

func StartSandBox() error {

	if len(os.Args) == 1 && os.Args[0] == "init" {
		if err := syscall.Chroot(Rootfs); err != nil {
			return err
		}

		if err := os.Chdir("/tmp"); err != nil {
			return err
		}

		if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
			return err
		}

		if err := dropToUser(defaultUser); err != nil {
			return err
		}

		logger.Info("You are now inside the isolated enviornemnt.")

		cmd := exec.Command("/bin/bash")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()

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

	if err := cmd.Run(); err != nil {
		return err
	}

	err := syscall.Unmount(Rootfs+"/proc", 0)

	return err
}
