FROM golang:1.21-bullseye AS builder
WORKDIR /app
RUN apt update \
    && apt install upx-ucl
COPY . ./
COPY internal ./internal
RUN go mod download \
    && go test ./... \
    && go build -ldflags="-s -w" -o /mailheap \
    && upx --best --lzma /mailheap

FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=builder /mailheap ./
EXPOSE 8080
ENTRYPOINT [ "/mailheap" ]
