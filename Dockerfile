# syntax=docker/dockerfile:experimental
# ---
FROM golang:1.23 AS build

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

WORKDIR /work
COPY . /work

# Build container-admission-webhook
RUN --mount=type=cache,target=/root/.cache/go-build,sharing=private \
    go build -o container-admission-webhook .

# ---
FROM scratch AS run

COPY --from=build /work/container-admission-webhook /usr/local/bin/

CMD ["container-admission-webhook"]