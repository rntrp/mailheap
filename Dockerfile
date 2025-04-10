FROM golang:1.24-bookworm AS builder
WORKDIR /app
RUN echo "deb http://deb.debian.org/debian bookworm-backports main" | tee -a /etc/apt/sources.list \
    && apt update \
    && apt install upx-ucl -y
COPY . ./
COPY internal ./internal
RUN go mod download \
    && go test ./... \
    && go build -ldflags="-s -w -X 'github.com/rntrp/mailheap/internal/config.defaultEnv=production'" -o /mailheap \
    && upx --best --lzma /mailheap

FROM gcr.io/distroless/base-nossl-debian12:nonroot
COPY --from=builder /mailheap /
EXPOSE 2525 8080
ENTRYPOINT [ "/mailheap" ]
