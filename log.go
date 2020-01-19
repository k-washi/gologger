package gologger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type logFormat struct {
	TimestampFormat string
}

//Format ログの形式を設定
func (f *logFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteByte('[')
	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteString("]:")
	b.WriteString(entry.Time.Format(f.TimestampFormat))

	b.WriteString(" [")
	b.WriteString(formatFilePath(entry.Caller.File))
	b.WriteString(":")
	fmt.Fprint(b, entry.Caller.Line)
	b.WriteString("] ")

	if entry.Message != "" {
		b.WriteString(" - ")
		b.WriteString(entry.Message)
	}

	if len(entry.Data) > 0 {
		b.WriteString(" || ")
	}
	for key, value := range entry.Data {
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteByte('{')
		fmt.Fprint(b, value)
		b.WriteString("}, ")
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

//init パッケージ読み込み時に実行される。
func init() {
	logrus.SetReportCaller(true) //Caller(実行ファイル(ex. main.go)を扱うため)
	formatter := logFormat{}
	formatter.TimestampFormat = "2006-01-02 15:04:05" //時刻設定

	logrus.SetFormatter(&formatter)

	//ログ出力ファイルの設定
	f, err := openFile("log.txt")
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, f))

	//ログレベルの設定
	logrus.SetLevel(logrus.InfoLevel)

}

//SetLevelDebug Debugレベルに設定
func SetLevelDebug() {
	logrus.SetLevel(logrus.DebugLevel)
}

//SetLevelInfo Set Infoレベルに設定
func SetLevelInfo() {
	logrus.SetLevel(logrus.InfoLevel)
}
