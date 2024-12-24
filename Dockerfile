FROM golang:1.23 AS builder

WORKDIR /build

RUN --mount=type=bind,target=. CGO_ENABLED=0 go build -o /bin/app .

FROM gcr.io/distroless/static-debian12

COPY --from=builder /bin/app /bin/app

ENTRYPOINT ["/bin/app"]
