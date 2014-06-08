package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const BACKUP_DST_BASE_DIR = "/tmp/test/dst"

type TargetFile struct {
	path string
}

func main() {
	for _, target := range readTargets() {
		fmt.Println(target)
		copy(target)
	}
}

func readTargets() []TargetFile {
	result := make([]TargetFile, 0)
	targetFile := TargetFile{
		path: "/tmp/test/src/a",
	}
	return append(result, targetFile)
}

func copy(target TargetFile) {
	if !exists(target.path) {
		return
	}
	dstDir, err := createDstDir(target.path)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println(dstDir)
	dstFilePath := dstDir + "/" + filepath.Base(target.path)
	copyFile(target.path, dstFilePath)
}

func exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func createDstDir(path string) (dstDir string, err error) {
	dstDir, err := toDstDir(path)
	if err != nil {
		return "", err
	}
	err = os.MkdirAll(dstDir, 0755)
	return dstDir, err
}

func toDstDir(targetFilePath string) (dstDir string, err error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return BACKUP_DST_BASE_DIR + "/" + hostname + filepath.Dir(targetFilePath), nil
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return
}
