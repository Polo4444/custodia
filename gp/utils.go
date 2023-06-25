package gp

import (
	"log"
	"time"
)

const (
	B   int64 = 1
	KiB int64 = 1024 * B
	MiB int64 = 1024 * KiB
	GiB int64 = 1024 * MiB
	TiB int64 = 1024 * GiB
	PiB int64 = 1024 * TiB
	EiB int64 = 1024 * PiB
	//ZiB int64 = 1024 * EiB
	//YiB int64 = 1024 * ZiB
)

func LogFatalIfErr(err error, mes ...string) {

	if err != nil {
		log.Fatal(err, mes)
	}
}

func CurrentTime() time.Time {
	return time.Now().UTC()
}

func ErrorToString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
