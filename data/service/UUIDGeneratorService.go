package service

import (
	. "github.com/microbusinesses/Micro-Businesses-Core/system"
)

// Defines the interface that will generate random UUID, will be used to inject UUID generator service implementation to Address service.
type UUIDGeneratorService interface {
	// Generates random UUID value.
	// Either the new random generated UUID or an error if something goes wrong.
	GenerateRandomUUID() (UUID, error)
}
