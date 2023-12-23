FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN ./scripts/build.sh

FROM gcr.io/distroless/static-debian11

COPY --from=builder /app/binary /binary

ENTRYPOINT ["/binary"]