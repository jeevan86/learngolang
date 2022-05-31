package base

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
)

const EbpfMapTypeFieldName = "typ"
const EbpfStructTagName = "ebpf"

var loaders = make(map[uint32]BpfLoader)

func Register(k uint32, l BpfLoader) {
	loaders[k] = l
}

type InjectedMapSpec struct {
	MapK uint32
	MapT ebpf.MapType
}

type InjectedFunSpec struct {
	Name string
	Maps []*InjectedMapSpec
}

type EbpfProgram struct {
	InjectedFunSpec
	Link link.Link
	Prog *ebpf.Program
}

// newKpSysFuncProgram
// *ebpf.Program#Info() require kernel 4.10
func NewEbpfProgram(n string, p *ebpf.Program, m []*InjectedMapSpec) *EbpfProgram {
	return &EbpfProgram{
		InjectedFunSpec: InjectedFunSpec{
			Name: n,
			Maps: m,
		},
		Prog: p,
	}
}

// EbpfMaps 一个模块中一个类型的Map只能有一个
type EbpfMapMap map[ebpf.MapType]*ebpf.Map
type ProgramMap map[string]*EbpfProgram
type InjectedFunMap map[string]*InjectedFunSpec

type BpfLoader interface {
	Start() error
	Close()
	GetPrograms() ProgramMap
	Lookup(*EbpfProgram, map[ebpf.MapType]interface{}) error
}
