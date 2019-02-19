FROM golang:1-alpine as builder
RUN apk --no-cache add \
	build-base \
	git
WORKDIR /go/src/github.com/TheoBrigitte/confluence

ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go install -ldflags '-s -w' ./cmd/confluence
#RUN go get -u -v github.com/anacrolix/confluence

FROM alpine as runtime
RUN apk --no-cache add \
        ca-certificates \
	libgcc \
	libstdc++
COPY --from=builder /go/bin/confluence /bin/confluence

EXPOSE 50007
ENTRYPOINT ["/bin/confluence"]
