package rest

import (
    "github.com/sirupsen/logrus"
	"github.com/femaref/helper"
	"github.com/femaref/toJson"
	"net/http"

	"fmt"

	"runtime"
)

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

    w.Header().Set("Content-Type", "application/json")       
    w.WriteHeader(code)
    if writer, ok := err.(toJson.JsonWriter); ok {
        err := toJson.WriteToJson(w, writer)
        if err != nil {
            if helper.Logger != nil {
                helper.Logger.WithFields(logrus.Fields{}).Error(err)
            }
        }        
    } else {
        toJson.WriteToJson(w, struct{ Text string }{err.Error()})
    }

    return true

}

func ErrorCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func ErrorWrapper(w http.ResponseWriter, err error) bool {
	return ShowError(w, err, 500)
}
