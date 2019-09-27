FROM golang:1.12-alpine
WORKDIR /go/src/github.com/inloco/sops-kustomize-generator-plugin
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on
RUN apk add git && \
    go get -d -v ./... && \
    go install -a -installsuffix cgo -ldflags '-extldflags "-static" -s -w' -tags netgo -v ./...

FROM alpine:3.8
COPY --from=0 /go/bin/sops-kustomize-generator-plugin /root/.config/kustomize/plugin/inloco.com.br/v1beta1/sops/SOPS
CMD [ "/root/.config/kustomize/plugin/inloco.com.br/v1beta1/sops/SOPS" ]
