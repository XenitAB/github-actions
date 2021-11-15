FROM golang:1.15 as tf-prepare-builder
WORKDIR /workspace

COPY ./go-tf-prepare/go.mod ./go-tf-prepare/go.sum ./
RUN go mod download
COPY ./go-tf-prepare/main.go main.go
COPY ./go-tf-prepare/pkg/ pkg/
RUN GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o tf-prepare main.go

FROM debian:11.1-slim

#Base
RUN apt update -y && \
    apt install -y git curl openssl && \
    apt install -y pip && \
    apt install -y make unzip gpg wget && \
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
RUN /usr/src/install-scripts/tflint.sh --version="v0.24.1" --sha="2dbe3b423f5d3e0bb458d51761c97d51a4fd6c3d7bd1efd87c4aa3dc5199e7b2"
COPY config/.tflint.hcl /work/.tflint.d/.tflint.hcl

# Install tflint ruleset
COPY install-scripts/tflint-ruleset.sh /usr/src/install-scripts/tflint-ruleset.sh
RUN /usr/src/install-scripts/tflint-ruleset.sh --ruleset="azurerm" --version="v0.8.2" --sha="4ef97bbc847bde194401c3206eb127fffaf4ce430127e0408878a8a833242a30"
RUN /usr/src/install-scripts/tflint-ruleset.sh --ruleset="aws" --version="v0.2.1" --sha="ec2a992a8413227e2321d985b62cde34bc34287599894f966b0fc8904aba0d8a"

# Install terraform (tfenv)
COPY install-scripts/tfenv.sh /usr/src/install-scripts/tfenv.sh
RUN /usr/src/install-scripts/tfenv.sh --latest-terraform-version="0.15.3" --tfenv-version="v2.2.2"

# Install tfsec
COPY install-scripts/tfsec.sh /usr/src/install-scripts/tfsec.sh
RUN /usr/src/install-scripts/tfsec.sh --version="v0.39.5" --sha="60e52ef9a2b2eb5aebf74fbebfbaf0d30fa107816a8bbc2759cfe5d5c2a9021d"

# Install Open Policy Agent
COPY install-scripts/opa.sh /usr/src/install-scripts/opa.sh
RUN /usr/src/install-scripts/opa.sh --version="v0.24.0" --sha="e40bde4cca8a5819518e3c35862bc5b6c388bc2904d412227059af29170f79e9"

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
