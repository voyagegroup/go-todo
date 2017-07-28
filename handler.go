// Copyright 2017 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.
//
// This handler implementation is based on gddo.
// https://github.com/golang/gddo/tree/master/gddo-server
package base

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/suzuken/go-todo/httputil"
)

type handler func(w http.ResponseWriter, r *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	runHandler(w, r, h, handleError)
}

type errFn func(w http.ResponseWriter, r *http.Request, status int, err error)

func logError(req *http.Request, err error, rv interface{}) {
	if err != nil {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Error serving %s: %v\n", req.URL, err)
		if rv != nil {
			fmt.Fprintln(&buf, rv)
			buf.Write(debug.Stack())
		}
		log.Print(buf.String())
	}
}

func runHandler(w http.ResponseWriter, r *http.Request,
	fn func(w http.ResponseWriter, r *http.Request) error, errfn errFn) {
	defer func() {
		if rv := recover(); rv != nil {
			err := errors.New("handler panic")
			logError(r, err, rv)
			errfn(w, r, http.StatusInternalServerError, err)
		}
	}()

	r.Body = http.MaxBytesReader(w, r.Body, 2048)
	r.ParseForm()
	var buf httputil.ResponseBuffer
	err := fn(&buf, r)
	if err == nil {
		buf.WriteTo(w)
	} else if e, ok := err.(*httputil.HTTPError); ok {
		if e.Status >= 500 {
			logError(r, err, nil)
		}
		errfn(w, r, e.Status, e.Err)
	} else {
		logError(r, err, nil)
		errfn(w, r, http.StatusInternalServerError, err)
	}
}

func errorText(err error) string {
	return "Internal Server error."
}

func handleError(w http.ResponseWriter, r *http.Request, status int, err error) {
	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": errorText(err),
		})
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	io.WriteString(w, errorText(err))
}
