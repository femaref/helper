package rest

import (
	"github.com/Sirupsen/logrus"
	"github.com/femaref/toJson"
	"net/http"

	"fmt"

	"runtime"
)

var Logger *logrus.Logger

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
	if err != nil {
	    if Logger != nil {
    		fname, fpath, fline, ferr := callerInfo()

            if ferr != nil {
                Logger.WithFields(logrus.Fields{}).Error(ferr)
            } else {
                Logger.WithFields(logrus.Fields{"func": fname, "path": fpath, "line": fline}).Error(err)
            }
		}

		w.WriteHeader(code)
		
		if writer, ok := err.(toJson.JsonWriter); ok {
			err := toJson.WriteToJson(w, writer)
			if err == nil {
			    return true
			}
			if Logger != nil {
    			Logger.WithFields(logrus.Fields{}).Error(err)
			}
			
		}

    	toJson.WriteToJson(w, struct{ Text string }{err.Error()})
		return true
	}

	return false

}

func ErrorCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func ErrorWrapper(w http.ResponseWriter, err error) bool {
	return ShowError(w, err, 500)
}
