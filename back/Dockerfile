FROM golang:1.24.3

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
    if [ "$TARGETARCH" = "amd64" ]; then \
    URL="https://7-zip.org/a/7z2500-linux-x64.tar.xz"; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
    URL="https://7-zip.org/a/7z2500-linux-arm64.tar.xz"; \
    elif [ "$TARGETARCH" = "arm" ]; then \
    URL="https://7-zip.org/a/7z2500-linux-arm.tar.xz"; \
    else \
    echo "Unsupported TARGETARCH: $TARGETARCH"; exit 1; \
    fi; \
    curl -L -o /tmp/7z.tar.xz "$URL"; \
    mkdir -p /usr/local/bin/7zz; \
    tar -xJf /tmp/7z.tar.xz -C /usr/local/bin/7zz --strip-components=0; \
    chmod +x /usr/local/bin/7zz/7zz; \
    ln -sf /usr/local/bin/7zz/7zz /usr/local/bin/7z; \
    rm /tmp/7z.tar.xz

RUN go mod tidy \
    && go build -o main .

CMD ["./main"]
