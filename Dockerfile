FROM golang:alpine3.18 as build
WORKDIR /app
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o build/chat cmd/chat/main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/build/chat .
CMD [ "./chat" ]
