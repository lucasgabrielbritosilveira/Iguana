package bus

import (
	"errors"
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

func validateAddress(address uint16) error {
	if address <= 0xFFFF {
		return nil
	}
	return errors.New("invalid Address")
}
