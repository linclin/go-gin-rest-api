# syntax=registry.cn-shenzhen.aliyuncs.com/dev-ops/dockerfile:1.9.0
FROM registry.cn-shenzhen.aliyuncs.com/dev-ops/golang:1.23.1-alpine3.20 as golang
ENV APP go-gin-rest-api  
RUN sed -i 's/https/http/' /etc/apk/repositories &&\
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories &&\
    apk update &&\
    apk add ca-certificates git bash openssh gcc musl-dev &&\
    rm -rf /var/cache/apk/*   /tmp/*  
ADD ./ /app/${APP}
ADD .git/ /app/${APP}/.git
WORKDIR /app/${APP}
RUN export GitBranch=$(git name-rev --name-only HEAD) &&\
    export GitRevision=$(git rev-parse --short HEAD) &&\
    export GitCommitLog=`git log --pretty=oneline -n 1` &&\
    export BuildTime=`date +'%Y.%m.%d.%H%M%S'` &&\
    export BuildGoVersion=`go version` &&\
    export LDFlags="-s -w -X 'main.GitBranch=${GitBranch}' -X 'main.GitRevision=${GitRevision}' -X 'main.GitCommitLog=${GitCommitLog}' -X 'main.BuildTime=${BuildTime}' -X 'main.BuildGoVersion=${BuildGoVersion}'"  &&\
    go build -tags=jsoniter -ldflags="$LDFlags" -o  ./${APP}

FROM registry.cn-shenzhen.aliyuncs.com/dev-ops/alpine:3.20.3
LABEL MAINTAINER="13579443@qq.com"
ENV APP go-gin-rest-api
ENV TZ='Asia/Shanghai' 
ENV LANG UTF-8
ENV LC_ALL zh_CN.UTF-8
ENV LC_CTYPE zh_CN.UTF-8
RUN TERM=linux && export TERM
RUN sed -i 's/https/http/' /etc/apk/repositories &&\
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories &&\
    apk update &&\
    apk add ca-certificates tzdata  bash  sudo busybox-extras curl  &&\
    echo "Asia/Shanghai" > /etc/timezone &&\ 
    rm -rf /var/cache/apk/*   /tmp/* &&\
    apk del tzdata &&\
    addgroup -g 1000 app &&\
    adduser -u 1000 -G app -D app &&\
    adduser -u 1001 -G app -D dev &&\
    chmod 4755 /bin/busybox  &&\
    sudo chmod +w /etc/sudoers &&\
    echo "app ALL=(ALL) NOPASSWD: NOPASSWD: /bin/su " >> /etc/sudoers &&\
    mkdir -p /app  &&\
    chown -R app:app /app 
WORKDIR /app/${APP}/
COPY --from=golang --chown=app:app --chmod=755 /app/${APP}/${APP} /app/${APP}/${APP} 
COPY --from=golang --chown=app:app --chmod=755 /app/${APP}/conf /app/${APP}/conf  
CMD ["./${APP}"]