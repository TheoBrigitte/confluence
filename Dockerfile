FROM golang:1-alpine as builder
RUN apk --no-cache add \
	build-base \
	git
RUN go get -u -v github.com/anacrolix/confluence

FROM alpine as runtime
RUN apk --no-cache add \
	libgcc \
	libstdc++
COPY --from=builder /go/bin/confluence /bin/confluence
EXPOSE 7803 50007
ENTRYPOINT ["/bin/confluence", "-addr=0.0.0.0:7803"]
