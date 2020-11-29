FROM alpine:3.12.1

RUN apk add git bash curl

RUN apk update && \
    apk add py-pip && \
    apk add --virtual=build gcc libffi-dev musl-dev openssl-dev python3-dev make && \
    pip --no-cache-dir install -U pip && \
    pip --no-cache-dir install azure-cli && \
    apk del --purge build

RUN pip --no-cache-dir install ansible

RUN wget https://releases.hashicorp.com/packer/1.6.5/packer_1.6.5_linux_amd64.zip && \
    unzip packer_1.6.5_linux_amd64.zip && \
    rm packer_1.6.5_linux_amd64.zip && \
    mv packer /usr/local/bin/packer

RUN wget https://github.com/terraform-linters/tflint/releases/download/v0.21.0/tflint_linux_amd64.zip && \
    unzip tflint_linux_amd64.zip && \
    rm tflint_linux_amd64.zip && \
    mv tflint /usr/local/bin/tflint

RUN git clone https://github.com/tfutils/tfenv.git ~/.tfenv && \
    ln -s ~/.tfenv/bin/* /usr/local/bin && \
    tfenv install 0.13.3 && \
    tfenv install 0.13.4 && \
    tfenv install 0.13.5 && \
    tfenv use 0.13.5

RUN wget https://github.com/open-policy-agent/opa/releases/download/v0.24.0/opa_linux_amd64 && \
    chmod +x opa_linux_amd64 && \
    mv opa_linux_amd64 /usr/local/bin/opa

COPY opa-policies /opt/opa-policies
COPY terraform.sh /opt/terraform.sh
