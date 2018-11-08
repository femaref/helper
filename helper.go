package helper

import (
	raven "github.com/femaref/raven-go"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger
var Raven *raven.Client
