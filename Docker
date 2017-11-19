FROM golang

ENV API_INSECURE false

RUN apt-get install -y git && \
    go get github.com/czerwonk/ovirt-zero-touch

CMD ovirt-zero-touch -api-url=$API_URL -username=$API_USER -password=$API_PASS -insecure=$API_INSECURE
EXPOSE 11337
