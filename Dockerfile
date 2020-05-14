FROM golang:1.14 AS snart
RUN go get -v -x github.com/superloach/mapper/cmd/mapper
CMD ["mapper"]
