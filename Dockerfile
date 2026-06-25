# ------------------------------------------------------------------------------
# Builder Stage
# ------------------------------------------------------------------------------
FROM golang:1.25-trixie AS build

# Set default build argument for CGO_ENABLED
ARG CGO_ENABLED=0
ENV CGO_ENABLED=${CGO_ENABLED}

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download
COPY Makefile main.go ./
COPY static/ static/
COPY pkg/ pkg/

# The commit SHA is injected via a build arg so the build does not require the
# .git directory in the build context (which many CI/blue-green pipelines
# strip out). Leave empty to omit git stamping; build time and Go version are
# still stamped automatically since they don't depend on git.
ARG GIT_COMMIT=
RUN make build GIT_COMMIT="${GIT_COMMIT}"

# ------------------------------------------------------------------------------
# Fetch signing key
# ------------------------------------------------------------------------------
FROM debian:trixie-slim AS keyring
ADD https://www.postgresql.org/media/keys/ACCC4CF8.asc keyring.asc
RUN apt-get update && \
    apt-get install -qq --no-install-recommends gpg
RUN gpg -o keyring.pgp --dearmor keyring.asc

# ------------------------------------------------------------------------------
# Release Stage
# ------------------------------------------------------------------------------
FROM debian:trixie-slim

ARG keyring=/usr/share/keyrings/postgresql-archive-keyring.pgp
COPY --from=keyring /keyring.pgp $keyring
RUN . /etc/os-release && \
    echo "deb [signed-by=${keyring}] http://apt.postgresql.org/pub/repos/apt/ ${VERSION_CODENAME}-pgdg main" > /etc/apt/sources.list.d/pgdg.list && \
    apt-get update && \
    apt-get install -qq --no-install-recommends ca-certificates openssl netcat-openbsd curl postgresql-client

COPY --from=build /build/pgweb /usr/bin/pgweb

RUN useradd --uid 1000 --no-create-home --shell /bin/false pgweb
USER pgweb

EXPOSE 8081
ENTRYPOINT ["/usr/bin/pgweb", "--bind=0.0.0.0", "--listen=8081"]
