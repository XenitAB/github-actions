# Image used for the build stage
FROM debian:stable-slim

# Set build arguments for target architecture
ARG TARGETARCH
ARG TARGETOS

# Install dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash \
    ca-certificates \
    git \
    curl \
    openssl \
    unzip \
    gpg \
    wget \
    jq \
    yq && \
    rm -rf /var/lib/apt/lists/*

# Install Azure CLI
RUN curl -sL https://aka.ms/InstallAzureCLIDeb | bash

# Install Terraform
RUN TERRAFORM_VERSION="1.11.0" && \
    curl -LO "https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip" && \
    unzip terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip -d /usr/local/bin && \
    rm terraform_${TERRAFORM_VERSION}_linux_${TARGETARCH}.zip

# Install tflint
RUN TFLINT_VERSION="v0.49.0" && \
    curl -L -o tflint.zip "https://github.com/terraform-linters/tflint/releases/download/${TFLINT_VERSION}/tflint_linux_${TARGETARCH}.zip" && \
    unzip tflint.zip -d /usr/local/bin && \
    rm tflint.zip

# Install tflint rulesets
COPY config/.tflint.hcl /work/.tflint.d/.tflint.hcl
RUN tflint --init --config=/work/.tflint.d/.tflint.hcl

# Install terraform (tfenv)
RUN TFENV_VERSION="1.5.7" && \
    git clone https://github.com/tfutils/tfenv.git ~/.tfenv && \
    echo 'export PATH="$HOME/.tfenv/bin:$PATH"' >> ~/.bashrc && \
    ~/.tfenv/bin/tfenv install ${TFENV_VERSION} && \
    ~/.tfenv/bin/tfenv use ${TFENV_VERSION}

# Install tfsec
RUN TFSEC_VERSION="v1.28.4" && \
    curl -L -o /usr/local/bin/tfsec "https://github.com/aquasecurity/tfsec/releases/download/${TFSEC_VERSION}/tfsec-linux-${TARGETARCH}" && \
    chmod +x /usr/local/bin/tfsec

# Install Open Policy Agent (OPA)
RUN OPA_VERSION="v0.43.0" && \
    curl -L -o opa "https://openpolicyagent.org/downloads/${OPA_VERSION}/opa_linux_${TARGETARCH}" && \
    chmod +x opa && \
    mv opa /usr/local/bin/

# Install sops
RUN SOPS_VERSION="v3.8.1" && \
    curl -L -o /usr/local/bin/sops "https://github.com/mozilla/sops/releases/download/${SOPS_VERSION}/sops-v${SOPS_VERSION}.linux.${TARGETARCH}" && \
    chmod +x /usr/local/bin/sops

# Install GitHub CLI
RUN GH_VERSION="2.39.2" && \
    curl -L -o gh.tar.gz "https://github.com/cli/cli/releases/download/v${GH_VERSION}/gh_${GH_VERSION}_linux_${TARGETARCH}.tar.gz" && \
    tar -xzf gh.tar.gz && \
    mv gh_*/bin/gh /usr/local/bin/ && \
    rm -rf gh_*

# Install kubectl
RUN KUBECTL_VERSION="v1.28.4" && \
    curl -LO "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/${TARGETARCH}/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/

# Install helm
RUN HELM_VERSION="v3.13.2" && \
    curl -fsSL -o get_helm.sh "https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3" && \
    chmod 700 get_helm.sh && \
    ./get_helm.sh --version "$HELM_VERSION" && \
    rm get_helm.sh

# Cleanup
RUN apt-get autoremove && \
    apt-get clean

# Copy additional files
COPY ./opa-policies /opt/opa-policies
COPY ./terraform.sh /opt/terraform.sh

# Set environment variables
ENV HOME=/work

# Set working directory
WORKDIR /work