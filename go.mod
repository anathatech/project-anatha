module github.com/anathatech/project-anatha

go 1.13

require (
	github.com/anathatech/cosmosd v0.2.2
	github.com/btcsuite/btcd v0.0.0-20190523000118-16327141da8c // indirect
	github.com/cosmos/cosmos-sdk v0.38.4
	github.com/gorilla/mux v1.7.3
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rcrowley/go-metrics v0.0.0-20190704165056-9c2d0518ed81 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.6
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.3
	github.com/tendermint/tm-db v0.5.0
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/cosmos/cosmos-sdk => github.com/anathatech/cosmos-sdk v0.38.5-0.20201102214748-0d5a47f063a8
