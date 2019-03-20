package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	logSavePath = ".logs/"
	logFileExt  = "log"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s/%s/", logSavePath, time.Now().Format("2006/01/02"))
}

func getLogFileFullPath(level Level) string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s.%s", levelFlags[level], logFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir(getLogFilePath())
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}

	return handle
}

func mkDir(filePath string) {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+filePath, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
