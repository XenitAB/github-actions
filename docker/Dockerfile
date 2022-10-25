FROM golang:1.19 as tf-prepare-builder
WORKDIR /workspace

COPY ./go-tf-prepare/go.mod ./go-tf-prepare/go.sum ./
RUN go mod download
COPY ./go-tf-prepare/main.go main.go
COPY ./go-tf-prepare/pkg/ pkg/
RUN GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o tf-prepare main.go

FROM debian:11.4-slim

#Base
RUN apt-get update -y && \
    apt-get install -y git curl openssl pip make unzip gpg wget && \
    apt-get install -y ansible=2.10.7+merged+base+2.10.8+dfsg-1

RUN pip --no-cache-dir install -U pip

RUN mkdir -p /tmp/install /usr/src /work
WORKDIR /tmp/install

# Install Azure CLI
COPY install-scripts/azure-cli.sh /usr/src/install-scripts/azure-cli.sh
RUN /usr/src/install-scripts/azure-cli.sh --version="2.39.0"

# Install AWS CLI
COPY install-scripts/aws-cli.sh /usr/src/install-scripts/aws-cli.sh
RUN /usr/src/install-scripts/aws-cli.sh --version="2.7.21"

# Install tflint
COPY install-scripts/tflint.sh /usr/src/install-scripts/tflint.sh
RUN /usr/src/install-scripts/tflint.sh --version="v0.42.1" --sha="11294cb2f395182fefdc7d35f49fc18966dda8860fcea1ecd67b6ab190ab7ea0"
COPY config/.tflint.hcl /work/.tflint.d/.tflint.hcl

# Install tflint ruleset
COPY install-scripts/tflint-ruleset.sh /usr/src/install-scripts/tflint-ruleset.sh
RUN /usr/src/install-scripts/tflint-ruleset.sh --ruleset="azurerm" --version="v0.18.0" --sha="39e3f45a7d00222abd151ec2dc0909d84cc5503bbb187d5a500a27643c1e8483"
RUN /usr/src/install-scripts/tflint-ruleset.sh --ruleset="aws" --version="v0.18.0" --sha="433916290a2603e0a388b8f05f4d462cec765d45222228e12e411b7007988916"

# Install terraform (tfenv)
COPY install-scripts/tfenv.sh /usr/src/install-scripts/tfenv.sh
RUN /usr/src/install-scripts/tfenv.sh --latest-terraform-version="1.3.0" --tfenv-version="v3.0.0"

# Install tfsec
COPY install-scripts/tfsec.sh /usr/src/install-scripts/tfsec.sh
RUN /usr/src/install-scripts/tfsec.sh --version="v1.27.1" --sha="edf06ce4897a3113dda6393b31345aea8b70626dac4c67df87ef0b69fd6c83f0"

# Install Open Policy Agent
COPY install-scripts/opa.sh /usr/src/install-scripts/opa.sh
RUN /usr/src/install-scripts/opa.sh --version="v0.43.0" --sha="d5337139a7ccb04149bd9f96ab7a1641a3e3c39f6e1fffa610c7a5c054b0881f"

# Install sops
COPY install-scripts/sops.sh /usr/src/install-scripts/sops.sh
RUN /usr/src/install-scripts/sops.sh --version="v3.7.3" --sha="53aec65e45f62a769ff24b7e5384f0c82d62668dd96ed56685f649da114b4dbb"

# Install GitHub CLI
COPY install-scripts/github-cli.sh /usr/src/install-scripts/github-cli.sh
RUN /usr/src/install-scripts/github-cli.sh --version="2.14.3" --sha="95fe2e93bbd7a9f02d5b22cee01a53ab6d581fda8fa170364b668d142840ef58"

# Install jq
COPY install-scripts/jq.sh /usr/src/install-scripts/jq.sh
RUN /usr/src/install-scripts/jq.sh --version="1.6" --sha="af986793a515d500ab2d35f8d2aecd656e764504b789b66d7e1a0b727a124c44"

# Install yq
COPY install-scripts/yq.sh /usr/src/install-scripts/yq.sh
RUN /usr/src/install-scripts/yq.sh --version="3.1.0"

# Install kubectl
COPY install-scripts/kubectl.sh /usr/src/install-scripts/kubectl.sh
RUN /usr/src/install-scripts/kubectl.sh --version="v1.24.3" --sha="8a45348bdaf81d46caf1706c8bf95b3f431150554f47d444ffde89e8cdd712c1"

# Install helm
COPY install-scripts/helm.sh /usr/src/install-scripts/helm.sh
RUN /usr/src/install-scripts/helm.sh --version="v3.9.2" --sha="3f5be38068a1829670440ccf00b3b6656fd90d0d9cfd4367539f3b13e4c20531"

COPY --from=tf-prepare-builder /workspace/tf-prepare /usr/local/bin/tf-prepare
RUN chmod +x /usr/local/bin/tf-prepare

# Install packer
COPY install-scripts/packer.sh /usr/src/install-scripts/packer.sh
RUN /usr/src/install-scripts/packer.sh --version="1.8.3" --sha="0587f7815ed79589cd9c2b754c82115731c8d0b8fd3b746fe40055d969facba5"

#Cleanup
RUN apt-get autoremove && \
    apt-get clean

RUN rm -rf /tmp/install

COPY opa-policies /opt/opa-policies
COPY terraform.sh /opt/terraform.sh
COPY packer.sh /opt/packer.sh

ENV HOME=/work

WORKDIR /work
