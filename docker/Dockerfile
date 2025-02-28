# Image used for the build stage
FROM debian:stable-slim

# Set build arguments for target architecture
ARG TARGETARCH
ARG TARGETOS


# Install dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    # ansible \
    ca-certificates \
    git \
    curl \
    openssl \
    # pip \
    unzip \
    gpg \
    wget \
    jq \
    yq && \
    rm -rf /var/lib/apt/lists/*

# Install Azure CLI
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash

# Install tflint
RUN if [ "$TARGETARCH" = "amd64" ]; then \
      wget https://github.com/terraform-linters/tflint/releases/download/v0.49.0/tflint_linux_amd64.zip -O tflint.zip; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
      wget https://github.com/terraform-linters/tflint/releases/download/v0.49.0/tflint_linux_arm64.zip -O tflint.zip; \
    fi && \
    unzip tflint.zip -d /usr/local/bin && \
    rm tflint.zip

# Install tflint rulesets
COPY config/.tflint.hcl /work/.tflint.d/.tflint.hcl
RUN tflint --init --config=/work/.tflint.d/.tflint.hcl

# Install terraform (tfenv)
RUN git clone https://github.com/tfutils/tfenv.git ~/.tfenv && \
    echo 'export PATH="$HOME/.tfenv/bin:$PATH"' >> ~/.bashrc && \
    ~/.tfenv/bin/tfenv install 1.5.7 && \
    ~/.tfenv/bin/tfenv use 1.5.7

# Install tfsec
RUN if [ "$TARGETARCH" = "amd64" ]; then \
      wget https://github.com/aquasecurity/tfsec/releases/download/v1.28.4/tfsec-linux-amd64 -O /usr/local/bin/tfsec; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
      wget https://github.com/aquasecurity/tfsec/releases/download/v1.28.4/tfsec-linux-arm64 -O /usr/local/bin/tfsec; \
    fi && \
    chmod +x /usr/local/bin/tfsec

# Install Open Policy Agent (OPA)
RUN if [ "$TARGETARCH" = "amd64" ]; then \
      curl -L -o opa https://openpolicyagent.org/downloads/v0.43.0/opa_linux_amd64; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
      curl -L -o opa https://openpolicyagent.org/downloads/v0.43.0/opa_linux_arm64; \
    fi && \
    chmod +x opa && \
    mv opa /usr/local/bin/

# Install sops
RUN if [ "$TARGETARCH" = "amd64" ]; then \
      wget https://github.com/mozilla/sops/releases/download/v3.8.1/sops-v3.8.1.linux.amd64 -O /usr/local/bin/sops; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
      wget https://github.com/mozilla/sops/releases/download/v3.8.1/sops-v3.8.1.linux.arm64 -O /usr/local/bin/sops; \
    fi && \
    chmod +x /usr/local/bin/sops

# Install GitHub CLI
RUN if [ "$TARGETARCH" = "amd64" ]; then \
      wget https://github.com/cli/cli/releases/download/v2.39.2/gh_2.39.2_linux_amd64.tar.gz; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
      wget https://github.com/cli/cli/releases/download/v2.39.2/gh_2.39.2_linux_arm64.tar.gz; \
    fi && \
    tar -xzf gh_*.tar.gz && \
    mv gh_*/bin/gh /usr/local/bin/ && \
    rm -rf gh_*

# Install kubectl
RUN if [ "$TARGETARCH" = "amd64" ]; then \
      curl -LO "https://dl.k8s.io/release/v1.28.4/bin/linux/amd64/kubectl"; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
      curl -LO "https://dl.k8s.io/release/v1.28.4/bin/linux/arm64/kubectl"; \
    fi && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/

# Install helm
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 && \
    chmod 700 get_helm.sh && \
    ./get_helm.sh --version v3.13.2 && \
    rm get_helm.sh

# Cleanup
RUN apt-get autoremove && \
    apt-get clean

RUN rm -rf /tmp/install

# Copy additional files
COPY ./opa-policies /opt/opa-policies
# COPY ./terraform.sh /opt/terraform.sh

# Copy the terraform.sh script and config.yaml
COPY ./tf-prepare/sh/tf-prepare.sh /opt/tf-prepare.sh
COPY ./tf-wrapper.sh /opt/tf-wrapper.sh

# COPY ./sh/config.yaml /opt/config.yaml

# Make the script executable
RUN chmod +x /opt/tf-prepare.sh
RUN chmod +x /opt/tf-wrapper.sh

# Set environment variables
ENV HOME=/work

# Set working directory
WORKDIR /work

# Set the entrypoint to run the terraform.sh script
# ENTRYPOINT ["/opt/tf-wrapper.sh"]
