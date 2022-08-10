### alpine
```shell
docker pull alpine:3.16.0
docker pull alpine:3.9.6
docker pull alpine:3.8.5
```

### alpine-utils
#### alpine-utils:3.16.0
`docker build . -t alpine-utils:3.16.0`

```dockerfile
FROM alpine:3.16.0
RUN ln -s /lib /lib64                                         \
 && cp -f /etc/apk/repositories /etc/apk/origin_repositories  \
 && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g'   \
     /etc/apk/repositories                                    \
 && apk add --no-cache                                        \
     musl libgcc libstdc++ libc6-compat libffi zlib libxml2   \
     openssh zsh curl linux-pam iproute2 tzdata               \
 && cp -f /etc/apk/origin_repositories /etc/apk/repositories  \
 && rm -rf /var/cache/apk/*
```
#### alpine-utils:3.9.6
`docker build . -t alpine-utils:3.9.6`

```dockerfile
FROM alpine:3.9.6
RUN ln -s /lib /lib64                                         \
 && cp -f /etc/apk/repositories /etc/apk/origin_repositories  \
 && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g'   \
     /etc/apk/repositories                                    \
 && apk add --no-cache                                        \
     musl libgcc libstdc++ libc6-compat libffi zlib libxml2   \
     openssh zsh curl linux-pam iproute2 tzdata               \
 && cp -f /etc/apk/origin_repositories /etc/apk/repositories  \
 && rm -rf /var/cache/apk/*
```
#### alpine-utils:3.8.5
`docker build . -t alpine-utils:3.8.5`

```dockerfile
FROM alpine:3.8.5
RUN ln -s /lib /lib64                                         \
 && cp -f /etc/apk/repositories /etc/apk/origin_repositories  \
 && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g'   \
     /etc/apk/repositories                                    \
 && apk add --no-cache                                        \
     musl libgcc libstdc++ libc6-compat libffi zlib libxml2   \
     openssh zsh curl linux-pam iproute2 tzdata               \
 && cp -f /etc/apk/origin_repositories /etc/apk/repositories  \
 && rm -rf /var/cache/apk/*
```

### alpine-utils-build
#### alpine-utils-build:v3.16-k5.16
`docker build . -t alpine-utils-build:v3.16-k5.16`

```dockerfile
FROM alpine-utils:3.16.0
RUN ln -s /lib /lib64                                         \
 && cp -f /etc/apk/repositories /etc/apk/origin_repositories  \
 && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g'   \
     /etc/apk/repositories                                    \
 && apk add --no-cache                                        \
     llvm13-libs llvm13-dev llvm13-static llvm13              \
     llvm13-test-utils linux-headers                          \
     clang gcc g++ cmake make binutils git                    \
 && cp -f /etc/apk/origin_repositories /etc/apk/repositories  \
 && rm -rf /var/cache/apk/*
```
#### alpine-utils-build:v3.9-k4.18
`docker build . -t alpine-utils-build:v3.9-k4.18`

```dockerfile
FROM alpine-utils:3.9.6
RUN ln -s /lib /lib64                                         \
 && cp -f /etc/apk/repositories /etc/apk/origin_repositories  \
 && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g'   \
     /etc/apk/repositories                                    \
 && apk add --no-cache                                        \
     llvm5-libs llvm5-dev llvm5-static llvm5-test-utils llvm5 \
     clang gcc g++ linux-headers cmake make binutils git      \
 && cp -f /etc/apk/origin_repositories /etc/apk/repositories  \
 && rm -rf /var/cache/apk/*
```
#### alpine-utils-build:v3.8-k4.4
`docker build . -t alpine-utils-build:v3.8-k4.4`

```dockerfile
FROM alpine-utils:3.8.5
RUN ln -s /lib /lib64                                         \
 && cp -f /etc/apk/repositories /etc/apk/origin_repositories  \
 && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g'   \
     /etc/apk/repositories                                    \
 && apk add --no-cache                                        \
     llvm5-libs llvm5-dev llvm5-static llvm5-test-utils llvm5 \
     clang gcc g++ linux-headers cmake make binutils git      \
 && cp -f /etc/apk/origin_repositories /etc/apk/repositories  \
 && rm -rf /var/cache/apk/*
```

### golang-develop

```
golang-develop$ ls
authorized_keys	  go1.17.11.linux-amd64.tar.gz	scripts
dockerfiles			  id_rsa				                sshd_config
env.config		    id_rsa.pub			              start-sshd.sh
```

start-sshd.sh

```shell
#!/bin/sh

ssh-keygen -A

if [ -z "${SSHD_PORT}" ]; then
  /usr/sbin/sshd -D
else
  /usr/sbin/sshd -D -p ${SSHD_PORT}
fi

exit 0
```

#### golang-develop:v3.16-k5.16-go17

`docker build . -f ./dockerfiles/df-v3.16-k5.16-go17 -t golang-develop:v3.16-k5.16-go17`

```dockerfile
FROM alpine-utils-build:v3.16-k5.16
VOLUME /data
ENV GO111MODULE=auto                            \
    GOROOT=/opt/golang/go                       \
    GOPATH=/opt/golang/gopath                   \
    GONOSUMDB=*                                 \
#   GOPRIVATE=com.ffcs.*                        \
#   GOPROXY=https://goproxy.cn,direct           \
    GO_BASEDIR=/opt/golang                      \
    GO_TARBALL=go1.17.11.linux-amd64.tar.gz     \
    PATH=${PATH}:${GOROOT}/bin
RUN mkdir /root/.ssh && chmod 700 /root/.ssh                                     \
 && cp -f /etc/apk/repositories /etc/apk/origin_repositories                     \
 && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g'/etc/apk/repositories \
 && apk add --no-cache libpcap-doc libpcap-dev libpcap git protobuf              \
 && cp -f /etc/apk/origin_repositories /etc/apk/repositories
COPY ${GO_TARBALL} /opt/${GO_TARBALL}
RUN mkdir -p ${GOROOT}                              \
    && ln -s /data/gopath ${GOPATH}                 \
    && tar -xzf /opt/${GO_TARBALL} -C ${GO_BASEDIR} \
    && rm -rf /var/cache/apk/* /opt/${GO_TARBALL}
ADD authorized_keys id_rsa id_rsa.pub /root/.ssh/
ADD sshd_config /etc/ssh/sshd_config
ADD start-sshd.sh /start-sshd.sh
RUN chmod 700 /start-sshd.sh && chmod 600 /root/.ssh/* && chmod 644 /root/.ssh/*.pub
EXPOSE 1622
CMD ["/start-sshd.sh"]
```

#### golang-develop:v3.9-k4.18-go17

`docker build . -f ./dockerfiles/df-v3.9-k4.18-go17 -t golang-develop:v3.9-k4.18-go17`

```dockerfile
FROM alpine-utils-build:v3.9-k4.18
VOLUME /data
ENV GO111MODULE=auto                            \
    GOROOT=/opt/golang/go                       \
    GOPATH=/opt/golang/gopath                   \
    GONOSUMDB=*                                 \
#   GOPRIVATE=com.ffcs.*                        \
#   GOPROXY=https://goproxy.cn,direct           \
    GO_BASEDIR=/opt/golang                      \
    GO_TARBALL=go1.17.11.linux-amd64.tar.gz     \
    PATH=${PATH}:${GOROOT}/bin
RUN mkdir /root/.ssh && chmod 700 /root/.ssh                                     \
 && cp -f /etc/apk/repositories /etc/apk/origin_repositories                     \
 && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g'/etc/apk/repositories \
 && apk add --no-cache libpcap-doc libpcap-dev libpcap git protobuf              \
 && cp -f /etc/apk/origin_repositories /etc/apk/repositories
COPY ${GO_TARBALL} /opt/${GO_TARBALL}
RUN mkdir -p ${GOROOT}                              \
    && ln -s /data/gopath ${GOPATH}                 \
    && tar -xzf /opt/${GO_TARBALL} -C ${GO_BASEDIR} \
    && rm -rf /var/cache/apk/* /opt/${GO_TARBALL}
ADD authorized_keys id_rsa id_rsa.pub /root/.ssh/
ADD sshd_config /etc/ssh/sshd_config
ADD start-sshd.sh /start-sshd.sh
RUN chmod 700 /start-sshd.sh && chmod 600 /root/.ssh/* && chmod 644 /root/.ssh/*.pub
EXPOSE 3922
CMD ["/start-sshd.sh"]
```

#### golang-develop:v3.8-k4.4-go17

`docker build . -f ./dockerfiles/df-v3.8-k4.4-go17 -t golang-develop:v3.8-k4.4-go17`

```dockerfile
FROM alpine-utils-build:v3.8-k4.4
VOLUME /data
ENV GO111MODULE=auto                            \
    GOROOT=/opt/golang/go                       \
    GOPATH=/opt/golang/gopath                   \
    GONOSUMDB=*                                 \
#   GOPRIVATE=com.ffcs.*                        \
#   GOPROXY=https://goproxy.cn,direct           \
    GO_BASEDIR=/opt/golang                      \
    GO_TARBALL=go1.17.11.linux-amd64.tar.gz     \
    PATH=${PATH}:${GOROOT}/bin
RUN mkdir /root/.ssh && chmod 700 /root/.ssh                                     \
 && cp -f /etc/apk/repositories /etc/apk/origin_repositories                     \
 && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g'/etc/apk/repositories \
 && apk add --no-cache libpcap-doc libpcap-dev libpcap git protobuf              \
 && cp -f /etc/apk/origin_repositories /etc/apk/repositories
COPY ${GO_TARBALL} /opt/${GO_TARBALL}
RUN mkdir -p ${GOROOT}                              \
    && ln -s /data/gopath ${GOPATH}                 \
    && tar -xzf /opt/${GO_TARBALL} -C ${GO_BASEDIR} \
    && rm -rf /var/cache/apk/* /opt/${GO_TARBALL}
ADD authorized_keys id_rsa id_rsa.pub /root/.ssh/
ADD sshd_config /etc/ssh/sshd_config
ADD start-sshd.sh /start-sshd.sh
RUN chmod 700 /start-sshd.sh && chmod 600 /root/.ssh/* && chmod 644 /root/.ssh/*.pub
EXPOSE 3822
CMD ["/start-sshd.sh"]
```

### start

```shell
docker create --privileged                      \
  --name 3.9-jeevan-go                          \
  --volume /data/volumes/jeevan-go:/data        \
  --network host                                \
  --publish 3922:3922                           \
  --env SSHD_PORT=3922                          \
  --env GO111MODULE=auto                        \
  --env GOROOT=/opt/golang/go                   \
  --env GOPATH=/opt/golang/gopath               \
  --env GONOSUMDB="*"                           \
  --env GO_BASEDIR=/opt/golang                  \
  --env https_proxy=http://192.168.109.33:7890  \
  --env http_proxy=http://192.168.109.33:7890   \
  --env all_proxy=socks5://192.168.109.33:7890  \
  --env PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/golang/go/bin:/data/gopath/bin" \
jeevan86/golang-develop:v3.9-k4.18-go17

docker start 3.9-jeevan-go

docker exec -it 3.9-jeevan-go zsh
```
