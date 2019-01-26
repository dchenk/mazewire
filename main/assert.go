package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dchenk/mazewire/pkg/log"
)

func assertType(r *http.Request, assertResult bool, neededType string, original interface{}) {
	if !assertResult {
		msg := fmt.Sprintf("main: could not type-assert %T to a "+neededType, original)
		log.Critical(r, msg, errors.New("invalid type"))
		panic(msg)
	}
}
