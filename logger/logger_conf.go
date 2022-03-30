package logger

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var ServiceLogDir string = ServicePath("service/log") + "/"

var ServiceLogfle string = ServiceLogDir + "service.log"

func ServicePath(path string) string {
	rootCode := strings.Split(path, "/")[0]
	cleanPath := filepath.Clean(path)
	targetFilePath := cleanPath[len(rootCode)+1:]

	var pwd string
	if pwdTmp, err := os.Getwd(); err != nil {
		UtilLog.Errorln(err)
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
			UtilLog.Errorln(err)
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
		UtilLog.Errorln(err)
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
			UtilLog.Fatalf("Cannot open %s: %+v", moduleFilePath, err)
		} else {
			file = fileTmp
		}
		defer func() {
			if err := file.Close(); err != nil {
				UtilLog.Warnf("File %s cannot close: %v", moduleFilePath, err)
			}
		}()

		reader := bufio.NewReader(file)
		moduleDeclearation, _, err := reader.ReadLine()
		if err != nil {
			UtilLog.Warnf("Read Line failed: %+v", err)
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

func init() {
	if err := os.MkdirAll(ServiceLogDir, 0775); err != nil {
		log.Printf("Mkdir %s failed: %+v", ServiceLogDir, err)
	}

	// Create log file or if it already exist, check if user can access it
	f, fileOpenErr := os.OpenFile(ServiceLogfle, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if fileOpenErr != nil {
		// user cannot access it.
		log.Printf("Cannot Open %s\n", ServiceLogfle)
	} else {
		// user can access it
		if err := f.Close(); err != nil {
			log.Printf("File %s cannot been closed\n", ServiceLogfle)
		}
	}

	sudoUID, errUID := strconv.Atoi(os.Getenv("SUDO_UID"))
	sudoGID, errGID := strconv.Atoi(os.Getenv("SUDO_GID"))

	if errUID == nil && errGID == nil {
		// if using sudo to run the program, errUID will be nil and sudoUID will get the uid who run sudo
		// else errUID will not be nil and sudoUID will be nil
		// If user using sudo to run the program and create log file, log will own by root,
		// here we change own to user so user can view and reuse the file
		if err := os.Chown(ServiceLogDir, sudoUID, sudoGID); err != nil {
			log.Printf("Dir %s chown to %d:%d error: %v\n", ServiceLogDir, sudoUID, sudoGID, err)
		}

		if fileOpenErr == nil {
			if err := os.Chown(ServiceLogfle, sudoUID, sudoGID); err != nil {
				log.Printf("File %s chown to %d:%d error: %v\n", ServiceLogfle, sudoUID, sudoGID, err)
			}
		}
	}
}
