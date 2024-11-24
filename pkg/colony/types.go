package colony

// Room represents a room in the ant farm
type Room struct {
	Name    string
	X       int
	Y       int
	IsStart bool
	IsEnd   bool
}

// Tunnel represents a connection between rooms
type Tunnel struct {
	From string
	To   string
}

// Colony represents the entire ant farm structure
type Colony struct {
	NumAnts int
	Rooms   map[string]*Room
	Tunnels []Tunnel
	Start   string
	End     string
	Input   []string
}

// NewColony creates a new Colony instance
func NewColony() *Colony {
	return &Colony{
		Rooms:   make(map[string]*Room),
		Tunnels: make([]Tunnel, 0),
		Input:   make([]string, 0),
	}
}
