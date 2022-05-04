package logging

import (
	"ess/utils/file"
	"ess/utils/setting"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// returns the folder which contains the logs file, according to the project root path
func getLogFilePath() string {
	return fmt.Sprintf("%s%s", setting.AppSetting.RuntimeRootPath, setting.AppSetting.LogSavePath)
}

// returns a log file name
func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s", setting.AppSetting.LogSaveName, time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt)
}

// create the folder according the filePath (if not exists), create and open the logFile according to the fileName,
// and return this file
func openLogFile(fileName string, filePath string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}
	src := dir + string(filepath.Separator) + filePath
	// check the permission is important!
	if perm := file.CheckPermission(src); perm {
		return nil, fmt.Errorf("file checkPermission denied src: %s", src)
	}
	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}
	// create the file
	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file: %v", err)
	}
	return f, nil

}
