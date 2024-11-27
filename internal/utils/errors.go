package utils

import (
	"log"
)

func Must[T any](x T, err error) T {
	if err != nil {
		panic(err)
	}
	return x
}

func MustErr(err error) {
	if err != nil {
		panic(err)
	}
}

func MustOk(ok bool, msg string) {
	if !ok {
		panic(msg)
	}
}

func Fatal[T any](x T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return x
}

func FatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Log[T any](x T, err error) T {
	if err != nil {
		log.Println(err)
	}
	return x
}

func LogErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func LogErrMsg(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
	}
}
