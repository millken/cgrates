package console

import (
	"fmt"
	"path/filepath"
)

// Console Command interface
type Commander interface {
	FromArgs(args []string) error // Load data from os arguments or flag.Args()
	Usage(string) string          // usage message
	RpcMethod() string            // Method which should be called remotely
	RpcParams() interface{}       // Parameters to send out on rpc
	RpcResult() interface{}       // Only requirement is to have a String method to print on console
	defaults() error              // set defaults wherever necessary
}

// Process args and return right command Value or error
func GetCommandValue(args []string) (Commander, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("\n\tUsage: %s [cfg_opts...{-h}] <command>\n", filepath.Base(args[0]))
	}
	cmd := args[1]
	var cmdVal Commander
	switch cmd {
	case "status":
		cmdVal = &CmdStatus{}
	case "get_balance":
		cmdVal = &CmdGetBalance{}
	default:
		return nil, fmt.Errorf("\n\tUsage: %s [cfg_opts...{-h}] <status|get_balance>\n", filepath.Base(args[0]))
	}
	if err := cmdVal.FromArgs(args); err != nil {
		return nil, err
	}
	return cmdVal, nil
}
