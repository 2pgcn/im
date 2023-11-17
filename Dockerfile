FROM debian:stable-slim

WORKDIR /app
COPY ./conf/logic.yaml /app
COPY ./conf/comet.yaml /app
COPY ./bin /app

COPY --from=builder /src/conf/logic.yaml /app/conf/logic.yaml
COPY --from=builder /src/conf/comet.yaml /app/conf/comet.yaml

EXPOSE 8000
EXPOSE 9000
