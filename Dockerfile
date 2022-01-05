FROM public.ecr.aws/amazonlinux/amazonlinux:latest AS hugo_build_env

ARG NODEJS_LTS_DOWNLOAD_URL=https://nodejs.org/dist/v16.13.1/node-v16.13.1-linux-x64.tar.xz
ARG NODEJS_LTS_DOWNLOAD_SHA256=a3721f87cecc0b52b0be8587c20776ac7305db413751db02c55aa2bffac15198

ARG HUGO_DOWNLOAD_URL=https://github.com/gohugoio/hugo/releases/download/v0.91.2/hugo_extended_0.91.2_Linux-64bit.tar.gz
ARG HUGO_DOWNLOAD_SHA256=e9e2b35ebef6ed41581eb18909b8ee02ee9285d209f7d9ecc5caf5207b7dc8e5

RUN yum update -y && yum install -y tar xz gzip git coreutils make
RUN curl -L -o nodejs.tar.xz ${NODEJS_LTS_DOWNLOAD_URL} && \
    echo "${NODEJS_LTS_DOWNLOAD_SHA256} nodejs.tar.xz" | sha256sum -c -

RUN mkdir -p /opt/nodejs && \
    tar --strip-components=1 -xJf nodejs.tar.xz -C /opt/nodejs && rm -f nodejs.tar.xz

RUN curl -L -o hugo.tar.gz ${HUGO_DOWNLOAD_URL} && \
    echo "${HUGO_DOWNLOAD_SHA256} hugo.tar.gz" | sha256sum -c -

RUN tar -xzf hugo.tar.gz hugo && \
    mv hugo /usr/local/bin && rm -f hugo.tar.gz

FROM hugo_build_env

ARG SITE_ENV=production

ADD . /aws-sdk-go-v2-docs

WORKDIR /aws-sdk-go-v2-docs

ENV PATH /opt/nodejs/bin:${PATH}

RUN make setup generate SITE_ENV=${SITE_ENV}
