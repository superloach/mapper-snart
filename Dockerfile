FROM golang:1.14
RUN GO111MODULE=on go get -u -v -x github.com/superloach/mapper/cmd/mapper
CMD ["mapper"]
