FROM golang:1.14
RUN go get -d -v -x github.com/go-snart/example github.com/superloach/mapper
RUN printf 'package main\n\nimport _ "github.com/superloach/mapper"\n' > /go/src/github.com/go-snart/example/plugins.go
RUN go build -v -x -o /go/bin/mapper github.com/go-snart/example
CMD ["mapper"]
