FROM superloach/snart

# download
RUN go get -d -v -x github.com/superloach/mapper

# build
RUN go build -v -x -o /plugins/mapper -buildmode=plugin github.com/superloach/mapper
