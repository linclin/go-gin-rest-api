FROM registry.cn-shenzhen.aliyuncs.com/dev-ops/golang:1.20.7-alpine3.18 as golang
ENV APP go-gin-rest-api   
ADD ./ /data/${APP}/
ADD .git/ /data/${APP}/.git
WORKDIR /data/${APP}/
RUN export GitBranch=$(git name-rev --name-only HEAD) && \
    export GitRevision=$(git rev-parse --short HEAD)  && \ 
    export GitCommitLog=`git log --pretty=oneline -n 1`  && \
    export BuildTime=`date +'%Y.%m.%d.%H%M%S'`  && \ 
    export BuildGoVersion=`go version`  && \ 
    export LDFlags="-s -w -X 'main.GitBranch=${GitBranch}' -X 'main.GitRevision=${GitRevision}' -X 'main.GitCommitLog=${GitCommitLog}' -X 'main.BuildTime=${BuildTime}' -X 'main.BuildGoVersion=${BuildGoVersion}'"  && \
    go build -tags=jsoniter -ldflags="$LDFlags" -o  ./${APP} 

FROM registry.cn-shenzhen.aliyuncs.com/dev-ops/alpine:3.18.2
LABEL MAINTAINER="13579443@qq.com"
ENV APP go-gin-rest-api
ENV TZ='Asia/Shanghai' 
RUN TERM=linux && export TERM
WORKDIR /data/${APP}/
COPY --from=golang /data/${APP}/${APP} /data/${APP}/${APP} 
COPY --from=golang /data/${APP}/conf /data/${APP}/conf   
CMD ["./${APP}"]