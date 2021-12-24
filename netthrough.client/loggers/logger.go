package loggers

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	//设置日期格式
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	Log.SetFormatter(customFormatter)

	fileWriter := &lumberjack.Logger{
		Filename: "Log/langtian.dcs.login/log.log",
		MaxAge:   30,
		MaxSize:  200,
	}
	fileAndStdoutWriter := io.MultiWriter(fileWriter, os.Stdout) //设置输出到文件和控制台
	Log.SetOutput(fileAndStdoutWriter)
	Log.SetLevel(logrus.TraceLevel)
	fmt.Println("初始化日志成功")
}
