FROM golang:1.23-bookworm AS build
WORKDIR /app
COPY . .
RUN go mod download \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .


FROM gcr.io/distroless/static-debian12:latest AS production
WORKDIR /app
COPY --from=busybox:1.37.0-uclibc /bin/busybox /bin/busybox
COPY --from=busybox:1.37.0-uclibc /bin/cp       /bin/cp
COPY --from=busybox:1.37.0-uclibc /bin/rm       /bin/rm
COPY --from=build /app/sssg .

ENTRYPOINT ["./sssg"]