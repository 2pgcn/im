package pprof

import (
	"errors"
	"github.com/grafana/pyroscope-go"
	"os"
	"runtime"
	"sync"
	"time"
)

var profiler *pyroscope.Profiler
var one sync.Once
var isInit bool
var InitErr = errors.New("pyroscope be init, not allow init many")

func GetProfiler() *pyroscope.Profiler {
	return profiler
}
func InitPyroscope(appName string, version string, endpoint string, log pyroscope.Logger) (err error) {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)
	if isInit {
		return
	}
	one.Do(func() {
		isInit = true
	})
	profiler, err = pyroscope.Start(pyroscope.Config{
		ApplicationName: appName,
		ServerAddress:   endpoint,
		Logger:          log,
		Tags:            map[string]string{"hostname": os.Getenv("HOSTNAME"), "environment": "test", "version": version},
		UploadRate:      time.Second * 10,
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
	return err

}
