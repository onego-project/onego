package onego

import (
	"github.com/onego-project/onego/services"
	"github.com/onego-project/xmlrpc"
	"net/http"
)

// Client structure contains the services to manage resources
type Client struct {
	UserService  services.UserService
	GroupService services.GroupService
}

// CreateClient creates Client with endpoint, token and http client
func CreateClient(endpoint, token string, client *http.Client) *Client {
	rpc := &services.RPC{Client: xmlrpc.NewClient(endpoint, client), Token: token}

	return &Client{UserService: services.UserService{Service: services.Service{RPC: rpc}},
		GroupService: services.GroupService{Service: services.Service{RPC: rpc}}}
}
