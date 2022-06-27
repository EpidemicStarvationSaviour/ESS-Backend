FROM golang:1.18-bullseye
EXPOSE 8080
WORKDIR /workspace
RUN sed -i "s@http://\(deb\|security\).debian.org@https://mirrors.xxx.com@g" /etc/apt/sources.list
RUN apt update
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
WORKDIR /app
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
RUN go mod tidy
# src code
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build -o ess-backend .
    RUN chmod +x ess-backend

ENV TZ=Asia/Shanghai
ENTRYPOINT ["./ess-backend"]