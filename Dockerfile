FROM golang:1.13.4

COPY go.mod go.sum /go/src/github.com/deepinbytes/go_voucher/
WORKDIR /go/src/github.com/deepinbytes/go_voucher/
RUN go mod download
COPY . /go/src/github.com/deepinbytes/go_voucher/

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

RUN ls -la

# Build tbe binary and give permissions
#RUN go build main.go
RUN chmod 777 main

# Set binary as entrypoint
ENTRYPOINT ./main

EXPOSE 3000

CMD ./main


