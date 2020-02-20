FROM golang:1.13 AS build

WORKDIR /app
COPY . /app

RUN go build -o ./bin/writelog cmd/writelog/main.go && \
    go build -o ./bin/logmonitor cmd/logmonitor/main.go


FROM golang:1.13

WORKDIR /app
COPY --from=build /app/bin/* /app/
RUN touch /tmp/access.log

# Inifinitely write lines to /tmp/access.log
CMD ["bash", "-c","while true; do /app/writelog --lines=3000 --duration=180 --path=/tmp/access.log; done"]