package loyalty_system

// Loyalty service
type Loyalty struct {
	Addr string
}

// New Loyalty instance
func New(addr string) *Loyalty {
	return &Loyalty{
		Addr: addr,
	}
}
