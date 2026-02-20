# Image used for the go-tf-prepare stage
FROM golang:1.24 AS tf-prepare-builder
WORKDIR /workspace

COPY ./go-tf-prepare/go.mod ./go-tf-prepare/go.sum ./
RUN go mod download
COPY ./go-tf-prepare/main.go main.go
COPY ./go-tf-prepare/pkg/ pkg/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} GO111MODULE=on go build -ldflags='-s -w' -o tf-prepare main.go

# ------------------------------

# Image used for the build stage
FROM debian:stable-slim

# Set build arguments for target architecture
ARG TARGETARCH
ARG TARGETOS

# MKDIR
RUN mkdir -p /work

# Install dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    apt-transport-https \
    ca-certificates \
    curl \
    git \
    make \
    unzip \
    gpg \
    gnupg \
    lsb-release \
    jq && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /etc/apt/keyrings

# Install Azure CLI
RUN AZ_CLI_VERSION="2.83.0" && \
    AZ_CLI_DISTRO="$(lsb_release -cs)" && \
    if [ "${AZ_CLI_DISTRO}" = "trixie" ]; then AZ_CLI_DISTRO="bookworm"; fi && \
    curl -sSL https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor -o /etc/apt/keyrings/microsoft.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/microsoft.gpg] https://packages.microsoft.com/repos/azure-cli/ ${AZ_CLI_DISTRO} main" \
    > /etc/apt/sources.list.d/azure-cli.list && \
    apt-get update && \
    if apt-cache madison azure-cli | awk '{ print $3 }' | grep -qx "${AZ_CLI_VERSION}-1~${AZ_CLI_DISTRO}"; then \
    apt-get install -y --no-install-recommends "azure-cli=${AZ_CLI_VERSION}-1~${AZ_CLI_DISTRO}"; \
    else \
    apt-get install -y --no-install-recommends azure-cli; \
    fi && \
    rm -rf /var/lib/apt/lists/*

# Install Terraform
RUN TERRAFORM_VERSION="1.14.5" && \
    curl -L -o terraform.zip "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip" && \
    unzip terraform.zip -d /usr/local/bin && \
    rm terraform.zip

# Install tflint
RUN TFLINT_VERSION="v0.61.0" && \
    curl -L -o tflint.zip "https://github.com/terraform-linters/tflint/releases/download/${TFLINT_VERSION}/tflint_linux_${TARGETARCH}.zip" && \
    unzip tflint.zip -d /usr/local/bin && \
    rm tflint.zip && \
    mkdir -p /work/.tflint.d/plugins/

# Install tflint rulesets && copy tflint configuration
COPY config/.tflint.hcl /work/.tflint.d/.tflint.hcl

RUN AZURERM_RULESET_VERSION="v0.31.0" && \
    curl -L -o tflint-ruleset-azurerm.zip "https://github.com/terraform-linters/tflint-ruleset-azurerm/releases/download/${AZURERM_RULESET_VERSION}/tflint-ruleset-azurerm_linux_${TARGETARCH}.zip" && \
    unzip tflint-ruleset-azurerm.zip -d /work/.tflint.d/plugins/ && \
    rm tflint-ruleset-azurerm.zip

# Install tfsec
RUN TFSEC_VERSION="v1.28.14" && \
    curl -L -o /usr/local/bin/tfsec "https://github.com/aquasecurity/tfsec/releases/download/${TFSEC_VERSION}/tfsec-linux-${TARGETARCH}" && \
    chmod +x /usr/local/bin/tfsec

# Install Open Policy Agent (OPA)
RUN OPA_VERSION="v1.13.2" && \
    curl -L -o /usr/local/bin/opa "https://openpolicyagent.org/downloads/${OPA_VERSION}/opa_linux_${TARGETARCH}_static" && \
    chmod +x /usr/local/bin/opa

# Install sopsh
RUN SOPS_VERSION="v3.11.0" && \
    curl -L -o /usr/local/bin/sops "https://github.com/getsops/sops/releases/download/${SOPS_VERSION}/sops-${SOPS_VERSION}.linux.${TARGETARCH}" && \
    chmod +x /usr/local/bin/sops

# Cleanup
RUN apt-get autoremove && \
    apt-get clean

# Copy additional files
COPY ./opa-policies /opt/opa-policies
COPY ./terraform.sh /opt/terraform.sh
COPY --from=tf-prepare-builder /workspace/tf-prepare /usr/local/bin/tf-prepare
RUN chmod +x /usr/local/bin/tf-prepare

# Set environment variables
ENV HOME=/work

# Set working directory
WORKDIR /work