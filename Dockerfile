FROM golang:1.21-bullseye AS builder
WORKDIR /app
RUN apt update \
    && apt install upx-ucl -y
COPY . ./
COPY internal ./internal
RUN go mod download \
    && go test ./... \
    && go build -ldflags="-s -w" -o /mailheap \
    && upx --best --lzma /mailheap

FROM gcr.io/distroless/base-nossl-debian11:nonroot
COPY --from=builder /mailheap /
EXPOSE 2525 8080
ENV MAILHEAP_ENV production
ENTRYPOINT [ "/mailheap" ]
