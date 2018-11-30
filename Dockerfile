FROM golang:1.11 as builder
RUN mkdir -p /go/src/github.com/Ayesha049/Rest-API
WORKDIR /go/src/github.com/Ayesha049/Rest-API
COPY restapi.go .
RUN go get github.com/gorilla/mux
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM busybox
COPY --from=builder /go/src/github.com/Ayesha049/Rest-API/server /usr/bin/server
ENTRYPOINT ["server"]
