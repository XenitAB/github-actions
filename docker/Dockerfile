FROM golang:1.21 as tf-prepare-builder
WORKDIR /workspace

COPY ./go-tf-prepare/go.mod ./go-tf-prepare/go.sum ./
RUN go mod download
COPY ./go-tf-prepare/main.go main.go
COPY ./go-tf-prepare/pkg/ pkg/
RUN GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o tf-prepare main.go

FROM debian:12.2-slim

#Base
RUN apt-get update -y
RUN apt-get install -y git curl openssl pip make unzip gpg wget apt-utils
RUN apt-get install -y ansible=7.3.0+dfsg-1

RUN mkdir -p /tmp/install /usr/src /work
WORKDIR /tmp/install

# Install Azure CLI
COPY install-scripts/azure-cli.sh /usr/src/install-scripts/azure-cli.sh
RUN /usr/src/install-scripts/azure-cli.sh --version="2.55.0"

# Install AWS CLI
COPY install-scripts/aws-cli.sh /usr/src/install-scripts/aws-cli.sh
RUN /usr/src/install-scripts/aws-cli.sh --version="2.14.5"

# Install tflint
COPY install-scripts/tflint.sh /usr/src/install-scripts/tflint.sh
RUN /usr/src/install-scripts/tflint.sh --version="v0.49.0" --sha="56d862054e8f71e3ba392c6617fca2d730b5ecdf8a4a0768ba6087406dcc7d63"
COPY config/.tflint.hcl /work/.tflint.d/.tflint.hcl

# Install tflint ruleset
COPY install-scripts/tflint-ruleset.sh /usr/src/install-scripts/tflint-ruleset.sh
RUN /usr/src/install-scripts/tflint-ruleset.sh --ruleset="azurerm" --version="v0.25.1" --sha="a0812b02b93dfbfcfb9521a665be249b025ee074f60943e5b5eb1869dbb36ee0"
RUN /usr/src/install-scripts/tflint-ruleset.sh --ruleset="aws" --version="v0.28.0" --sha="f099944a68e464665dfd4067b663ff6c68229b703d5dff9faea789d3974b078a"

# Install terraform (tfenv)
COPY install-scripts/tfenv.sh /usr/src/install-scripts/tfenv.sh
RUN /usr/src/install-scripts/tfenv.sh --latest-terraform-version="1.5.7" --tfenv-version="v3.0.0"

# Install tfsec
COPY install-scripts/tfsec.sh /usr/src/install-scripts/tfsec.sh
RUN /usr/src/install-scripts/tfsec.sh --version="v1.28.4" --sha="ecd6c40122835356b78d9488328c2c1c186a1a03c8ae00c6a33ba1c08909a4d9"

# Install Open Policy Agent
COPY install-scripts/opa.sh /usr/src/install-scripts/opa.sh
RUN /usr/src/install-scripts/opa.sh --version="v0.43.0" --sha="d5337139a7ccb04149bd9f96ab7a1641a3e3c39f6e1fffa610c7a5c054b0881f"

# Install sops
COPY install-scripts/sops.sh /usr/src/install-scripts/sops.sh
RUN /usr/src/install-scripts/sops.sh --version="v3.8.1" --sha="d6bf07fb61972127c9e0d622523124c2d81caf9f7971fb123228961021811697"

# Install GitHub CLI
COPY install-scripts/github-cli.sh /usr/src/install-scripts/github-cli.sh
RUN /usr/src/install-scripts/github-cli.sh --version="2.39.2" --sha="460d270662f80b8314928ff80ac36db619b8bbac281a1331afd0d6e822f40426"

# Install jq
COPY install-scripts/jq.sh /usr/src/install-scripts/jq.sh
RUN /usr/src/install-scripts/jq.sh --version="1.7" --sha="2f312b9587b1c1eddf3a53f9a0b7d276b9b7b94576c85bda22808ca950569716"

# Install yq
COPY install-scripts/yq.sh /usr/src/install-scripts/yq.sh
RUN /usr/src/install-scripts/yq.sh --version="3.1.0-3"

# Install kubectl
COPY install-scripts/kubectl.sh /usr/src/install-scripts/kubectl.sh
RUN /usr/src/install-scripts/kubectl.sh --version="v1.28.4" --sha="893c92053adea6edbbd4e959c871f5c21edce416988f968bec565d115383f7b8"

# Install helm
COPY install-scripts/helm.sh /usr/src/install-scripts/helm.sh
RUN /usr/src/install-scripts/helm.sh --version="v3.13.2" --sha="55a8e6dce87a1e52c61e0ce7a89bf85b38725ba3e8deb51d4a08ade8a2c70b2d"

COPY --from=tf-prepare-builder /workspace/tf-prepare /usr/local/bin/tf-prepare
RUN chmod +x /usr/local/bin/tf-prepare

# Install packer
COPY install-scripts/packer.sh /usr/src/install-scripts/packer.sh
RUN /usr/src/install-scripts/packer.sh --version="1.9.4" --sha="6cd5269c4245aa8c99e551d1b862460d63fe711c58bec618fade25f8492e80d9"

#Cleanup
RUN apt-get autoremove && \
    apt-get clean

RUN rm -rf /tmp/install

COPY opa-policies /opt/opa-policies
COPY terraform.sh /opt/terraform.sh
COPY packer.sh /opt/packer.sh

ENV HOME=/work

WORKDIR /work
