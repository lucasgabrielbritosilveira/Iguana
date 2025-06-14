package bus

import (
	"iguana/memory"
)

type Bus struct {
	RAM memory.Data
}

func (bus *Bus) Write(address uint16, data uint8) {

}

func (bus *Bus) Read(address uint16) uint8 {
	return 0xF
}
