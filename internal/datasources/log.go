package datasources

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogFileWriter struct {
	LogPath      string
	FileName     string
	PrintConsole bool
}

func (c *LogFileWriter) Write(body []byte) (n int, err error) {
	if c.PrintConsole {
		log.Printf("%+s", body)
	}

	var logPath string
	var logFileName string

	if len(c.FileName) == 0 {
		logPath = "./log"
	} else {
		logPath = c.LogPath
	}
	err = os.MkdirAll(c.LogPath, 0755)
	if err != nil {
		return
	}

	if len(c.FileName) == 0 {
		logFileName = fmt.Sprintf("access-%s.log", time.Now().Format("2006-02-01"))
	} else {
		logFileName = c.FileName
	}
	logFile, err := os.OpenFile(fmt.Sprintf("%s/%s", logPath, logFileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return
	}
	defer logFile.Close()

	return logFile.Write(body)
}