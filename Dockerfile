FROM golang:1.23 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /build/app cmd/main.go

FROM gcr.io/distroless/static-debian12

COPY --from=builder /build/app /bin/app

ENTRYPOINT ["/bin/app"]
