FROM golang:1.14 AS snart
RUN bash -c "go get -d -v -x github.com/{go-snart/{db,route,bot,snart,plugin-{help,admin}},superloach/mapper}"
WORKDIR /go/src/github.com/superloach/mapper
RUN go install -i -v -x .
CMD ["mapper"]
