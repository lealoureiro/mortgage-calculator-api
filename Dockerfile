FROM golang:1.13.5-alpine3.10 as builder

RUN apk update && apk upgrade && apk add --no-cache git dep

LABEL maintainer="Leandro Loureiro <leandroloureiro@pm.me>"

WORKDIR $GOPATH/src/github.com/lealoureiro/mortgage-calculator-api

COPY controller controller
COPY model model
COPY monthlypayments monthlypayments
COPY utils utils
COPY mortgage-calculator-api.go .
COPY Gopkg.toml .

RUN dep ensure

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/mortgage-calculator-api .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/bin/mortgage-calculator-api .

EXPOSE 5000

CMD ["./mortgage-calculator-api"]