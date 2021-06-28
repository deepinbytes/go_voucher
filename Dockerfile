FROM golang:1.13.4 as builder

COPY go.mod go.sum /go/src/github.com/deepinbytes/go_voucher/
WORKDIR /go/src/github.com/deepinbytes/go_voucher/
RUN go mod download
COPY . /go/src/github.com/deepinbytes/go_voucher/

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest AS production
COPY --from=builder /go/src/github.com/deepinbytes/go_voucher/ .
EXPOSE 3000

CMD ["./main"]



