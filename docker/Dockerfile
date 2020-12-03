FROM alpine:3.12.1

ENV USER="tools"
ENV GROUP="tools"

RUN mkdir -p /tmp/install
WORKDIR /tmp/install

RUN addgroup -g 1000 -S ${GROUP} && \
    adduser -S --ingroup ${GROUP} --uid 1000 ${USER}

RUN mkdir -p /usr/src

RUN apk update && \
    apk add bash

# Base should always be first
COPY install-scripts/base.sh /usr/src/install-scripts/base.sh
RUN /usr/src/install-scripts/base.sh

# Install Azure CLI
COPY install-scripts/azure-cli.sh /usr/src/install-scripts/azure-cli.sh
RUN /usr/src/install-scripts/azure-cli.sh --version="2.15.1"

# Install Ansible
COPY install-scripts/ansible.sh /usr/src/install-scripts/ansible.sh
RUN /usr/src/install-scripts/ansible.sh --version="2.10.3"

# Install Packer
COPY install-scripts/packer.sh /usr/src/install-scripts/packer.sh
RUN /usr/src/install-scripts/packer.sh --version="1.6.5" --sha="a49f6408a50c220fe3f1a6192ea21134e2e8f31092c507614cd27ad4f913234b"

# Install tflint
COPY install-scripts/tflint.sh /usr/src/install-scripts/tflint.sh
RUN /usr/src/install-scripts/tflint.sh --version="v0.21.0" --sha="f80a85dbe88d39231faef5ad10f3a2ecccca70496476ebbc608a78bfabcaa904"

# Install terraform (tfenv)
COPY install-scripts/tfenv.sh /usr/src/install-scripts/tfenv.sh
RUN /usr/src/install-scripts/tfenv.sh --latest-terraform-version="0.13.5" --tfenv-version="v2.0.0" --user="${USER}" --group="${GROUP}"

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

# Cleanup should always be last
COPY install-scripts/cleanup.sh /usr/src/install-scripts/cleanup.sh
RUN /usr/src/install-scripts/cleanup.sh

RUN rm -rf /tmp/install

COPY opa-policies /opt/opa-policies
COPY terraform.sh /opt/terraform.sh
COPY packer.sh /opt/packer.sh

USER ${USER}

WORKDIR /home/${USER}
