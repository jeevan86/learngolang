package panics

import (
	"errors"
	"github.com/jeevan86/learngolang/pkg/util/str"
	"reflect"
	"testing"
)

func TestSafeRet(t *testing.T) {
	type args struct {
		f func() interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet interface{}
		wantSta string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				f: func() interface{} {
					return "success"
				},
			},
			wantRet: "success",
			wantSta: str.EMPTY(),
			wantErr: false,
		},
		{
			name: "panic",
			args: args{
				f: func() interface{} {
					panic(errors.New("panic"))
				},
			},
			wantRet: nil,
			wantSta: "err => panic",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, gotSta, err := SafeRet(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeRet() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("SafeRet() gotRet = %v, want %v", gotRet, tt.wantRet)
			}
			if tt.wantSta == str.EMPTY() {
				if gotSta != tt.wantSta {
					t.Errorf("SafeRet() gotSta = %v, want %v", gotSta, tt.wantSta)
				}
			}
		})
	}
}

func TestSafeRun(t *testing.T) {
	type args struct {
		f func()
	}
	tests := []struct {
		name    string
		args    args
		wantSta string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				f: func() {

				},
			},
			wantSta: str.EMPTY(),
			wantErr: false,
		},
		{
			name: "panic",
			args: args{
				f: func() {
					panic(errors.New("panic"))
				},
			},
			wantSta: "err => panic",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSta, err := SafeRun(tt.args.f)
			if (err != nil) != tt.wantErr {
				t.Errorf("SafeRun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantSta == str.EMPTY() {
				if gotSta != tt.wantSta {
					t.Errorf("SafeRet() gotSta = %v, want %v", gotSta, tt.wantSta)
				}
			}
		})
	}
}
