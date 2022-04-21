FROM golang:1.18.1-bullseye AS builder

WORKDIR /main
RUN go env -w GOPROXY=direct
RUN go env -w CGO_ENABLED=0
ADD main.go go.mod go.sum .
RUN go mod download
RUN go build -ldflags="-w -s" -o main main.go

FROM scratch
COPY --from=builder /main/main .
CMD ["./main"]