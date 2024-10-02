FROM --platform=linux/amd64 golang:alpine as packager

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

FROM --platform=linux/amd64 golang:alpine as builder

WORKDIR /app

COPY --from=packager /go/pkg/mod /go/pkg/mod

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o tln ./main.go

FROM --platform=linux/amd64 alpine as runner

WORKDIR /app

COPY --from=builder /app/tln .

EXPOSE 3000

ENTRYPOINT [ "/app/tln" ]
