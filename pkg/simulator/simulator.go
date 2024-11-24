package simulator

import (
	"fmt"
	"strings"

	"lem-in/pkg/colony"
	"lem-in/pkg/pathfinder"
)

// Ant represents an ant in the simulation
type Ant struct {
	ID       int
	Path     pathfinder.Path
	Position int
}

// SimulateMovement simulates the movement of ants through the colony
func SimulateMovement(c *colony.Colony) []string {
	paths := pathfinder.FindPaths(c)
	if len(paths) == 0 {
		return []string{"ERROR: invalid data format"}
	}

	// Print the input
	for _, line := range c.Input {
		fmt.Println(line)
	}
	fmt.Println()

	var moves []string
	ants := make([]*Ant, c.NumAnts)
	for i := range ants {
		ants[i] = &Ant{
			ID:       i + 1,
			Path:     paths[0],
			Position: -1,
		}
	}

	for {
		turnMoves := make([]string, 0)
		moveMade := false
		roomOccupancy := make(map[string]bool)

		// Try to move each ant
		for _, ant := range ants {
			if ant.Position == len(ant.Path)-1 {
				continue // Ant has reached the end
			}

			var nextRoom string
			if ant.Position == -1 {
				nextRoom = ant.Path[0]
			} else {
				nextRoom = ant.Path[ant.Position+1]
			}

			// Check if room is available
			if nextRoom != c.Start && nextRoom != c.End && roomOccupancy[nextRoom] {
				continue
			}

			if ant.Position == -1 {
				ant.Position = 0
				turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", ant.ID, nextRoom))
				moveMade = true
				roomOccupancy[nextRoom] = true
			} else {
				ant.Position++
				turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", ant.ID, nextRoom))
				moveMade = true
				roomOccupancy[nextRoom] = true
			}
		}

		if !moveMade {
			break
		}

		moves = append(moves, strings.Join(turnMoves, " "))
	}

	return moves
}
