module github.com/jeevan86/learngolang

go 1.17

require (
	github.com/cilium/ebpf v0.9.0
	github.com/google/gopacket v1.1.19
	github.com/iovisor/gobpf v0.2.0
	github.com/jeevan86/lf4go v0.4.0
	github.com/reactivex/rxgo/v2 v2.5.0
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0
)

require (
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/natefinch/lumberjack/v3 v3.0.0-alpha // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	github.com/teivah/onecontext v1.3.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.11 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	github.com/cilium/ebpf v0.9.0 => ../../cilium/ebpf
	// 先本地测试
	github.com/jeevan86/lf4go v0.4.0 => ../lf4go
)
