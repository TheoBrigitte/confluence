FROM golang:alpine as builder
RUN apk --no-cache add \
	build-base \
	git

WORKDIR /go/src/github.com/TheoBrigitte/confluence
COPY vendor vendor
COPY go.mod .
COPY go.sum .
COPY cmd cmd
COPY pkg pkg
COPY Makefile .

RUN make build


FROM alpine:latest as runtime
RUN apk --no-cache add \
        ca-certificates \
	libgcc \
	libstdc++
COPY --from=builder /go/bin/confluence /bin/confluence
EXPOSE 7803 50007
ENTRYPOINT ["/bin/confluence"]
CMD ["-addr=0.0.0.0:7803"]
