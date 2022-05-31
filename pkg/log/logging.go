package log

import (
	"fmt"
	"github.com/jeevan86/learngolang/pkg/config"
	"github.com/jeevan86/learngolang/pkg/util/str"
	"github.com/jeevan86/lf4go/factory"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

var reflectVar = int8(0)

var logging = config.GetConfig().Logging
var loggerFactory *factory.LoggerFactory
var mutex = sync.Mutex{}

func initLogging() {
	if loggerFactory != nil {
		return
	}
	mutex.Lock()
	defer mutex.Unlock()
	if loggerFactory != nil {
		return
	}
	loggerFactory = factory.NewLoggerFactory(logging.Factory, callerDetector)
}

var callerDetector = func(caller string) string {
	projectName := str.EMPTY()
	rootName := logging.RootName
	if len(rootName) > 0 {
		projectName = rootName
	} else {
		// 会使用：go.mod文件中的module值 + "/" + 包名
		myPackage := reflect.TypeOf(reflectVar).PkgPath()
		projectName = myPackage[:strings.LastIndex(myPackage, "/")]
	}
	// 当go.mod文件中的module值，与当前项目的目录名称不一样时，这里会有问题
	projectNameIdx := strings.Index(caller, projectName)
	if projectNameIdx < 0 {
		fmt.Println("FATAL!: 项目名称与go.mod中的module不一致，需手动配置config.yml:log.root-name为源码项目目录的名称")
		os.Exit(-1)
	}
	callerPackage := caller[projectNameIdx:]
	firstSlash := strings.Index(callerPackage, str.SLASH())
	lastSlash := strings.LastIndex(callerPackage, str.SLASH())
	if firstSlash < lastSlash {
		callerPackage = callerPackage[firstSlash+1 : lastSlash]
	} else {
		callerPackage = callerPackage[firstSlash+1:]
	}
	return callerPackage
}

var NewLogger = func() *factory.Logger {
	_, callFilePath, _, _ := runtime.Caller(1)
	initLogging()
	return loggerFactory.NewLogger(callFilePath, logging.Formatter, config.GetConfig().Logging.Appenders)
}
