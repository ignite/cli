module github.com/tendermint/starport

go 1.16

require (
	github.com/AlecAivazis/survey/v2 v2.1.1
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/briandowns/spinner v1.11.1
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/charmbracelet/glow v1.4.0
	github.com/cosmos/cosmos-sdk v0.42.0
	github.com/cosmos/go-bip39 v1.0.0
	github.com/cosmos/relayer v1.0.0-rc1.0.20210326125444-76eb658fb20a
	github.com/dariubs/percent v0.0.0-20200128140941-b7801cf1c7e2
	github.com/emicklei/proto v1.9.0
	github.com/fatih/color v1.10.0
	github.com/gertd/go-pluralize v0.1.7
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/go-git/go-git/v5 v5.1.0
	github.com/gobuffalo/genny v0.6.0
	github.com/gobuffalo/packd v1.0.0
	github.com/gobuffalo/plush v3.8.3+incompatible
	github.com/gobuffalo/plushgen v0.1.2
	github.com/goccy/go-yaml v1.8.0
	github.com/gookit/color v1.2.7
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/iancoleman/strcase v0.1.3
	github.com/imdario/mergo v0.3.11
	github.com/improbable-eng/grpc-web v0.13.0
	github.com/jpillora/chisel v1.7.3
	github.com/kr/pretty v0.1.0
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-zglob v0.0.3
	github.com/mwitkow/grpc-proxy v0.0.0-20181017164139-0f1106ef9c76
	github.com/olekukonko/tablewriter v0.0.5
	github.com/otiai10/copy v1.4.2
	github.com/pelletier/go-toml v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/radovskyb/watcher v1.0.7
	github.com/rakyll/statik v0.1.7
	github.com/rdegges/go-ipify v0.0.0-20150526035502-2d94a6a86c40
	github.com/rs/cors v1.7.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/spn v0.0.0-20210323090228-7c451f9fe1b5
	github.com/tendermint/tendermint v0.34.8
	github.com/tendermint/vue v0.0.0
	golang.org/x/crypto v0.0.0-20210317152858-513c2a44f670 // indirect
	golang.org/x/mod v0.4.1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d
	google.golang.org/grpc v1.35.0
)

replace github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4

replace github.com/tendermint/vue => ../vue

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

replace github.com/cosmos/relayer => github.com/cosmos/relayer v0.9.0
