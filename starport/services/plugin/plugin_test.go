package plugin

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.Println("Build plugin")

	// TODO: Build plugin.
	// _ = exec.Command("pushd", "fixtures").Run()
	// cmd := exec.Command("go", "build", "-buildmode=plugin", "./")
	// err := cmd.Run()

	// if err != nil {
	// 	log.Panic(err)
	// }

	// _ = exec.Command("popd").Run()
}

func Test_LoadPlugin(t *testing.T) {
	tests := []struct {
		Desc      string
		Name      string
		Symbol    string
		ExpectErr error
	}{
		{
			Desc:      "Load successfully",
			Name:      "fixture",
			Symbol:    "./fixtures/fixtures.so",
			ExpectErr: nil,
		},

		{
			Desc:      "Failed to load not existing symbol",
			Name:      "not-exist",
			Symbol:    "./fixtures/not-exist-symbol.so",
			ExpectErr: ErrSymbolNotExist,
		},
	}

	for _, test := range tests {
		p, err := LoadPlugin(test.Name, test.Symbol)

		assert.Equal(t, test.ExpectErr, err)

		_ = p
		if test.ExpectErr == nil {
			assert.NotNil(t, p)
			assert.Equal(t, test.Name, p.Name())

			_ = p.List()
		} else {
			assert.Nil(t, p)
		}
	}
}

func Test_Execute(t *testing.T) {
	tests := []struct {
		Desc      string
		FuncName  string
		Params    []string
		IsSuccess bool
	}{
		{
			Desc:      "Success",
			FuncName:  "Add",
			Params:    []string{"10", "1"},
			IsSuccess: true,
		},

		{
			Desc:      "Failed not exist",
			FuncName:  "NotExistfunction",
			Params:    []string{},
			IsSuccess: false,
		},

		{
			Desc:      "Invalid paramete type",
			FuncName:  "Add",
			Params:    []string{"10", "a"},
			IsSuccess: false,
		},
	}

	p, err := LoadPlugin("fixture", "./fixtures/fixtures.so")
	assert.Nil(t, err)

	for _, test := range tests {
		err = p.Execute(test.FuncName, test.Params)

		if test.IsSuccess {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}

}
