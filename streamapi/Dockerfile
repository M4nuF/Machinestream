FROM golang
COPY . /go/src/github.com/Machinestream/
WORKDIR /go/src/github.com/Machinestream/streamapi
RUN go build -o streamapi

FROM alpine
COPY --from=0 /go/src/github.com/Machinestream/streamapi/streamapi /usr/bin/
EXPOSE 8080
CMD streamapi