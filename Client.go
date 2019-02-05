package onego

import (
	"net/http"

	"github.com/onego-project/onego/services"
	"github.com/onego-project/xmlrpc"
)

// Client structure contains the services to manage resources
type Client struct {
	UserService             services.UserService
	TokenService            services.TokenService
	GroupService            services.GroupService
	DatastoreService        services.DatastoreService
	ImageService            services.ImageService
	HostService             services.HostService
	ClusterService          services.ClusterService
	VirtualNetworkService   services.VirtualNetworkService
	ReservationService      services.ReservationService
	AddressRangeService     services.AddressRangeService
	LeaseService            services.LeaseService
	TemplateService         services.TemplateService
	VirtualMachineService   services.VirtualMachineService
	NetworkInterfaceService services.NetworkInterfaceService
}

// CreateClient creates Client with endpoint, token and http client
func CreateClient(endpoint, token string, client *http.Client) *Client {
	rpc := &services.RPC{Client: xmlrpc.NewClient(endpoint, client), Token: token}

	return &Client{UserService: services.UserService{Service: services.Service{RPC: rpc}},
		TokenService:            services.TokenService{Service: services.Service{RPC: rpc}},
		GroupService:            services.GroupService{Service: services.Service{RPC: rpc}},
		DatastoreService:        services.DatastoreService{Service: services.Service{RPC: rpc}},
		ImageService:            services.ImageService{Service: services.Service{RPC: rpc}},
		HostService:             services.HostService{Service: services.Service{RPC: rpc}},
		ClusterService:          services.ClusterService{Service: services.Service{RPC: rpc}},
		VirtualNetworkService:   services.VirtualNetworkService{Service: services.Service{RPC: rpc}},
		ReservationService:      services.ReservationService{Service: services.Service{RPC: rpc}},
		AddressRangeService:     services.AddressRangeService{Service: services.Service{RPC: rpc}},
		LeaseService:            services.LeaseService{Service: services.Service{RPC: rpc}},
		TemplateService:         services.TemplateService{Service: services.Service{RPC: rpc}},
		VirtualMachineService:   services.VirtualMachineService{Service: services.Service{RPC: rpc}},
		NetworkInterfaceService: services.NetworkInterfaceService{Service: services.Service{RPC: rpc}},
	}
}
