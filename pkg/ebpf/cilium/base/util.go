package base

import (
	"github.com/cilium/ebpf"
	"reflect"
	"strings"
)

func ReflectGetEbpfMaps(bpfMaps interface{}) EbpfMapMap {
	bpfMapsType := reflect.TypeOf(bpfMaps)
	ebpfMaps := make(EbpfMapMap)
	for i := 0; i < bpfMapsType.NumField(); i++ {
		field := bpfMapsType.Field(i)
		ebpfMap, _ := reflect.ValueOf(bpfMaps).FieldByName(field.Name).Interface().(*ebpf.Map)
		// panic: reflect: call of reflect.Value.FieldByName on ptr Value
		// panic: reflect.Value.Interface: cannot return value obtained from unexported field or method
		ebpfMapType := ebpf.MapType(reflect.ValueOf(*ebpfMap).FieldByName(EbpfMapTypeFieldName).Uint())
		//ebpfMapName := field.Tag.Get("ebpf") // kprobe_connect
		ebpfMaps[ebpfMapType] = ebpfMap
	}
	return ebpfMaps
}

func ReflectGetPrograms(injectedFunMap InjectedFunMap, bpfPrograms interface{}, prefix string) ProgramMap {
	programs := make(ProgramMap)
	bpfProgramsType := reflect.TypeOf(bpfPrograms)
	// panic: reflect: NumField of non-struct type *kprobe.bpfPrograms
	for i := 0; i < bpfProgramsType.NumField(); i++ {
		field := bpfProgramsType.Field(i)
		program, _ := reflect.ValueOf(bpfPrograms).FieldByName(field.Name).Interface().(*ebpf.Program)
		ebpfKpFuncName := field.Tag.Get(EbpfStructTagName) // kprobe_connect
		sysFuncName := strings.ReplaceAll(ebpfKpFuncName, prefix, "")
		spec, exists := injectedFunMap[sysFuncName]
		if exists && spec != nil {
			programs[sysFuncName] = NewEbpfProgram(spec.Name, program, spec.Maps)
		}
	}
	return programs
}

func Lookup(m *ebpf.Map, k uint32) interface{} {
	switch m.Type() {
	case ebpf.PerfEventArray:
		var value uint64
		_ = m.Lookup(k, &value)
		return value
	case ebpf.Array:
		var value uint64
		_ = m.Lookup(k, &value)
		return value
	default:
		return nil
	}
}
