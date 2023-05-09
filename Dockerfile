FROM public.ecr.aws/amazonlinux/amazonlinux:latest AS hugo_build_env

ARG NODEJS_LTS_DOWNLOAD_URL=https://nodejs.org/dist/v16.16.0/node-v16.16.0-linux-x64.tar.xz
ARG NODEJS_LTS_DOWNLOAD_SHA256=edcb6e9bb049ae365611aa209fc03c4bfc7e0295dbcc5b2f1e710ac70384a8ec

ARG HUGO_DOWNLOAD_URL=https://github.com/gohugoio/hugo/releases/download/v0.101.0/hugo_extended_0.101.0_Linux-64bit.tar.gz
ARG HUGO_DOWNLOAD_SHA256=8c3adf2ace1604468325a6dd094bcc41c141c4a28a0c1ebbeb0022e714897595

ARG GO_VERSION=go1.19

RUN yum update -y && yum install -y tar xz gzip git make golang

RUN curl -L -o nodejs.tar.xz ${NODEJS_LTS_DOWNLOAD_URL} && \
    echo "${NODEJS_LTS_DOWNLOAD_SHA256} nodejs.tar.xz" | sha256sum -c -

RUN mkdir -p /opt/nodejs && \
    tar --strip-components=1 -xJf nodejs.tar.xz -C /opt/nodejs && rm -f nodejs.tar.xz

RUN go install golang.org/dl/${GO_VERSION}@latest && \
    $HOME/go/bin/${GO_VERSION} download

RUN curl -L -o hugo.tar.gz ${HUGO_DOWNLOAD_URL} && \
    echo "${HUGO_DOWNLOAD_SHA256} hugo.tar.gz" | sha256sum -c -

RUN tar -xzf hugo.tar.gz hugo && \
    mv hugo /usr/local/bin && rm -f hugo.tar.gz

FROM hugo_build_env

ARG GO_VERSION=go1.19
ARG SITE_ENV=production
ENV SITE_ENV=${SITE_ENV}
ENV PATH="${HOME}/sdk/${GO_VERSION}:${PATH}"

ADD . /aws-sdk-go-v2-docs

WORKDIR /aws-sdk-go-v2-docs

ENV PATH /opt/nodejs/bin:${PATH}

CMD ["make", "setup", "generate"]
