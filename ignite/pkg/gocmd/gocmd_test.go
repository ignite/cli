package gocmd_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ignite/cli/ignite/pkg/gocmd"
)

func TestIsInstallError(t *testing.T) {
	assert.False(t, gocmd.IsInstallError(errors.New("oups")))

	err := errors.New(`error while running command go install github.com/cosmos/gogoproto/protoc-gen-gocosmos google.golang.org/protobuf/cmd/protoc-gen-go github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2: no required module provides package github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2; to add it:
		go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2`)
	assert.True(t, gocmd.IsInstallError(err))
}
