package onego

import (
	"github.com/onego-project/xmlrpc"
	"github.com/owlet123/onego/virtualmachine"
)

// Client structure contains XML-RPC client and virtual machine
type Client struct {
	XMLRPCClient          xmlrpc.Client
	VirtualMachineService virtualmachine.Service
}
