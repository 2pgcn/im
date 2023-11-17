package main

import (
	"flag"
	"fmt"
	"github.com/2pgcn/gameim/internal/logic/server"
	"github.com/2pgcn/gameim/pkg/gamelog"
	"github.com/2pgcn/gameim/pkg/trace_conf"
	"github.com/go-kratos/kratos/v2"
	"net/http"
	"os"

	"github.com/2pgcn/gameim/conf"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	_ "go.uber.org/automaxprocs"
	_ "net/http/pprof"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	Name           string
	Version        string
	flagconf       string
	KubeConfigPath string

	id, _ = os.Hostname()
)

func init() {
	var err error
	if err != nil {
		panic(err)
	}
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&Name, "name", "logic/app001", "config path, eg: -name test")
}

func initLog() log.Logger {
	l := gamelog.GetZapLog()
	return gamelog.NewHelper(l)
	//
	//writeSyncer := zapcore.AddSync(os.Stdout)
	//encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	//core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	//z := zap.New(core)
	//logger := kratoszap.NewLogger(z)
	//log.SetLogger(log.NewStdLogger(os.Stdout))
	//return logger
}

func newApp(logger log.Logger, gs *grpc.Server, rs *server.OtherServer) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			rs,
		),
		//kratos.Registrar(reg),
	)
}

func main() {
	flag.Parse()
	lg := initLog()
	log.SetLogger(lg)
	logger := log.With(lg,
		//"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		//"service.id", id,
		"service.name", Name,
		//"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	if port := os.Getenv("ILOGTAIL_PROFILE_PORT"); len(port) > 0 {
		startPprof(fmt.Sprintf(":", port))
	}
	trace_conf.SetTraceConfig(bc.Server.TraceConf)
	app, cleanup, err := wireApp(bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()
	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func startPprof(port string) {
	go func() {
		_ = http.ListenAndServe(port, nil)
	}()
}
