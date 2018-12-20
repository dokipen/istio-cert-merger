FROM golang:1.11-alpine3.8 AS build
RUN apk add curl && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . /go/src/github.com/dokipen/istio-cert-merger/
WORKDIR /go/src/github.com/dokipen/istio-cert-merger
RUN dep ensure && go install github.com/dokipen/istio-cert-merger/cmd/istio-cert-merger

FROM alpine:3.8
COPY --from=build /go/bin/istio-cert-merger /usr/bin/istio-cert-merger
ENTRYPOINT ["/usr/bin/istio-cert-merger"]
