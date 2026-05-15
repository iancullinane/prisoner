FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /prisoner .

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /prisoner /app/prisoner

EXPOSE 5001

ENTRYPOINT ["/app/prisoner"]
CMD ["-server"]
