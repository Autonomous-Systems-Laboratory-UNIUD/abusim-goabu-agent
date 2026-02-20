package memory

import (
	"errors"

	resources "github.com/Autonomous-Systems-Laboratory-UNIUD/aburos/rosresources"
	vehicles "github.com/Autonomous-Systems-Laboratory-UNIUD/aburos/vehicles"
)

// New creates a new memory, based on the passed memory controller and items
func New(controller string, items map[string]map[string]any) (resources.ROSresources, error) {
	// I check the controller type and I return the correct implementation
	switch controller {
	case "basic":
		base, err := NewBasicMemory(items)
		if err != nil {
			return nil, err
		}
		mem := resources.NewBaseResourceController(items["Text"]["id"].(string))
		mem.Enclose(base.GetResources())
		return mem, nil
	case "copter":
		base, err := NewBasicMemory(items)
		if err != nil {
			return nil, err
		}
		vec, err := vehicles.NewCopterVehicleV2(items["Text"]["id"].(string), "", nil, nil)
		if err != nil {
			return nil, err
		}
		mem := resources.NewCopterResourceV2(vec)
		mem.Resources.Enclose(base.GetResources())
		return mem, nil
	case "plane":
		return nil, errors.New("not yet supported")
	case "sub":
		base, err := NewBasicMemory(items)
		if err != nil {
			return nil, err
		}
		vec, err := vehicles.NewSubVehicle(items["Text"]["id"].(string), "", nil, nil)
		if err != nil {
			return nil, err
		}
		mem := resources.NewSubResource(vec)
		mem.Resources.Enclose(base.GetResources())
		return mem, nil
	case "rover":
		return nil, errors.New("not yet supported")

	default:
		// If an invalid controller is passed, I raise an error
		return nil, errors.New("unsupported controller")
	}
}
