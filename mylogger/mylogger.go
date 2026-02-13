package mylogger

import (
	"io"
	"log"
	"os"
)

// 自定义Logger同时输出文件和标准io。Logger本身内部做了互斥锁，因此是并发安全的。
func LoggerInit(logFilePath string) (*log.Logger, error) {

	//设置输出的文件
	logfile_ptr, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	//同时定向日志输出源为日志文件和标准输出
	mutiWriter := io.MultiWriter(os.Stdout, logfile_ptr)
	logger := log.New(mutiWriter, "<New>", log.Lshortfile|log.Ldate|log.Ltime)

	return logger, nil

}
