# FROM golang:1.8-onbuild


FROM golang:1.8

WORKDIR /go/src/github.com/prantoran/go-elastic-textsearch
COPY . .

RUN pwd
RUN ls

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]