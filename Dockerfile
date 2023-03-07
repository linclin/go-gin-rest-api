FROM golang:1.20.1-alpine3.17 as golang
ENV APP go-gin-rest-api
RUN sed -i 's/https/http/' /etc/apk/repositories && \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update --no-cache && \
    apk add --no-cache ca-certificates git bash gcc openssh && \ 
    rm -rf /var/cache/apk/*   /tmp/*     
ADD ./ /data/${APP}/${APP}
ADD .git/ /data/${APP}/${APP}/.git
WORKDIR /data/${APP}/${APP} 
RUN GitBranch=$(git name-rev --name-only HEAD) && \
    GitRevision=$(git rev-parse --short HEAD)  && \ 
    GitCommitLog=`git log --pretty=oneline -n 1`  && \ 
    GitCommitLog=${GitCommitLog//\'/\"}  && \ 
    BuildTime=`date +'%Y.%m.%d.%H%M%S'`  && \ 
    BuildGoVersion=`go version`  && \ 
    LDFlags="-w -X 'main.GitBranch=${GitBranch}' -X 'main.GitRevision=${GitRevision}' -X 'main.GitCommitLog=${GitCommitLog}' -X 'main.BuildTime=${BuildTime}' -X 'main.BuildGoVersion=${BuildGoVersion}'"  && \
    go build -ldflags="$LDFlags" -o  ./${APP} 

FROM alpine:3.17.2
LABEL MAINTAINER="13579443@qq.com"
ENV APP go-gin-rest-api
ENV TZ='Asia/Shanghai' 
RUN TERM=linux && export TERM
RUN sed -i 's/https/http/' /etc/apk/repositories && \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update --no-cache && \
    apk add  --no-cache  ca-certificates bash tzdata sudo busybox-extras curl net-tools && \ 
    echo "Asia/Shanghai" > /etc/timezone && \
    cp -r -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    rm -rf /var/cache/apk/*   /tmp/*   
WORKDIR /data/${APP}/
COPY --from=golang /data/${APP}/${APP}/${APP} /data/${APP}/${APP}/${APP} 
COPY --from=golang /data/${APP}/${APP}/conf /data/${APP}/${APP}/conf   
CMD ["./${APP}"]