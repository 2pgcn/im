package gerr

import (
	"fmt"
	"runtime"
)

// GetStack 获取stack返回map[string]string供err添加md,通过errencode里取出上报
func GetStack() map[string]string {
	_, file, line, _ := runtime.Caller(1)
	return map[string]string{"stack": fmt.Sprintf("%s:%d", file, line)}
}
