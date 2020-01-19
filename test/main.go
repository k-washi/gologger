package main

import (
	"github.com/k-washi/gologger/test/subpkg"

	_ "github.com/k-washi/gologger"
	log "github.com/sirupsen/logrus"
)

/*
[INFO]:2020-01-20 03:05:53 [main.go:12]  - info msg1
[INFO]:2020-01-20 03:05:53 [subpkg.go:11]  - sub info msg1

[INFO]:2020-01-20 03:06:19 [main.go:19]  - info msg1
[DEBUG]:2020-01-20 03:06:19 [main.go:20]  - debug msg1
[INFO]:2020-01-20 03:06:19 [subpkg.go:11]  - sub info msg1
[DEBUG]:2020-01-20 03:06:19 [subpkg.go:12]  - sub debug msg1
*/

func main() {
	//logSetter.SetLevelDebug()
	log.Info("info msg1")
	log.Debug("debug msg1")

	subpkg.SubPrint()

}
