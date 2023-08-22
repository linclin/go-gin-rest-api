FROM registry.cn-shenzhen.aliyuncs.com/dev-ops/golang:1.21.0-alpine3.18 as golang
ENV APP go-gin-rest-api   
ADD ./ /app/${APP}/
ADD .git/ /app/${APP}/.git
WORKDIR /app/${APP}/
RUN export GitBranch=$(git name-rev --name-only HEAD) && \
    export GitRevision=$(git rev-parse --short HEAD)  && \ 
    export GitCommitLog=`git log --pretty=oneline -n 1`  && \
    export BuildTime=`date +'%Y.%m.%d.%H%M%S'`  && \ 
    export BuildGoVersion=`go version`  && \ 
    export LDFlags="-s -w -X 'main.GitBranch=${GitBranch}' -X 'main.GitRevision=${GitRevision}' -X 'main.GitCommitLog=${GitCommitLog}' -X 'main.BuildTime=${BuildTime}' -X 'main.BuildGoVersion=${BuildGoVersion}'"  && \
    go build -tags=jsoniter -ldflags="$LDFlags" -o  ./${APP} 

FROM registry.cn-shenzhen.aliyuncs.com/dev-ops/alpine:3.18.3
LABEL MAINTAINER="13579443@qq.com"
ENV APP go-gin-rest-api
ENV TZ='Asia/Shanghai' 
RUN TERM=linux && export TERM
WORKDIR /app/${APP}/
COPY --from=golang /app/${APP}/${APP} /app/${APP}/${APP} 
COPY --from=golang /app/${APP}/conf /app/${APP}/conf   
CMD ["./${APP}"]