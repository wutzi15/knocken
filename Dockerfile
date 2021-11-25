FROM golang:latest

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build knocken.go

CMD ["/go/src/app/knocken"]