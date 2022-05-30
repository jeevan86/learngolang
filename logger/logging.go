package logging

import (
	"github.com/jeevan86/lf4go/factory"
	"gopackettest/config"
	"reflect"
	"strings"
)

type EMPTY int8

var LoggerFactory = factory.NewLoggerFactory(
	func(caller string) string {
		projectName := ""
		rootName := config.Config.Logging.RootName
		if len(rootName) > 0 {
			projectName = rootName
		} else {
			myPackage := reflect.TypeOf(EMPTY(0)).PkgPath()
			projectName = myPackage[:strings.LastIndex(myPackage, "/")]
		}
		callerPackage := caller[strings.Index(caller, projectName):]
		callerPackage = callerPackage[strings.Index(callerPackage, "/"):strings.LastIndex(callerPackage, "/")]
		return callerPackage
	},
)
