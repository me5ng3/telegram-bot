FROM golang:1.14-alpine AS build
RUN apk --no-cache add ca-certificates

ENV GO111MODULE=on

WORKDIR /src/

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/binary

EXPOSE 8080

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /bin/binary /bin/binary
COPY --from=build /src/config.json /config.json
ENTRYPOINT ["/bin/binary"]
CMD ["./binary"]