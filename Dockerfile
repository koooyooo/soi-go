FROM golang:alpine AS build
WORKDIR /go/src/github.com/koooyooo/soi-go
COPY . .
RUN GOOS=linux go build -o soi-server cmd/srv/soi-server.go


FROM golang:alpine
WORKDIR /
COPY --from=build /go/src/github.com/koooyooo/soi-go/soi-server soi-server
EXPOSE 8080
CMD ["./soi-server"]
