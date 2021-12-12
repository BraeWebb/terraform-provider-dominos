FROM golang

RUN go install github.com/hashicorp/terraform@v0.12.19

RUN mkdir -p /workdir/tf-dominos
COPY . /workdir/tf-dominos
RUN cd /workdir/tf-dominos && \
    go mod init terraform-provider-dominos && \
    go get github.com/hashicorp/terraform@v0.12.19 && \
    make

# RUN mkdir -p /workdir/tf-dominos && \
#     git clone --depth 1 https://github.com/ndmckinley/terraform-provider-dominos.git && \
#     cd terraform-provider-dominos && \
#     go mod init terraform-provider-dominos && \
#     go get github.com/hashicorp/terraform@v0.12.19

# RUN cd terraform-provider-dominos && \
#     make

RUN mkdir -p ~/.terraform.d/plugins && \
    mv /workdir/tf-dominos/terraform-provider-dominos ~/.terraform.d/plugins && \
    chmod +x ~/.terraform.d/plugins/terraform-provider-dominos

ENV TF_LOG=trace
ENTRYPOINT [ "terraform" ]

# RUN mkdir -p /home/.terraform.d/plugins && \
#     wget https://github.com/ndmckinley/terraform-provider-dominos/raw/master/bin/terraform-provider-dominos -O /home/.terraform.d/plugins/terraform-provider-dominos && \
#     chmod +x /home/.terraform.d/plugins/terraform-provider-dominos
