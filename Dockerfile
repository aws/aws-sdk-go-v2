FROM public.ecr.aws/amazonlinux/amazonlinux:latest

ARG SITE_ENV=production

RUN yum update -y && yum install -y tar xz gzip git coreutils make
RUN curl -L -o nodejs.tar.xz https://nodejs.org/dist/v14.15.1/node-v14.15.1-linux-x64.tar.xz && \
    mkdir -p /opt/nodejs && \
    tar --strip-components=1 -xJf nodejs.tar.xz -C /opt/nodejs && rm -f nodejs.tar.xz
RUN curl -L -o hugo.tar.gz https://github.com/gohugoio/hugo/releases/download/v0.79.0/hugo_extended_0.79.0_Linux-64bit.tar.gz && \
    tar -xzf hugo.tar.gz hugo && \
    mv hugo /usr/local/bin && rm -f hugo.tar.gz

ADD . /aws-sdk-go-v2-docs

WORKDIR /aws-sdk-go-v2-docs

ENV PATH /opt/nodejs/bin:${PATH}

RUN make setup generate SITE_ENV=${SITE_ENV}
