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

func Fatal[T any](x T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return x
}
