module github.com/tendermint/starport

go 1.16

require (
	github.com/99designs/keyring v1.1.6
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/Microsoft/hcsshim v0.8.17 // indirect
	github.com/blang/semver v3.5.1+incompatible
	github.com/briandowns/spinner v1.11.1
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/charmbracelet/glow v1.4.0
	github.com/containerd/containerd v1.5.2 // indirect
	github.com/cosmos/cosmos-sdk v0.44.0
	github.com/cosmos/go-bip39 v1.0.0
	github.com/docker/docker v20.10.7+incompatible
	github.com/emicklei/proto v1.9.0
	github.com/fatih/color v1.10.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-git/go-git/v5 v5.1.0
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gobuffalo/genny v0.6.0
	github.com/gobuffalo/logger v1.0.3
	github.com/gobuffalo/packd v1.0.0
	github.com/gobuffalo/packr/v2 v2.8.1 // indirect
	github.com/gobuffalo/plush v3.8.3+incompatible
	github.com/gobuffalo/plushgen v0.1.2
	github.com/goccy/go-yaml v1.9.2
	github.com/gogo/protobuf v1.3.3
	github.com/google/go-github/v37 v37.0.0
	github.com/gookit/color v1.2.7
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/rpc v1.2.0
	github.com/hashicorp/go-getter v1.5.9 // indirect
	github.com/iancoleman/strcase v0.1.3
	github.com/imdario/mergo v0.3.11
	github.com/jpillora/chisel v1.7.3
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-zglob v0.0.3
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/otiai10/copy v1.6.0
	github.com/pelletier/go-toml v1.9.3
	github.com/pkg/errors v0.9.1
	github.com/radovskyb/watcher v1.0.7
	github.com/rogpeppe/go-internal v1.7.0 // indirect
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/flutter v1.0.2
	github.com/tendermint/spm v0.1.7
	github.com/tendermint/spn v0.1.1-0.20211109105629-7de2146b06d9
	github.com/tendermint/tendermint v0.34.13
	github.com/tendermint/vue v0.1.53
	golang.org/x/mod v0.4.2
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
)
