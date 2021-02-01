# builder image
FROM golang:1.15 as builder
RUN mkdir /build
WORKDIR /build
ADD *.go /build/
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o imps-mutating-webhook .

# generate clean, final image for end users
FROM alpine:latest
COPY --from=builder /build/imps-mutating-webhook ./bin
ENTRYPOINT [ "imps-mutating-webhook" ]
