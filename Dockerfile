FROM golang:1.25.0 AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o /build/untear .

FROM scratch

COPY --from=builder /build/untear /untear

WORKDIR /worlds

ENTRYPOINT [ "/untear" ]