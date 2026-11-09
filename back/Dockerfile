FROM golang:1.27.3

ARG TARGETARCH
ENV TARGETARCH=${TARGETARCH}

WORKDIR /back
COPY . .

RUN apt-get update && apt-get install -y --no-install-recommends \
    git \
    poppler-utils \
    ghostscript \
    curl \
    xz-utils \
    && rm -rf /var/lib/apt/lists/*

RUN set -eux; \
    case "${TARGETARCH}" in \
      amd64) FILE="7z2501-linux-x64.tar.xz" ;; \
      arm64) FILE="7z2501-linux-arm64.tar.xz" ;; \
      arm) FILE="7z2501-linux-arm.tar.xz" ;; \
      *) echo "Unsupported TARGETARCH: ${TARGETARCH}"; exit 1 ;; \
    esac; \
    URL="https://7-zip.org/a/${FILE}"; \
    echo "$URL"; \
    curl -fsSL "$URL" -o /tmp/7z.tar.xz; \
    mkdir -p /tmp/7zextract; \
    tar -xf /tmp/7z.tar.xz -C /tmp/7zextract; \
    install -m 755 /tmp/7zextract/7zz /usr/local/bin/7zz; \
    ln -sf /usr/local/bin/7zz /usr/local/bin/7z; \
    rm -rf /tmp/7z.tar.xz /tmp/7zextract

RUN go mod tidy \
    && go build -o main .

CMD ["./main"]
