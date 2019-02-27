FROM node:alpine as builder
RUN apk --no-cache add \
	build-base \
	git

WORKDIR /usr/src/app
COPY node_modules node_modules
COPY package.json .
COPY package-lock.json .
COPY public public
COPY src src
COPY Makefile .

RUN make build-client


FROM scratch as runtime
VOLUME /usr/app
COPY --from=builder /usr/src/app/build /usr/app
CMD ["noop"]
