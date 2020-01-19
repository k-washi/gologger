# Golang logrusによるロギング

# 概要

golangのロギングパッケージである、logrusを用いて、以下のようにログを出力させる設定方法について記載する。

```bash
[INFO]:2020-01-20 03:06:19 [main.go:19]  - info msg1
```

出力内容は、レベル、時刻、ファイル名、ライン位置、メッセージである。

また、以下についても記載する。

+ ログファイルへの出力
+ ログレベルの変更

Git: [k-washi/gologger](https://github.com/k-washi/gologger)

# インストール

```bash
go get github.com/sirupsen/logrus
```

# logrusの設定ファイル

まず、設定に使用する関数を作成する。
詳細はプログラム内に記載。

```go:utils.go
package gologger

import (
	"os"
	"strings"
)

//openFile ログを出力するファイルを設定する。
//ファイルが存在する場合、ファイルにログを追記。
//ファイルが存在しない場合、ファイルを作成し、ログを出力。
func openFile(fileName string) (*os.File, error) {
	if exists(fileName) {
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0777)
		return f, err
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	return f, err
}

//formatFilePath ログに記載するファイル名の抽出
func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]

}

//exists　ファイルが存在するか確認する。
func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
```

[logrusのformat](https://github.com/x-cray/logrus-prefixed-formatter/blob/master/formatter.go)を参考にして、ロギングフォーマットを設定。
詳細はプログラム内に記載。

```go:log.go
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

```

# ロギング

以下のように、main.goで設定パッケージを読み込むことで、logrusの設定が読み込まれる。
もし、ログレベルを変更したい場合は、コメントアウトした関数を使用する。

```go:main.go
package main

import (
	"github.com/k-washi/gologger/test/subpkg"

	_ "github.com/k-washi/gologger"
	log "github.com/sirupsen/logrus"
)

func main() {
	//logSetter.SetLevelDebug()
	log.Info("info msg1")
	log.Debug("debug msg1")

	subpkg.SubPrint()

}
```

```go:subpkg.go
package subpkg

import (
	log "github.com/sirupsen/logrus"

)

func SubPrint() {

	log.Info("sub info msg1")
	log.Debug("sub debug msg1")

}
```