#!/usr/bin/env sh

if [ -z "${GOPATH}" ]; then
  echo "Exiting, env GOPATH is need."
  exit
fi

CURR_DIR=$(dirname "$0") && cd "${CURR_DIR}" && CURR_DIR=$(pwd)
export PRJ_EBPF_MODULE=${CURR_DIR##*/}
export GOPACKAGE=${PRJ_EBPF_MODULE}
echo "Generating ${PRJ_EBPF_MODULE}"

export LINUX_SRC=${GOPATH}/src/linux/current
export LINUX_SRC_HEADERS=${LINUX_SRC}/headers
export LINUX_SRC_TOOLS=${LINUX_SRC}/tools

export LIBBPF_REPO=${GOPATH}/src/github.com/libbpf/libbpf
export LIBBPF_DEST=${LIBBPF_REPO}/dest
export LIBBPF_HEADERS=${LIBBPF_DEST}/usr/include
export LIBBPF_OBJECTS=${LIBBPF_DEST}/usr/lib64
export CILIUM_EBPF=${GOPATH}/src/github.com/cilium/ebpf
export PRJ_ROOT=${GOPATH}/src/github.com/jeevan86/learngolang
export PRJ_EBPF_MODULE_PATH=${PRJ_ROOT}/pkg/ebpf/cilium/${PRJ_EBPF_MODULE}

cd "${PRJ_EBPF_MODULE_PATH}" || exit

rm -f ./bpf_bpfe*.go ./bpf_bpfe*.o ./bpf_bpfe*.o.d

# --disable-x86asm                                                \ # disable x86 asm
# --disable-inline-asm                                            \ # disable inline asm
# -mkernel=4.18                                                   \

# 1 - little endian - 如x86_64
clang -O2 -mcpu=v1 -g -fno-ident -MD -MP -MF ./bpf_bpfel.o.d     \
 -target bpfel                                                   \
 -fdebug-prefix-map="${PRJ_EBPF_MODULE_PATH}"=.                  \
 -fdebug-compilation-dir .                                       \
 -Wunused-command-line-argument                                  \
 -no-integrated-as                                               \
 -D__BPF_TARGET_MISSING="\"Please provide -target bpfeb|bpfel\"" \
 -I"${LINUX_SRC_HEADERS}"                                        \
 -I"${LINUX_SRC_TOOLS}"                                          \
 -I"${CILIUM_EBPF}"/examples/headers                             \
 -I"${LIBBPF_HEADERS}"                                           \
 -I"${PRJ_EBPF_MODULE_PATH}"/c                                   \
 -L"${LIBBPF_OBJECTS}"                                           \
 -c "${PRJ_EBPF_MODULE_PATH}"/c/"${PRJ_EBPF_MODULE}".c           \
 -o "${PRJ_EBPF_MODULE_PATH}"/bpf_bpfel.o

#"$CLANG_EXEC -D__KERNEL__ -D__NR_CPUS__=$NR_CPUS "\
#		"-DLINUX_VERSION_CODE=$LINUX_VERSION_CODE "	\
#		"$CLANG_OPTIONS $PERF_BPF_INC_OPTIONS $KERNEL_INC_OPTIONS " \
#		"-Wno-unused-value -Wno-pointer-sign "		\
#		"-working-directory $WORKING_DIR "		\
#		"-c \"$CLANG_SOURCE\" -target bpf $CLANG_EMIT_LLVM -O2 -o - $LLVM_OPTIONS_PIPE"

# URL = https://lore.kernel.org/netdev/20220414134134.247912-1-yangjihong1@huawei.com/T/
# @ the end of ${URL}
#clang -O2 -emit-llvm -Xclang -disable-llvm-passes     \
# -D__GNUC__=1 -D__clang__=1                           \
# -I"${LINUX_SRC_HEADERS}"                             \
# -I"${CILIUM_EBPF}"/examples/headers                  \
# -I"${LIBBPF_HEADERS}"                                \
# -I"${PRJ_EBPF_MODULE_PATH}"/c                        \
# -L"${LIBBPF_OBJECTS}"                                \
#-c "${PRJ_EBPF_MODULE_PATH}"/c/"${PRJ_EBPF_MODULE}".c \
#-o "${PRJ_EBPF_MODULE_PATH}"/bpf_bpfel.o

#opt -O2 -mtriple=bpf-pc-linux | $(LLVM_DIS) | \
#llc -march=bpf $(LLC_FLAGS) -filetype=obj -o $@

"${GOPATH}"/bin/bpfo2go -target bpfel bpf

# 2 - big endian - 如arm
# -MF/dev/fd/3
#clang -O2 -mcpu=v1 -g -fno-ident -MD -MP -MF ./bpf_bpfeb.o.d     \
# -target bpfeb                                                   \
# -fdebug-prefix-map="${PRJ_EBPF_MODULE_PATH}"=.                  \
# -fdebug-compilation-dir .                                       \
# -Wunused-command-line-argument                                  \
# -D__BPF_TARGET_MISSING="\"Please provide -target bpfeb|bpfel\"" \
# -I"${LINUX_SRC_HEADERS}"                                        \
# -I"${CILIUM_EBPF}"/examples/headers                             \
# -I"${LIBBPF_HEADERS}"                                           \
# -I"${PRJ_EBPF_MODULE_PATH}"/c                                   \
# -L"${LIBBPF_OBJECTS}"                                           \
# -c "${PRJ_EBPF_MODULE_PATH}"/c/"${PRJ_EBPF_MODULE}".c           \
# -o "${PRJ_EBPF_MODULE_PATH}"/bpf_bpfeb.o

# ${GOPATH}/bin/bpfo2go bpfeb
# $GOPATH/bin/dlv --listen=:8606 --headless=true --api-version=2 --accept-multiclient exec $GOPATH/bin/bpfo2go -- -target bpfel bpf

exit 0
