package memory

import (
	"fmt"
	"strconv"
	"time"

	"github.com/abu-lang/goabu/memory"
)

// New creates a new memory, based on the basic resources
func NewBasicMemory(items map[string]map[string]any) (memory.ResourceController, error) {
	// I create an empty basic memory...
	mem := memory.MakeResources()
	// ... and I range over the items to initialize it, with the provided initialization value or a default
	for t, m := range items {
		switch t {
		case "String":
			for name, v := range m {
				val, ok := v.(string)
				if !ok {
					return nil, fmt.Errorf("value for key %q is not a string", name)
				}
				mem.Text[name] = val
			}
		case "Integer":
			for name, v := range m {
				val, ok := v.(int64)
				if !ok {
					return nil, fmt.Errorf("value for key %q is not a string", name)
				}
				mem.Integer[name] = val
			}
		case "Bool":
			for name, v := range m {
				val, ok := v.(bool)
				if !ok {
					return nil, fmt.Errorf("value for key %q is not a string", name)
				}
				mem.Bool[name] = val
			}
		case "Float":
			for name, v := range m {
				val, ok := v.(float64)
				if !ok {
					return nil, fmt.Errorf("value for key %q is not a string", name)
				}
				mem.Float[name] = val
			}
		}
	}
	return mem, nil
}

func getBasicMemoryBool(initvalue string) (bool, error) {
	if initvalue == "" {
		return false, nil
	} else {
		return strconv.ParseBool(initvalue)
	}
}

func getBasicMemoryInteger(initvalue string) (int64, error) {
	if initvalue == "" {
		return 0, nil
	} else {
		return strconv.ParseInt(initvalue, 10, 64)
	}
}

func getBasicMemoryFloat(initvalue string) (float64, error) {
	if initvalue == "" {
		return 0, nil
	} else {
		return strconv.ParseFloat(initvalue, 64)
	}
}

func getBasicMemoryText(initvalue string) (string, error) {
	return initvalue, nil
}

func getBasicMemoryTime(initvalue string) (time.Time, error) {
	if initvalue == "" {
		return time.Now(), nil
	} else {
		return time.Parse(time.Stamp, initvalue)
	}
}
