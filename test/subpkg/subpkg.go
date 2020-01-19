package subpkg

import (
	log "github.com/sirupsen/logrus"
)

func SubPrint() {

	log.Info("sub info msg1")
	log.Debug("sub debug msg1")

}
