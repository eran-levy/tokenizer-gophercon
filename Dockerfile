FROM golang:1.15 as builder
WORKDIR /build
COPY go.* ./
RUN go mod download
#Build
COPY . .

# `skaffold debug` sets SKAFFOLD_GO_GCFLAGS to disable compiler optimizations
ARG SKAFFOLD_GO_GCFLAGS
RUN GOOS=linux CGO_ENABLED=0 go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -a -o /app ./cmd/tokenizer

FROM golang:1.15.8-alpine3.13
WORKDIR /opt
# Define GOTRACEBACK to mark this container as using the Go language runtime
# for `skaffold debug` (https://skaffold.dev/docs/workflows/debug/).
ENV GOTRACEBACK=single
EXPOSE 8080 7070

COPY --from=builder /app app
COPY --from=builder /build/helm helm

CMD ["/opt/app"]