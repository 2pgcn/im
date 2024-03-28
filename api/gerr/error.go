package gerr

import (
	"fmt"
	"runtime"
)

// 获取stack返回map[string]string供err添加md,通过errencode里取出上报
func GetStack() map[string]string {
	_ = runtime.StartTrace()
	var callers string
	pcs := make([]uintptr, 2)
	depth := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:depth])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		s := fmt.Sprintf("%s: %d %s \r\n", frame.File, frame.Line, frame.Function)
		callers = s + "/r/n" + callers
		if !more {
			break
		}
	}
	return map[string]string{"stack": callers}
}
