FROM --platform=$BUILDPLATFORM golang:1.24.3-bookworm AS build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY go.mod /workspace
COPY go.sum /workspace
RUN go mod download
COPY internal /workspace/internal
COPY main.go /workspace
COPY LICENSE /workspace

RUN  <<EOF
CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bin/gostubby main.go
EOF

FROM --platform=$BUILDPLATFORM gcr.io/distroless/base-debian12:latest
WORKDIR /app
COPY --chmod=100 --from=build /workspace/bin/gostubby /app/gostubby
ENTRYPOINT [ "/app/gostubby" ]
CMD [ "--host", "0.0.0.0" ]
