FROM golang:alpine AS build

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-s -w -extldflags '-static'" -o main ./cmd/app/main.go

FROM alpine

WORKDIR /build/app/

COPY --from=build /build/ .

RUN apk add --no-cache tzdata

RUN chmod +x /bin/goose

CMD /build/app/main
