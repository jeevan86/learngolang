```zsh
#!/usr/bin/env sh

WORK_DIR=$(dirname "$0") && cd "${WORK_DIR}" && WORK_DIR=$(pwd)
SRC_DIR=${WORK_DIR}/src
BUILD_DIR=${WORK_DIR}/build
DEST_DIR=${WORK_DIR}/dest

export WORK_DIR SRC_DIR BUILD_DIR DEST_DIR

mkdir -p ${BUILD_DIR}
mkdir -p ${DEST_DIR}

cd ${SRC_DIR} && BUILD_STATIC_ONLY=y OBJDIR=${BUILD_DIR} DESTDIR=${DEST_DIR} make install

exit 0
```

```zsh
export CILIUM_EBPF=${GOPATH}/src/github.com/cilium/ebpf
cd ${CILIUM_EBPF}
CC=clang CGO_ENABLED=1                                           \
go build -gcflags "all=-N -l"                                    \
 -a -ldflags '-linkmode external -extldflags "-fno-PIC -static"' \
 -o ${GOPATH}/bin/bpf2go ${CILIUM_EBPF}/cmd/bpf2go
 
export CILIUM_EBPF=${GOPATH}/src/github.com/cilium/ebpf
cd ${CILIUM_EBPF}
CC=clang CGO_ENABLED=1                                           \
go build -gcflags "all=-N -l"                                    \
 -a -ldflags '-linkmode external -extldflags "-fno-PIC -static"' \
 -o ${GOPATH}/bin/bpfo2go ${CILIUM_EBPF}/cmd/bpfo2go
```

```zsh
#!/usr/bin/env sh

if [ -z "${GOPATH}" ]; then
  echo "Exiting, env GOPATH is need."
  exit
fi

CURR_DIR=$(dirname "$0") && cd "${CURR_DIR}" && CURR_DIR=$(pwd)
export PRJ_EBPF_MODULE=${CURR_DIR##*/}
export GOPACKAGE=${PRJ_EBPF_MODULE}
echo "Generating ${PRJ_EBPF_MODULE}"

export LIBBPF_REPO=${GOPATH}/src/github.com/libbpf/libbpf
export LIBBPF_DEST=${LIBBPF_REPO}/dest
export LIBBPF_HEADERS=${LIBBPF_DEST}/usr/include
export LIBBPF_OBJECTS=${LIBBPF_DEST}/usr/lib64
export CILIUM_EBPF=${GOPATH}/src/github.com/cilium/ebpf
export PRJ_ROOT=${GOPATH}/src/github.com/jeevan86/learngolang
export PRJ_EBPF_MODULE_PATH=${PRJ_ROOT}/pkg/ebpf/cilium/${PRJ_EBPF_MODULE}

cd ${PRJ_EBPF_MODULE_PATH} || exit

- rm -f ./bpf_bpfe*.go ./bpf_bpfe*.o ./bpf_bpfe*.o.d

# 1 - little endian - 如x86_64
clang -O2 -mcpu=v1 -g -fno-ident -MD -MP -MF ./bpf_bpfel.o.d     \
 -target bpfel                                                   \
 -fdebug-prefix-map=${PRJ_EBPF_MODULE_PATH}=.                    \
 -fdebug-compilation-dir .                                       \
 -Wunused-command-line-argument                                  \
 -D__BPF_TARGET_MISSING="\"Please provide -target bpfeb|bpfel\"" \
 -I${CILIUM_EBPF}/examples/headers                               \
 -I${LIBBPF_HEADERS}                                             \
 -L${LIBBPF_OBJECTS}                                             \
 -c ${PRJ_EBPF_MODULE_PATH}/c/kprobe.c                           \
 -o ${PRJ_EBPF_MODULE_PATH}/bpf_bpfel.o

${GOPATH}/bin/bpfo2go -target bpfel bpf

# 2 - big endian - 如arm
# -MF/dev/fd/3
#clang -O2 -mcpu=v1 -g -fno-ident -MD -MP -MF ./bpf_bpfeb.o.d     \
# -target bpfeb                                                   \
# -fdebug-prefix-map=${PRJ_EBPF_MODULE_PATH}=.                    \
# -fdebug-compilation-dir .                                       \
# -D__BPF_TARGET_MISSING="\"Please provide -target bpfeb|bpfel\"" \
# -I${CILIUM_EBPF}/examples/headers                               \
# -I${LIBBPF_HEADERS}                                             \
# -L${LIBBPF_OBJECTS}                                             \
# -c ${PRJ_EBPF_MODULE_PATH}/c/kprobe.c                           \
# -o ${PRJ_EBPF_MODULE_PATH}/bpf_bpfeb.o

# ${GOPATH}/bin/bpfo2go bpfeb
# $GOPATH/bin/dlv --listen=:8606 --headless=true --api-version=2 --accept-multiclient exec $GOPATH/bin/bpfo2go -- -target bpfel bpf

exit 0
```