FROM golang:1.19 as tf-prepare-builder
WORKDIR /workspace

COPY ./go-tf-prepare/go.mod ./go-tf-prepare/go.sum ./
RUN go mod download
COPY ./go-tf-prepare/main.go main.go
COPY ./go-tf-prepare/pkg/ pkg/
RUN GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o tf-prepare main.go

FROM debian:11.1-slim

#Base
RUN apt-get update -y && \
    apt-get install -y git curl openssl pip make unzip gpg wget && \
    apt-get install -y ansible=2.10.7+merged+base+2.10.8+dfsg-1 && \
    apt-get install -y packer=1.6.6+ds1-2+b3

RUN pip --no-cache-dir install -U pip

RUN mkdir -p /tmp/install /usr/src /work
WORKDIR /tmp/install

# Install Azure CLI
COPY install-scripts/azure-cli.sh /usr/src/install-scripts/azure-cli.sh
RUN /usr/src/install-scripts/azure-cli.sh --version="2.30.0"

# Install AWS CLI
COPY install-scripts/aws-cli.sh /usr/src/install-scripts/aws-cli.sh
RUN /usr/src/install-scripts/aws-cli.sh --version="2.1.10"

# Install tflint
COPY install-scripts/tflint.sh /usr/src/install-scripts/tflint.sh
RUN /usr/src/install-scripts/tflint.sh --version="v0.34.1" --sha="5a3d254c5e8222e15afea5da61a7ab66d560c31f302a47cbb90ed69c6dc973ec"
COPY config/.tflint.hcl /work/.tflint.d/.tflint.hcl

# Install tflint ruleset
COPY install-scripts/tflint-ruleset.sh /usr/src/install-scripts/tflint-ruleset.sh
RUN /usr/src/install-scripts/tflint-ruleset.sh --ruleset="azurerm" --version="v0.14.0" --sha="63a1e9b2502516144d62643d37e2ff6bba1e795305b0e7e7d6a8bb8f0cabf2b2"
RUN /usr/src/install-scripts/tflint-ruleset.sh --ruleset="aws" --version="v0.12.0" --sha="1f8e54f5504c91f1bba4cf47ad749d8bec64b7f1786837f3d39921246bdfe15a"

# Install terraform (tfenv)
COPY install-scripts/tfenv.sh /usr/src/install-scripts/tfenv.sh
RUN /usr/src/install-scripts/tfenv.sh --latest-terraform-version="1.1.7" --tfenv-version="v2.2.2"

# Install tfsec
COPY install-scripts/tfsec.sh /usr/src/install-scripts/tfsec.sh
RUN /usr/src/install-scripts/tfsec.sh --version="v1.8.0" --sha="6f0ffcc204a73991fb2b47561bbdf473e98a9566036c0f36b95d9574fd5ab4ce"

# Install Open Policy Agent
COPY install-scripts/opa.sh /usr/src/install-scripts/opa.sh
RUN /usr/src/install-scripts/opa.sh --version="v0.40.0" --sha="fc77b0fd2efffb60f0d65d22c251b780c15a8e298d5034f41439a7f3d3183f97"

# Install sops
COPY install-scripts/sops.sh /usr/src/install-scripts/sops.sh
RUN /usr/src/install-scripts/sops.sh --version="v3.6.1" --sha="b2252aa00836c72534471e1099fa22fab2133329b62d7826b5ac49511fcc8997"

# Install GitHub CLI
COPY install-scripts/github-cli.sh /usr/src/install-scripts/github-cli.sh
RUN /usr/src/install-scripts/github-cli.sh --version="1.3.0" --sha="56e540ddc978908fd236d53b00855c3526936392976e18bf429161963bbd45ec"

# Install jq
COPY install-scripts/jq.sh /usr/src/install-scripts/jq.sh
RUN /usr/src/install-scripts/jq.sh --version="1.6" --sha="af986793a515d500ab2d35f8d2aecd656e764504b789b66d7e1a0b727a124c44"

# Install yq
COPY install-scripts/yq.sh /usr/src/install-scripts/yq.sh
RUN /usr/src/install-scripts/yq.sh --version="2.11.1"

# Install kubectl
COPY install-scripts/kubectl.sh /usr/src/install-scripts/kubectl.sh
RUN /usr/src/install-scripts/kubectl.sh --version="v1.19.0" --sha="79bb0d2f05487ff533999a639c075043c70a0a1ba25c1629eb1eef6ebe3ba70f"

# Install helm
COPY install-scripts/helm.sh /usr/src/install-scripts/helm.sh
RUN /usr/src/install-scripts/helm.sh --version="v3.4.1" --sha="538f85b4b73ac6160b30fd0ab4b510441aa3fa326593466e8bf7084a9c288420"

COPY --from=tf-prepare-builder /workspace/tf-prepare /usr/local/bin/tf-prepare
RUN chmod +x /usr/local/bin/tf-prepare

#Cleanup
RUN apt-get autoremove && \
    apt-get clean

RUN rm -rf /tmp/install

COPY opa-policies /opt/opa-policies
COPY terraform.sh /opt/terraform.sh
COPY packer.sh /opt/packer.sh

ENV HOME=/work

WORKDIR /work
