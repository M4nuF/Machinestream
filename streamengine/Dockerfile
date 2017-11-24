FROM golang
COPY . /go/src/github.com/Machinestream/
WORKDIR /go/src/github.com/Machinestream/streamengine
RUN go build -o streamengine

FROM alpine
COPY --from=0 /go/src/github.com/Machinestream/streamengine/streamengine /usr/bin/
CMD streamengine
