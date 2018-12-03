FROM golang:1.11 as builder
RUN mkdir -p /home/ayesha/Desktop/rest
WORKDIR /home/ayesha/Desktop/rest
COPY restapi.go .
RUN go get github.com/gorilla/mux
RUN CGO_ENABLE=0 GOOS=linux go build -o restserver .


FROM golang:1.11

COPY --from=builder /home/ayesha/Desktop/rest/restserver /usr/bin/restserver

ENTRYPOINT ["restserver"]

