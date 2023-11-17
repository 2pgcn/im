FROM golang:1.21.4 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn  go build -gcflags="all=-N -l"  -o /src/logic /src/cmd/logic...
RUN GOPROXY=https://goproxy.cn  go build -gcflags="all=-N -l"  -o /src/comet /src/cmd/comet...
# 放一块编译为了自己方便,最好分开编译成不同的镜像
FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /src/comet /app
COPY --from=builder /src/logic /app

COPY --from=builder /src/conf/logic.yaml /app/conf/logic.yaml
COPY --from=builder /src/conf/comet.yaml /app/conf/comet.yaml


EXPOSE 8000
EXPOSE 9000
