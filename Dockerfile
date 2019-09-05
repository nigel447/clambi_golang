# multi stage build
FROM golang:1
# creates and subsequent commands from here
WORKDIR /go/src/github.com/lambci/docker-lambda
RUN curl -sSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/local/bin/dep
COPY aws-lambda-mock.go Gopkg.toml Gopkg.lock ./
RUN dep ensure
RUN GOARCH=amd64 GOOS=linux go build aws-lambda-mock.go

# multi stage 2
# delete this image only for rebuild
FROM lambci/lambda-base

ENV AWS_EXECUTION_ENV=AWS_Lambda_go1.x

RUN rm -rf /var/runtime /var/lang && \
  curl https://lambci.s3.amazonaws.com/fs/go1.x.tgz | tar -zx -C /
# multi stage magic copys binaries from stage 1 only => from=0
COPY --from=0 /go/src/github.com/lambci/docker-lambda/aws-lambda-mock /var/runtime/aws-lambda-go

# mock aws lambda looks for binary here
COPY golami /var/task
USER sbx_user1051

ENTRYPOINT ["/var/runtime/aws-lambda-go"]
