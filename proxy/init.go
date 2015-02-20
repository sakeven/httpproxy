package proxy

import (
	"github.com/op/go-logging"
	"httpproxy/cache"
	"httpproxy/config"
	stdlog "log"
)

var log = logging.MustGetLogger("proxy")
var cnfg config.Config

//setLog() sets log output format.
func setLog() {
	var level logging.Level
	if cnfg.Log == 1 {
		level = logging.DEBUG
	} else {
		level = logging.INFO
	}

	var format logging.Formatter
	if level == logging.DEBUG {
		format = logging.MustStringFormatter("%{shortfile} %{level} %{message}")
	} else {
		format = logging.MustStringFormatter("%{level} %{message}")
	}
	logging.SetFormatter(format)
	logging.SetLevel(level, "proxy")
}

func init() {
	Caches = make(map[cache.Checksum]*cache.Cache)
	err := cnfg.GetConfig("config/config.json")
	if err != nil {
		stdlog.Fatal(err)
	}
	setLog()
	go CheckCaches()
}
