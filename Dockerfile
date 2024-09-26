FROM golang:1.23

WORKDIR /julien

COPY contract /julien/contract
COPY driver /julien/driver
COPY form /julien/form
COPY fs /julien/fs
COPY web /julien/web
COPY pager /julien/pager
COPY utils /julien/utils
COPY julien /julien/julien
COPY template /julien/template
COPY main.go /julien/main.go

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum /julien/
RUN go mod download && go mod verify
RUN go build -v -o /usr/local/bin/julien

WORKDIR /www
RUN julien --help

