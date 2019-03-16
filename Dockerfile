FROM golang:1-alpine3.9 AS build
ARG tsuru_workdir="${GOPATH}/src/github.com/tsuru/tsuru"
COPY ./ ${tsuru_workdir}
WORKDIR ${tsuru_workdir}
RUN apk add --no-cache ca-certificates gcc make musl-dev && \
    make tsurud TSR_BIN=/bin/tsurud TSR_BUILD_EXTRAFLAGS='-i -v --ldflags "-linkmode external -extldflags \"-static\""'

FROM alpine:3.9
COPY --from=build /bin/tsurud /bin/
COPY ./etc/tsuru-custom.conf /etc/tsuru/tsuru.conf
RUN apk add --no-cache ca-certificates
EXPOSE 8080
ENTRYPOINT ["/bin/tsurud"]
CMD ["api"]
