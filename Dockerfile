FROM registry.cn-shenzhen.aliyuncs.com/dev-ops/golang:1.20.7-alpine3.18 as golang
ENV APP go-gin-rest-api   
ADD ./ /data/${APP}/${APP}
ADD .git/ /data/${APP}/${APP}/.git
WORKDIR /data/${APP}/${APP} 
RUN GitBranch=$(git name-rev --name-only HEAD) && \
    GitRevision=$(git rev-parse --short HEAD)  && \ 
    GitCommitLog=`git log --pretty=oneline -n 1`  && \ 
    GitCommitLog=${GitCommitLog//\'/\"}  && \ 
    BuildTime=`date +'%Y.%m.%d.%H%M%S'`  && \ 
    BuildGoVersion=`go version`  && \ 
    LDFlags="-s -w -X 'main.GitBranch=${GitBranch}' -X 'main.GitRevision=${GitRevision}' -X 'main.GitCommitLog=${GitCommitLog}' -X 'main.BuildTime=${BuildTime}' -X 'main.BuildGoVersion=${BuildGoVersion}'"  && \
    go build -tags=jsoniter -ldflags="$LDFlags" -o  ./${APP} 

FROM registry.cn-shenzhen.aliyuncs.com/dev-ops/alpine:3.18.2
LABEL MAINTAINER="13579443@qq.com"
ENV APP go-gin-rest-api
ENV TZ='Asia/Shanghai' 
RUN TERM=linux && export TERM
WORKDIR /data/${APP}/
COPY --from=golang /data/${APP}/${APP}/${APP} /data/${APP}/${APP}/${APP} 
COPY --from=golang /data/${APP}/${APP}/conf /data/${APP}/${APP}/conf   
CMD ["./${APP}"]