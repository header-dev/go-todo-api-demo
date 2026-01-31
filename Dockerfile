FROM golang:1.25.6-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

ENV GOARCH=arm64

RUN go build \
    -o /go/bin/app

FROM gcr.io/distroless/base-debian11

COPY --from=build /go/bin/app /app
COPY local.env /local.env

EXPOSE 8081

USER nonroot:nonroot

CMD ["/app"]