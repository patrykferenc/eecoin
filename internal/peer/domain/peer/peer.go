package peer

import (
	"fmt"
)

type Status int

const (
	StatusHealthy Status = iota
	StatusUnhealthy
	StatusUnknown
)

var statusName = map[Status]string{
	StatusHealthy:   "healthy",
	StatusUnhealthy: "unhealthy",
	StatusUnknown:   "unknown",
}

func (s Status) String() string {
	return statusName[s]
}

type Peer struct {
	Host   string
	Status Status
}

func (p Peer) String() string {
	return fmt.Sprintf("%s (%s)", p.Host, p.Status)
}
