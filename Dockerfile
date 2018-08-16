FROM golang as builder
RUN go get github.com/czerwonk/ovirt-zero-touch


FROM alpine:latest

ENV API_INSECURE false
ENV DEBUG false

RUN apk --no-cache add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/bin/ovirt-zero-touch .

CMD ovirt-zero-touch -debug=$DEBUG -insecure=$API_INSECURE -api-url=$API_URL -username=$API_USER -password=$API_PASS
EXPOSE 11337
