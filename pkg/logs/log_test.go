package logs_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zercle/gofiber-skelton/pkg/logs"
)

func TestWriteLog(t *testing.T) {
	logpath := "./log"
	logFileName := fmt.Sprintf("access-%s.log", time.Now().Format("2006-01-02"))
	logWriter := logs.LogFileWriter{
		LogPath:  logpath,
		FileName: logFileName,
	}

	t.Run("success", func(t *testing.T) {
		_, err := logWriter.Write([]byte("log test"))
		assert.NoError(t, err)

		// clear unuse resources
		err = os.RemoveAll(logpath)
		assert.NoError(t, err)
	})
}
