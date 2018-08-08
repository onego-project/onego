package onego

import (
	"github.com/onego-project/xmlrpc"
	"github.com/onego-project/onego/services"
)

// Client structure contains XML-RPC client and virtual machine
type Client struct {
	XMLRPCClient          xmlrpc.Client
	VirtualMachineService services.Service
}
