FROM alpine:3.12.1

RUN wget https://github.com/terraform-linters/tflint/releases/download/v0.21.0/tflint_linux_amd64.zip && \
    unzip tflint_linux_amd64.zip && \
    rm tflint_linux_amd64.zip && \
    mv tflint /usr/bin/tflint

RUN wget https://releases.hashicorp.com/terraform/0.13.5/terraform_0.13.5_linux_amd64.zip && \
    unzip terraform_0.13.5_linux_amd64.zip && \
    rm terraform_0.13.5_linux_amd64.zip && \
    mv terraform /usr/bin/terraform

COPY terraform.sh .
