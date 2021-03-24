FROM golang:1.16.2-alpine3.13 as builder

LABEL maintainer="Leandro Loureiro <leandroloureiro@pm.me>"

ARG APP_VERSION=dev

WORKDIR $GOPATH/src/github.com/lealoureiro/mortgage-calculator-api

COPY controller controller
COPY model model
COPY monthlypayments monthlypayments
COPY utils utils
COPY config config
COPY mortgage-calculator-api.go .
COPY go.mod .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-X 'github.com/lealoureiro/mortgage-calculator-api/config.Version=${APP_VERSION}'" -installsuffix cgo -o /go/bin/mortgage-calculator-api .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/bin/mortgage-calculator-api .

EXPOSE 5000

CMD ["./mortgage-calculator-api"]