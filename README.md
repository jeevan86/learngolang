# pcap

## gopacket

```zsh
brew install FiloSottile/musl-cross/musl-cross
brew install libpcap

CGO_ENABLED=1                                           \
CGO_CFLAGS="-I/usr/local/Cellar/libpcap/1.10.1/include" \
CGO_LDFLAGS="-L/usr/local/Cellar/libpcap/1.10.1/lib"    \
CC=x86_64-linux-musl-gcc                                \
CXX=x86_64-linux-musl-g++                               \
GOOS=linux GOARCH=amd64                                 \
go build --ldflags "-extldflags -static" -o dist/gopcap-agent          cmd/pcap/main.go
go build --ldflags "-extldflags -static" -o dist/gopcap-collector-grpc cmd/collector/grpc.go

go build -gcflags="-d pctab=pctoinline" ...

go build -gcflags="-d pctab=pctoinline" -o dist/gopcap-agent          cmd/pcap/main.go
go build -gcflags="-d pctab=pctoinline" -o dist/gopcap-collector-grpc cmd/collector/grpc.go


go generate ebpf/cilium/kprobe/


export CILIUM_EBPF=${GOPATH}/src/github.com/cilium/ebpf
go run ${CILIUM_EBPF}/cmd/bpf2go/main.go         \
 -cc $BPF_CLANG -cflags $BPF_CFLAGS bpf kprobe.c \
-- -I${CILIUM_EBPF}/examples/headers

go build -o ${GOPATH}/bin/ ebpf/cilium/kprobe/
```

# ebpf

## cilium-ebpf

### bpf2go

```zsh
export CILIUM_EBPF=${GOPATH}/src/github.com/cilium/ebpf
cd ${CILIUM_EBPF}
CC=clang CGO_ENABLED=0                                           \
go build -gcflags "all=-N -l"                                    \
 -a -ldflags '-linkmode external -extldflags "-fno-PIC -static"' \
 -o ${GOPATH}/bin/bpf2go ${CILIUM_EBPF}/cmd/bpf2go
```

### go generate

```zsh
export CILIUM_EBPF=${GOPATH}/src/github.com/cilium/ebpf
export PRJ_ROOT=${GOPATH}/src/github.com/jeevan86/learngolang
export CILIUM_EBPF_KPROBE=kprobe
export GOPACKAGE=${CILIUM_EBPF_KPROBE}
cd ${PRJ_ROOT}/pkg/ebpf/cilium/${CILIUM_EBPF_KPROBE}
rm -f ./bpf_bpfe*.go ./bpf_bpfe*.o
export BPF_CLANG=clang
export BPF_CFLAGS="-I../headers"
${GOPATH}/bin/bpf2go     \
   -cc ${BPF_CLANG}      \
   -cflags ${BPF_CFLAGS} \
  bpf ./kprobe.c         \
 -- -I${CILIUM_EBPF}/examples/headers -I/usr/include/linux
```

### build app

```zsh
export PRJ_ROOT=${GOPATH}/src/github.com/jeevan86/learngolang
go build -o ${GOPATH}/bin/ebpf_cilium ${PRJ_ROOT}/cmd/ebpf_cilium.go
```

### debug app

#### 首先使用静态编译delve

```zsh
cd ${GOPATH}/src/github.com/go-delve/delve/cmd/dlv
CGO_ENABLE=1 go build -a -ldflags '-extldflags "-static"' -o ${GOPATH}/bin/dlv ./main.go
```

#### 然后加参数并静态编译app

```zsh
export PRJ_ROOT=${GOPATH}/src/github.com/jeevan86/learngolang
CC=/usr/bin/x86_64-alpine-linux-musl-gcc CGO_ENABLE=1            \
go build -gcflags "all=-N -l"                                    \
 -a -ldflags '-linkmode external -extldflags "-fno-PIC -static"' \
-o ${GOPATH}/bin/ebpf_cilium ${PRJ_ROOT}/cmd/ebpf/cilium/main.go

##### or #####

export PRJ_ROOT=${GOPATH}/src/github.com/jeevan86/learngolang
CC=clang CGO_ENABLE=1                                            \
go build -gcflags "all=-N -l"                                    \
 -a -ldflags '-linkmode external -extldflags "-fno-PIC -static"' \
-o ${GOPATH}/bin/ebpf_cilium ${PRJ_ROOT}/cmd/ebpf/cilium/main.go
```

#### 最后在主机上通过delve启动app

```zsh
./dlv --listen=:8606 --headless=true --api-version=2 --accept-multiclient exec ./ebpf_cilium
```

# grpc
```zsh
# https://github.com/protocolbuffers/protobuf/releases/download/v21.1/protoc-21.1-linux-x86_64.zip
# https://github.com/protocolbuffers/protobuf/releases/download/v21.1/protoc-21.1-osx-x86_64.zip
```
### MacOS版本问题

最新的protoc的二进制是在Mac OS X 11.3之上编译的，以下的版本无法执行

```zsh
$GOPATH/bin/protoc ./helloworld.proto --go_out=plugins=grpc:./
dyld: lazy symbol binding failed: Symbol not found: ___darwin_check_fd_set_overflow
  Referenced from: /Users/huangjian/Documents/workspace/golang/gopath/bin/protoc (which was built for Mac OS X 11.3)
  Expected in: /usr/lib/libSystem.B.dylib

dyld: Symbol not found: ___darwin_check_fd_set_overflow
  Referenced from: /Users/huangjian/Documents/workspace/golang/gopath/bin/protoc (which was built for Mac OS X 11.3)
  Expected in: /usr/lib/libSystem.B.dylib

Abort trap: 6
brew install protobuf
```
自己编译一个

```zsh
git clone git@github.com:protocolbuffers/protobuf.git
cd protobuf && mkdir release && cd release
cmake -Dprotobuf_BUILD_TESTS=OFF ..
make
cp -f protoc $GOPATH/bin/
```
### go的两个插件

```zsh
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
```
### 生成\*.pb.go、\*.grpc.pb.go

自己执行命令或使用go generate

```zsh
# example: https://github.com/grpc/grpc-go/tree/v1.46.2/examples/helloworld
export PATH=$PATH:$GOPATH/bin

cd ./pkg/collect/api/grpc/pb
protoc ./*.proto --go_out=./
protoc ./*.proto --go-grpc_out=./

----------- or ------------
go generate ./pkg/collect/api/grpc/pb
```
cat ./pkg/collect/api/grpc/pb/generate.go
```go
package pb
// protoc protoc-gen-go protoc-gen-go-grpc in env PATH
//go:generate echo running protoc with $PWD $GOARCH $GOOS $GOPACKAGE $GOFILE $GOLINE $DOLLAR
//go:generate protoc $PWD/collect.proto --proto_path=$PWD --go_out=$PWD/
//go:generate protoc $PWD/collect.proto --proto_path=$PWD --go-grpc_out=$PWD/
```
