package client

import (
	"github.com/jeevan86/learngolang/pkg/capture/protocol/ip/base"
	"github.com/jeevan86/learngolang/pkg/util/panics"
	"strings"
)

type internalConverter interface {
	convert(m *base.OutputStruct) interface{}
}

type converter struct {
	internal internalConverter
}

func (c *converter) convert(m *base.OutputStruct) interface{} {
	ret, sta, err := panics.SafeRet(
		func() interface{} {
			return c.internal.convert(m)
		},
	)
	if err != nil {
		logger.Error("Failed to convert: %s, %s", err.Error(), sta)
	}
	return ret
}

func newConverter() *converter {
	var internal internalConverter
	serverType := *collectConfig.ServerType
	switch strings.ToLower(serverType) {
	case TypeGrpc:
		conv := grpcConverter("grpcConverter")
		internal = &conv
		break
	case TypeHttp:
		conv := defaultConverter("httpConverter")
		internal = &conv
		break
	case TypeLog:
		conv := defaultConverter("logConverter")
		internal = &conv
		break
	default:
		conv := defaultConverter("logConverter")
		internal = &conv
		break
	}
	return &converter{
		internal: internal,
	}
}
