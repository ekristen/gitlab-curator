FROM golang:1.14.7-alpine3.12 as builder
WORKDIR /src
ENV GO111MODULE=on
COPY . /src
RUN go build -o gitlab-curator main.go && chmod +x gitlab-curator

FROM alpine:3.12
ENTRYPOINT ["/usr/local/bin/gitlab-curator"]
COPY --from=builder /src/gitlab-curator /usr/local/bin/gitlab-curator
