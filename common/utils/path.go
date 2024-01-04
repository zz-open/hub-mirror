package utils

import "os"

func GetPwdPath() string {
	rootPath, _ := os.Getwd()
	return rootPath
}
