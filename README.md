[![Build Status](https://travis-ci.org/onego-project/onego.svg?branch=master)](https://travis-ci.org/onego-project/onego)

# Onego
OpenNebula Go library is designed to manage enterprise clouds and virtualized data centers.

## Requirements
* Go 1.10 or newer to compile
* OpenNebula instance

## Installation
The recommended way to install this library is using `go get`:
```
go get -u github.com/onego-project/onego
```

## Usage examples
Short usage example expects OpenNebula server running at `localhost:2633/RPC2`.
```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/onego-project/onego"
	"github.com/onego-project/onego/blueprint"
	"github.com/onego-project/onego/errors"
)

var (
	endpoint = "http://localhost:2633/RPC2"
	name     = "oneadmin"
	password = "qwerty123"
	token    = name + ":" + password
)

func main() {
	// create client
	client := onego.CreateClient(endpoint, token, &http.Client{})
	if client == nil {
		log.Fatal(errors.ErrNoClient)
	}

	// prepare virtual machine name, memory and CPU size
	virtualMachineBlueprint := blueprint.CreateAllocateVirtualMachineBlueprint()
	if virtualMachineBlueprint == nil {
		log.Fatal(errors.ErrNoVirtualMachineBlueprint)
	}

	virtualMachineBlueprint.SetName("test-allocate")
	virtualMachineBlueprint.SetMemory(2048)
	virtualMachineBlueprint.SetCPU(4)

	// allocate virtual machine
	virtualMachine, err := client.VirtualMachineService.Allocate(context.TODO(), 
	    virtualMachineBlueprint, true)
	if err != nil {
		log.Fatal(err)
	}

	// get virtual machine info - ID
	id, err := virtualMachine.ID()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ID:", id)
}

```

## Contributing
1. [Fork onego library](https://github.com/onego-project/onego/fork)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
