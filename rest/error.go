package rest

import (
	"net/http"

	"github.com/femaref/helper"
	"github.com/femaref/toJson"
	"github.com/sirupsen/logrus"

	"fmt"

	"runtime"
)

type RestError struct {
	Code int
	Err  error
}

func (this RestError) Error() string {
	return fmt.Sprintf("status %d: %v", this.Code, this.Err)
}

func callerInfo() (string, string, string, error) {

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return "", "", "", fmt.Errorf("not enough functions on stack")
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return "", "", "", fmt.Errorf("Could not find function")
	}
	path, line := fun.FileLine(fpcs[0] - 1)
	// return name and path:line_number
	return fun.Name(), path, fmt.Sprintf("%d", line), nil
}

func ShowError(w http.ResponseWriter, err error, code int) bool {
	if err == nil {
		return false
	}

	if helper.Logger != nil {
		fname, fpath, fline, ferr := callerInfo()

		if ferr != nil {
			helper.Logger.WithFields(logrus.Fields{}).Error(ferr)
		} else {
			helper.Logger.WithFields(logrus.Fields{"func": fname, "path": fpath, "line": fline}).Error(err)
		}
	}

	if helper.Raven != nil {
		_, err := helper.Raven.CaptureErrorAndWait(err, nil)
		if helper.Logger != nil {
			helper.Logger.WithFields(logrus.Fields{}).Error(err)
		}
	}

	if writer, ok := err.(toJson.JsonWriter); ok {
		err := toJson.WriteToJsonWithCode(w, writer, code)
		if err != nil {
			if helper.Logger != nil {
				helper.Logger.WithFields(logrus.Fields{}).Error(err)
			}
		}
	} else {
		toJson.WriteToJsonWithCode(w, struct{ Text string }{err.Error()}, code)
	}

	return true

}

func ErrorCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func ErrorWrapper(w http.ResponseWriter, err error) bool {
	if rerr, ok := err.(RestError); ok {
		return ShowError(w, rerr.Err, rerr.Code)
	}
	return ShowError(w, err, 500)
}
