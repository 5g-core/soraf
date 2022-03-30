package util

import (
	"bufio"
	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/logger"
	"os"
	"path/filepath"
	"strings"
)

var (
	ServiceLogPath           = ServicePath("service/sslkey.log")
	ServicePemPath           = ServicePath("service/support/TLS/.pem")
	ServiceKeyPath           = ServicePath("service/support/TLS/.key")
	DefaultServiceConfigPath = ServicePath("service/config/servicecfg.yaml")
)

/*
 *
 * This package is used to locate the root directory of   project
 * Compatible with Windows and Linux
 *

 *
 * Return value:
 * A string value of the relative path between the working directory and the root directory of the   project
 *
 * Usage:
 *  ServicePath("your file location  ")
 *
 * Example:
 *  ServicePath("service/abcdef/abcdef.pem")
 */
func ServicePath(path string) string {
	rootCode := strings.Split(path, "/")[0]
	cleanPath := filepath.Clean(path)
	targetFilePath := cleanPath[len(rootCode)+1:]

	var pwd string
	if pwdTmp, err := os.Getwd(); err != nil {
		logger.UtilLog.Errorln(err)
	} else {
		pwd = pwdTmp
	}
	currentPath := filepath.Clean(pwd)

	// Module mode
	target := ""
	if returnPath, ok := FindModuleRoot(currentPath, rootCode); ok {
		target = returnPath + filepath.Clean("/"+targetFilePath)
	}

	// Non-module mode
	if target == "" {
		binPathDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.UtilLog.Errorln(err)
		}

		rootPath := ""
		if strings.Contains(currentPath, rootCode) {
			if returnPath, ok := FindRoot(currentPath, rootCode, targetFilePath); ok {
				rootPath = returnPath
			} else if returnPath, ok := FindRoot(currentPath, rootCode, "lib"); ok {
				rootPath = returnPath
			}
		}
		if strings.Contains(binPathDir, rootCode) {
			if returnPath, ok := FindRoot(binPathDir, rootCode, targetFilePath); ok {
				rootPath = returnPath
			} else if returnPath, ok := FindRoot(binPathDir, rootCode, "lib"); ok {
				rootPath = returnPath
			}
		}

		if rootPath != "" {
			target = rootPath + cleanPath
		} else {
			binPathDirParent := GetParentDirectory(binPathDir)
			binPathDirParentWithTargetFilePath := binPathDirParent + filepath.Clean("/"+targetFilePath)
			target = binPathDirParentWithTargetFilePath
		}
	}

	location, err := filepath.Rel(currentPath, target)
	if err != nil {
		logger.UtilLog.Errorln(err)
	}

	return location
}

func Exists(fpath string) bool {
	_, err := os.Stat(fpath)
	return !os.IsNotExist(err)
}

func FindRoot(path string, rootCode string, objName string) (string, bool) {
	rootPath := path
	loc := strings.LastIndex(rootPath, rootCode)
	for loc != -1 {
		rootPath = rootPath[:loc+len(rootCode)]
		if Exists(rootPath + filepath.Clean("/"+objName)) {
			return rootPath[:loc], true
		}
		rootPath = rootPath[:loc]
		loc = strings.LastIndex(rootPath, rootCode)
	}
	return "", false
}

func FindModuleRoot(path string, rootCode string) (string, bool) {
	moduleFilePath := path + filepath.Clean("/go.mod")
	if Exists(moduleFilePath) {
		var file *os.File
		if fileTmp, err := os.Open(moduleFilePath); err != nil {
			logger.UtilLog.Fatalf("Cannot open %s: %+v", moduleFilePath, err)
		} else {
			file = fileTmp
		}
		defer func() {
			if err := file.Close(); err != nil {
				logger.UtilLog.Warnf("File %s cannot close: %v", moduleFilePath, err)
			}
		}()

		reader := bufio.NewReader(file)
		moduleDeclearation, _, err := reader.ReadLine()
		if err != nil {
			logger.UtilLog.Warnf("Read Line failed: %+v", err)
		}
		if string(moduleDeclearation) == "module "+rootCode {
			return path, true
		}
	}

	abs, err := filepath.Abs(path + string(filepath.Separator) + "..")
	if err != nil || abs == filepath.Clean("/") {
		return "", false
	}

	return FindModuleRoot(abs, rootCode)
}

func GetParentDirectory(dirctory string) string {
	return filepath.Dir(dirctory)
}
