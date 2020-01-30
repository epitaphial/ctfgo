#Version: 1.0

FROM golang:1.13.4

MAINTAINER Curled <a@stega.cn>

#创建工作目录
RUN mkdir -p /ctfgo
#进入工作目录
WORKDIR /ctfgo

#将当前目录下的所有文件复制到指定位置
COPY . /ctfgo

#下载bee工具
RUN export GO111MODULE=off && go get -u github.com/beego/bee

#下载beego模块和sqlite模块
RUN go env -w GOPROXY=https://goproxy.io,direct && go get -u github.com/astaxie/beego && go get -u github.com/mattn/go-sqlite3

RUN chmod 755 entrypoint.sh

#端口
EXPOSE 8080

#运行
ENTRYPOINT ["./entrypoint.sh"]