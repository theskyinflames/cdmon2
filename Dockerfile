FROM golang:1.11.5-alpine3.9

RUN apk add make

RUN mkdir -p $GOPATH/src/github.com/theskyinflames/cdmon2
WORKDIR $GOPATH/src/github.com/theskyinflames/cdmon2
COPY . ./
RUN ls -la

RUN make build 
CMD ["make","run"]