package cmdb

type ResIpPortMeta struct {
	CompId       string `json:"compId,omitempty" yaml:"comp-id"`
	CompName     string `json:"compName,omitempty" yaml:"comp-name"`
	CompTypeId   string `json:"compTypeId,omitempty" yaml:"comp-type-id"`
	CompTypeName string `json:"compTypeName,omitempty" yaml:"comp-type-name"`
}
