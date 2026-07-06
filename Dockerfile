FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1
COPY . .
RUN sqlc generate
RUN CGO_ENABLED=0 go build -o /prisoner .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /prisoner /app/prisoner

EXPOSE 5001

# The image never migrates on its own. Default run serves HTTP; migrations are a
# deliberate one-off invocation of the same image: `<image> migrate up`.
ENTRYPOINT ["/app/prisoner"]
CMD ["server"]
