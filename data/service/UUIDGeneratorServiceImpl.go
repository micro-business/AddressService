package service

import (
	"github.com/microbusinesses/Micro-Businesses-Core/system"
)

// UUIDGeneratorServiceImpl is the default implementation of UUID generator service
type UUIDGeneratorServiceImpl struct {
}

// GenerateRandomUUID generates random UUID value.
// Either the new random generated UUID or an error if something goes wrong.
func (UUIDGeneratorServiceImpl) GenerateRandomUUID() (system.UUID, error) {
	return system.RandomUUID()
}
