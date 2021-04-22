FROM public.ecr.aws/amazonlinux/amazonlinux:latest AS hugo_build_env

ARG NODEJS_LTS_DOWNLOAD_URL=https://nodejs.org/dist/v14.16.1/node-v14.16.1-linux-x64.tar.xz
ARG NODEJS_LTS_DOWNLOAD_SHA256=85a89d2f68855282c87851c882d4c4bbea4cd7f888f603722f0240a6e53d89df

ARG HUGO_DOWNLOAD_URL=https://github.com/gohugoio/hugo/releases/download/v0.82.0/hugo_extended_0.82.0_Linux-64bit.tar.gz
ARG HUGO_DOWNLOAD_SHA256=171b8f935acc60f74e1eb9edb73fc5e9afaa3affaed4ddafd072ada800ce8748

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
