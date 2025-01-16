FROM golang:1.23.4 AS build
WORKDIR /src
COPY . .
RUN go get -v
RUN CGO_ENABLED=0 go build -o /bin/service .

FROM alpine:latest
WORKDIR /app
COPY --from=build /bin/service /bin/service
CMD ["/bin/service"]