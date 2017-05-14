package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
    log "github.com/femaref/helper"
)


func Unmarshal(buffer io.ReadCloser, target interface{}) error {
	return json.NewDecoder(buffer).Decode(target)
}

func EnsureUnmarshal(buffer io.ReadCloser, target interface{}, required []string) error {
	var copy, original io.ReadCloser
	var err error

	buf, _ := ioutil.ReadAll(buffer)
	original = ioutil.NopCloser(bytes.NewBuffer(buf))
	copy = ioutil.NopCloser(bytes.NewBuffer(buf))

	// decode into map so we can check
	var check map[string]json.RawMessage
	err = Unmarshal(original, &check)
	if err != nil {
		return err
	}

	// check for each field
	for _, rfield := range required {
		if _, ok := check[rfield]; !ok {
			return fmt.Errorf("Could not find field %s in input", rfield)
		}
	}

	fname, fpath, fline, ferr := callerInfo()

	if ferr != nil {
		log.Logger.WithFields(logrus.Fields{}).Error(ferr)
	} else {
		log.Logger.WithFields(logrus.Fields{"func": fname, "path": fpath, "line": fline}).Info(string(buf))
	}

	err = Unmarshal(copy, &target)
	return err
}