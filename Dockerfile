FROM golang:1.25.3-trixie AS builder

ARG TARGETOS
ARG TARGETARCH

ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}

WORKDIR /go/src/github.com/kauche/cloud-run-api-emulator

COPY . .

RUN CGO_ENABLED=0 go build \
        -a \
        -trimpath \
        -ldflags "-s -w -extldflags -static" \
        -o /usr/bin/cloud-run-api-emulator \
        ./cmd/cloud-run-api-emulator

FROM gcr.io/distroless/base@sha256:df13a91fd415eb192a75e2ef7eacf3bb5877bb05ce93064b91b83feef5431f37
COPY --from=builder /usr/bin/cloud-run-api-emulator /usr/bin/cloud-run-api-emulator

CMD ["/usr/bin/cloud-run-api-emulator"]
