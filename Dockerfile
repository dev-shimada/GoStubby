FROM --platform=$BUILDPLATFORM golang:1.24.0-bookworm AS build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY . /workspace

RUN  <<EOF
CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bin/gostubby main.go
EOF

FROM --platform=$BUILDPLATFORM gcr.io/distroless/base-debian12:latest
WORKDIR /app
COPY --chmod=100 --from=build /workspace/bin/gostubby /app/gostubby
ENTRYPOINT [ "/app/gostubby" ]
