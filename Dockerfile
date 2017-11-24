FROM golang
COPY . /go/src/streamengine/
WORKDIR /go/src/streamengine/
RUN go build -o streamengine

FROM alpine
WORKDIR /usr/bin/MachineStream
COPY --from=0 /go/src/streamengine/streamengine /usr/bin/
CMD streamengine
