package panics

import (
	"bytes"
	"fmt"
	"github.com/jeevan86/learngolang/pkg/util/str"
	"runtime"
)

func SafeRet(f func() interface{}) (ret interface{}, sta string, err error) {
	defer func() {
		e := recover()
		if e != nil {
			ret = nil
			sta, err = catch(e)
		}
	}()
	ret = f()
	err = nil
	sta = str.EMPTY()
	return
}

func SafeRun(f func()) (sta string, err error) {
	defer func() {
		e := recover()
		if e != nil {
			sta, err = catch(e)
		}
	}()
	f()
	err = nil
	sta = str.EMPTY()
	return
}

func catch(err interface{}) (string, error) {
	if err != nil {
		switch err.(type) {
		case error:
			return stack(err), err.(error)
		default:
			return str.EMPTY(), nil
		}
	} else {
		return str.EMPTY(), nil
	}
}

func stack(err interface{}) string {
	if err != nil {
		buf := new(bytes.Buffer)
		_, _ = fmt.Fprintf(buf, "err => %v\n", err)
		for i := 1; ; i++ {
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		}
		return buf.String()
	}
	return str.EMPTY()
}
