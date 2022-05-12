FROM golang:1.15 as tf-prepare-builder
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
RUN /usr/src/install-scripts/tflint.sh --version="v0.35.0" --sha="9969323f8c6fd511bdd9267a34d89c5941eb9a16336fec368381555e9c2f2441"
COPY config/.tflint.hcl /work/.tflint.d/.tflint.hcl

# Install tflint ruleset
COPY install-scripts/tflint-ruleset.sh /usr/src/install-scripts/tflint-ruleset.sh
RUN /usr/src/install-scripts/tflint-ruleset.sh --ruleset="azurerm" --version="v0.15.0" --sha="995cf85133d7512bc36e060a235b2edb332520d21beba1fa05a9e092f05c5be4"
RUN /usr/src/install-scripts/tflint-ruleset.sh --ruleset="aws" --version="v0.13.3" --sha="51b35974f558e3fdd03b5eabb922bf7dc99027e85ba3136dddc254a210c478f5"

# Install terraform (tfenv)
COPY install-scripts/tfenv.sh /usr/src/install-scripts/tfenv.sh
RUN /usr/src/install-scripts/tfenv.sh --latest-terraform-version="1.1.7" --tfenv-version="v2.2.2"

# Install tfsec
COPY install-scripts/tfsec.sh /usr/src/install-scripts/tfsec.sh
RUN /usr/src/install-scripts/tfsec.sh --version="v1.18.0" --sha="15bc4bec72b8e7bffa283fc1be4c55f6fc35eaa4a09d49126f58d9be6f4bb9ad"

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
