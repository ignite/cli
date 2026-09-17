package runtime

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	"github.com/gnolang/gno/gnovm/stdlibs/chain/runtime/unsafe"
	"github.com/gnolang/gno/gnovm/stdlibs/internal/execctx"
	"github.com/gnolang/gno/tm2/pkg/crypto"
)

func TestPreviousRealmIsOrigin(t *testing.T) {
	var (
		user = gno.DerivePkgBech32Addr("user1.gno")
		ctx  = execctx.ExecContext{
			OriginCaller: user,
		}
		msgCallFrame = gno.Frame{LastPackage: &gno.PackageValue{PkgPath: "main"}}
		msgRunFrame  = gno.Frame{LastPackage: &gno.PackageValue{PkgPath: "gno.land/e/g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5/run"}}
	)
	type expectations struct {
		addr         crypto.Bech32Address
		pkgPath      string
		isOriginCall bool
		doesPanic    bool
	}
	tests := []struct {
		name                 string
		machine              *gno.Machine
		expectedAddr         crypto.Bech32Address
		expectedPkgPath      string
		expectedIsOriginCall bool
	}{
		{
			name: "no frames",
			machine: &gno.Machine{
				Context: ctx,
				Frames:  []gno.Frame{},
			},
			expectedAddr:         user,
			expectedPkgPath:      "",
			expectedIsOriginCall: false,
		},
		{
			name: "one frame w/o LastPackage",
			machine: &gno.Machine{
				Context: ctx,
				Frames: []gno.Frame{
					{LastPackage: nil},
				},
			},
			expectedAddr:         user,
			expectedPkgPath:      "",
			expectedIsOriginCall: false,
		},
		{
			name: "one package frame",
			machine: &gno.Machine{
				Context: ctx,
				Frames: []gno.Frame{
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/p/xxx"}},
				},
			},
			expectedAddr:         user,
			expectedPkgPath:      "",
			expectedIsOriginCall: false,
		},
		{
			name: "one realm frame",
			machine: &gno.Machine{
				Context: ctx,
				Frames: []gno.Frame{
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/r/xxx"}},
				},
			},
			expectedAddr:         user,
			expectedPkgPath:      "",
			expectedIsOriginCall: false,
		},
		{
			name: "one msgCall frame",
			machine: &gno.Machine{
				Context: ctx,
				Frames: []gno.Frame{
					msgCallFrame,
				},
			},
			expectedAddr:         user,
			expectedPkgPath:      "",
			expectedIsOriginCall: true,
		},
		{
			name: "one msgRun frame",
			machine: &gno.Machine{
				Context: ctx,
				Frames: []gno.Frame{
					msgRunFrame,
				},
			},
			expectedAddr:         user,
			expectedPkgPath:      "",
			expectedIsOriginCall: false,
		},
		{
			name: "one package frame and one msgCall frame",
			machine: &gno.Machine{
				Context: ctx,
				Frames: []gno.Frame{
					msgCallFrame,
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/p/xxx"}},
				},
			},
			expectedAddr:         user,
			expectedPkgPath:      "",
			expectedIsOriginCall: true,
		},
		{
			name: "one realm frame and one msgCall frame",
			machine: &gno.Machine{
				Context: ctx,
				Frames: []gno.Frame{
					msgCallFrame,
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/r/xxx"}},
				},
			},
			expectedAddr:         user,
			expectedPkgPath:      "",
			expectedIsOriginCall: true,
		},
		{
			name: "one package frame and one msgRun frame",
			machine: &gno.Machine{
				Context: ctx,
				Frames: []gno.Frame{
					msgRunFrame,
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/p/xxx"}},
				},
			},
			expectedAddr:         user,
			expectedPkgPath:      "",
			expectedIsOriginCall: false,
		},
		{
			name: "one realm frame and one msgRun frame",
			machine: &gno.Machine{
				Context: ctx,
				Frames: []gno.Frame{
					msgRunFrame,
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/r/xxx"}},
				},
			},
			expectedAddr:         user,
			expectedPkgPath:      "",
			expectedIsOriginCall: false,
		},
		{
			name: "multiple frames with one realm",
			machine: &gno.Machine{
				Context: ctx,
				Frames: []gno.Frame{
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/p/xxx"}},
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/p/xxx"}},
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/r/xxx"}},
				},
			},
			expectedAddr:         user,
			expectedPkgPath:      "",
			expectedIsOriginCall: false,
		},
		{
			name: "multiple frames with multiple realms",
			machine: &gno.Machine{
				Context: ctx,
				Frames: []gno.Frame{
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/r/zzz"}},
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/r/zzz"}},
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/r/yyy"}},
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/p/yyy"}},
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/p/xxx"}},
					{LastPackage: &gno.PackageValue{PkgPath: "gno.land/r/xxx"}},
				},
			},
			expectedAddr:         gno.DerivePkgBech32Addr("gno.land/r/yyy"),
			expectedPkgPath:      "gno.land/r/yyy",
			expectedIsOriginCall: false,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("fail", i)
				}
			}()
			assert := assert.New(t)

			addr, pkgPath := unsafe.X_getRealm(tt.machine, 1)
			isOrigin := isOriginCall(tt.machine)

			assert.Equal(string(tt.expectedAddr), addr)
			assert.Equal(tt.expectedPkgPath, pkgPath)
			assert.Equal(tt.expectedIsOriginCall, isOrigin)
		})
	}
}
