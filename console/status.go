package console

import (
	"fmt"
)

func init() {
	commands["status"] = &CmdStatus{}
}

type CmdStatus struct {
	rpcMethod string
	rpcParams string
	rpcResult string
}

func (self *CmdStatus) Usage(name string) string {
	return fmt.Sprintf("\n\tUsage: %s [cfg_opts...{-h}] status", name)
}

func (self *CmdStatus) defaults() error {
	self.rpcMethod = "Responder.Status"
	return nil
}

func (self *CmdStatus) FromArgs(args []string) error {
	self.defaults()
	return nil
}

func (self *CmdStatus) RpcMethod() string {
	return self.rpcMethod
}

func (self *CmdStatus) RpcParams() interface{} {
	return &self.rpcParams
}

func (self *CmdStatus) RpcResult() interface{} {
	return &self.rpcResult
}
