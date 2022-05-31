package client

import (
	"github.com/jeevan86/learngolang/pkg/pcap/protocol/ip/base"
	"testing"
)

func Test_converter_convert(t *testing.T) {
	type fields struct {
		internal internalConverter
	}
	type args struct {
		m *base.OutputStruct
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   interface{}
	}{
		{
			name: "success",
			fields: fields{
				internal: newConverter(),
			},
			args: args{
				m: new(base.OutputStruct),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &converter{
				internal: tt.fields.internal,
			}
			got := c.convert(tt.args.m)
			if got == nil {
				t.Errorf("convert() = nil, want <not nill>")
			}
		})
	}
}
