FROM golang:1.20.1 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn  go build -gcflags="all=-N -l"  -o /src/logic /src/cmd/logic...
RUN GOPROXY=https://goproxy.cn  go build -gcflags="all=-N -l"  -o /src/comet /src/cmd/comet...
RUN GOPROXY=https://goproxy.cn  go build -gcflags="all=-N -l"  -o /src/bench /src/benchmark...
# 放一块编译为了自己方便,最好分开编译成不同的镜像
FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /src/comet /app
COPY --from=builder /src/logic /app
COPY --from=builder /src/bench /app

#通过configmap构建进去
#COPY --from=builder /src/conf/logic.yaml /app/conf/logic.yaml
#COPY --from=builder /src/conf/comet.yaml /app/conf/comet.yaml

#EXPOSE 9001

#CMD ["/app/logic","-conf=/app/conf/logic.yaml"]
