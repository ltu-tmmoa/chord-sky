FROM golang:alpine

COPY . /go/src/github.com/ltu-tmmoa/chord-sky
RUN go install github.com/ltu-tmmoa/chord-sky

ENTRYPOINT ["/go/bin/chord-sky"]

EXPOSE 8080
