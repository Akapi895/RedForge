# ============================================================================
# CyberStrikeAI — Docker image (multi-stage)
#
# Build:      docker build -t cyberstrike-ai .
# Run:        docker compose up -d     (preferred, see docker-compose.yml)
#
# Notes:
#   - The server serves the SPA from disk (./web/static) and resolves tools/,
#     skills/, roles/, agents/ relative to the working directory, so the
#     runtime image runs from /app with those trees present.
#   - SQLite (mattn/go-sqlite3) needs cgo, so the build stage keeps CGO enabled
#     and the runtime image ships glibc + python3.
#   - Persistence is handled by mounting a volume at /app/data and a config
#     volume at /app/config.yaml (see docker-compose.yml).
# ============================================================================

# ---------------------------------------------------------------------------
# Stage 1: build Go binaries (server + admin-password helper)
# ---------------------------------------------------------------------------
FROM golang:1.25-bookworm AS builder

# goproxy.cn is what the repo's own run.sh defaults to; proxy.golang.org is
# often unreliable for this project's region. Override with --build-arg GOPROXY.
ARG GOPROXY="https://goproxy.cn,direct"
ENV CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="$GOPROXY"

WORKDIR /src

# Dependencies first for better layer caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build both binaries.
COPY cmd ./cmd
COPY internal ./internal

# Some antivirus products delete legitimate C2 source files (false positives),
# which would otherwise fail the build with cryptic "undefined: c2.OnelinerKind"
# errors. Detect it early and point the user at the portable restore helper.
RUN if [ ! -f internal/c2/payload_oneliner.go ]; then \
      echo "ERROR: internal/c2/payload_oneliner.go is missing."; \
      echo "This is usually an antivirus false positive that deleted the file."; \
      echo "Run from the repo root to restore it, then rebuild:"; \
      echo "    ./docker/prepare.sh"; \
      exit 1; \
    fi \
 && go build -trimpath -o /out/cyberstrike-ai      ./cmd/server \
 && go build -trimpath -o /out/set-admin-password ./cmd/set-admin-password

# ---------------------------------------------------------------------------
# Stage 2: runtime image
# ---------------------------------------------------------------------------
FROM debian:bookworm-slim

# git is used by several tool recipes; ca-certificates + curl for HTTPS and
# healthchecks. The binary itself needs libc (from libc6, always present) and
# the SQLite cgo build may pull libgcc.
RUN apt-get update \
 && apt-get install -y --no-install-recommends \
        python3 \
        python3-venv \
        python3-pip \
        ca-certificates \
        curl \
        git \
        sqlite3 \
        libgcc-s1 \
 && rm -rf /var/lib/apt/lists/*

ENV PIP_BREAK_SYSTEM_PACKAGES=1 \
    PYTHONUNBUFFERED=1 \
    GOPROXY="https://proxy.golang.org,direct"

WORKDIR /app

# Copy the compiled Go binaries.
COPY --from=builder /out/cyberstrike-ai      ./cyberstrike-ai
COPY --from=builder /out/set-admin-password ./set-admin-password

# Runtime content that the server reads from disk / resolves relative to /app.
COPY requirements.txt ./
COPY config.example.yaml ./
COPY web/        ./web/
COPY tools/      ./tools/
COPY skills/     ./skills/
COPY roles/      ./roles/
COPY agents/     ./agents/
COPY knowledge_base/ ./knowledge_base/
COPY prompts/      ./prompts/

# Python virtual environment for tool recipes that need it (venv/ is what the
# server expects, matching run.sh).
RUN python3 -m venv /app/venv \
 && /app/venv/bin/pip install --no-cache-dir -r /app/requirements.txt

# Entrypoint handles config bootstrap + optional default admin password.
COPY docker/entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Data directory is bind-mounted in compose; create it so WAL/artifacts work.
RUN mkdir -p /app/data
VOLUME ["/app/data"]

EXPOSE 7123 7134

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["--https"]
