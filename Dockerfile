FROM golang:1.22.10-bullseye 

ENV PATH="/root/.cargo/bin:${PATH}" 
ENV USER=root

# Instalar Rust
RUN apt-get update && apt-get install -y --no-install-recommends curl && \
    curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs -o rustup-init.sh && \
    sh rustup-init.sh -y && \
    . $HOME/.cargo/env && rustup default stable && \
    rustup target add x86_64-unknown-linux-gnu && \
    rm -rf /var/lib/apt/lists/* rustup-init.sh

WORKDIR /go/src

# Copiar os arquivos de módulo Go
RUN ln  -sf /bin/bash /bin/sh
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copiar o restante dos arquivos do aplicativo
COPY . .

# Configurar o comando de execução
CMD ["tail", "-f", "/dev/null"]
