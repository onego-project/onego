package requests

import (
	"fmt"
	"github.com/onego-project/onego/blueprint"
)

// ResizeRequest struct
type ResizeRequest struct {
	CPU    float32
	VCPU   int
	Memory int
}

// Render method
func (rr ResizeRequest) Render() string {
	bpmap := map[string]blueprint.Node{
		"cpu":    {Value: fmt.Sprintf("%.2f", rr.CPU)},
		"vcpu":   {Value: fmt.Sprintf("%d", rr.VCPU)},
		"memory": {Value: fmt.Sprintf("%d", rr.Memory)},
	}
	bp := blueprint.Blueprint{Values: bpmap, Root: "template"}

	return bp.Render()
}
