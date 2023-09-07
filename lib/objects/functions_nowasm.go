//go:build !wasm
// +build !wasm

package objects


type value struct{}

func FromSourceCode(this value, args []value) any {
	return value{}
}

func FromPackages(this value, args []value) any {
	return value{}
}

func GetRepositoryDetails(this value, args []value) any {
	return value{}
}

func FetchRepository(this value, args []value) any {
	return value{}
}
