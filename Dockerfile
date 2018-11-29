FROM golang:1.11

WORKDIR ~/Desktop/rest
COPY .git     .
COPY dockertest.go   .

RUN GIT_COMMIT=$(git rev-list -1 HEAD) && \
  go build -ldflags "-X main.GitCommit=$GIT_COMMIT"

CMD ["./rest"]
