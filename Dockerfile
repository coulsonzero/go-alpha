FROM golang:1.26

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

RUN apt-get update && apt-get install -y python3 python3-pip && \
    pip3 install requests --break-system-packages

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go-alpha

EXPOSE 8000
ENTRYPOINT [ "/go-alpha" ]
