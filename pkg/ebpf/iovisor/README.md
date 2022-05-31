### 构建命令行

```zsh
export PRJ_ROOT=${GOPATH}/src/github.com/jeevan86/learngolang
CC=clang CGO_ENABLED=0                                           \
go build -gcflags "all=-N -l"                                    \
 -a -ldflags '-linkmode external -extldflags "-fno-PIC -static -static-libstdc++ -static-libgcc"' \
 -o ${GOPATH}/bin/ebpf_iovisor ${PRJ_ROOT}/cmd/ebpf/iovisor/main.go
```