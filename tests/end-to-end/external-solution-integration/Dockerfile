FROM golang:1.13-alpine as builder

ENV SRC_DIR=/go/src/github.com/kyma-project/kyma/tests/end-to-end/external-solution-integration

WORKDIR $SRC_DIR
COPY . $SRC_DIR

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /e2e ./cmd/runner

FROM alpine:edge
LABEL source=git@github.com:kyma-project/kyma.git

COPY --from=builder /e2e /

RUN apk update && apk add ca-certificates && apk add curl && rm -rf /var/cache/apk/*

ENTRYPOINT ["/e2e"]
