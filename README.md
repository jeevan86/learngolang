brew install FiloSottile/musl-cross/musl-cross
brew install libpcap

CGO_ENABLED=1                                           \
CGO_CFLAGS="-I/usr/local/Cellar/libpcap/1.10.1/include" \
CGO_LDFLAGS="-L/usr/local/Cellar/libpcap/1.10.1/lib"    \
CC=x86_64-linux-musl-gcc                                \
CXX=x86_64-linux-musl-g++                               \
GOOS=linux GOARCH=amd64                                 \
go build -o dist/gopcap gopackettest/main

go build -gcflags="-d pctab=pctoinline" ...


go build -gcflags="-d pctab=pctoinline" -o dist/gopcap gopackettest/main