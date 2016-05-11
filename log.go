package kash

import (
	"github.com/op/go-logging"
	"os"
)

var (
	log = logging.MustGetLogger("kash")

	format = logging.MustStringFormatter(
		"%{level} %{time:2006-01-02T15:04:05.999999Z07:00} %{message}",
	)
)

func init() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	logging.SetBackend(backend1Leveled, backend2Formatter)
}