package simulator

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"lem-in/pkg/colony"
	"lem-in/pkg/pathfinder"
)

// Ant represents an ant in the simulation
type Ant struct {
	ID        int
	Path      pathfinder.Path
	Position  int
	PathIndex int
	Delay     int // Number of turns to wait before starting
}

// PathState tracks the state of a path for optimization
type PathState struct {
	path       pathfinder.Path
	antsCount  int
	lastAntEnd float64
	efficiency float64 // How efficiently ants move through this path
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

	// Initialize path states with efficiency calculations
	pathStates := initializePathStates(paths, c)

	// Distribute ants optimally across paths
	ants := distributeAnts(c.NumAnts, pathStates)

	// Simulate movements
	return simulateAntMovements(ants, c, pathStates)
}

// initializePathStates initializes path states with efficiency calculations
func initializePathStates(paths []pathfinder.Path, c *colony.Colony) []PathState {
	pathStates := make([]PathState, len(paths))
	
	for i, path := range paths {
		efficiency := calculatePathEfficiency(path, c)
		pathStates[i] = PathState{
			path:       path,
			antsCount:  0,
			lastAntEnd: float64(len(path) - 1),
			efficiency: efficiency,
		}
	}

	return pathStates
}

// calculatePathEfficiency calculates how efficiently ants can move through a path
func calculatePathEfficiency(path pathfinder.Path, c *colony.Colony) float64 {
	if len(path) <= 2 {
		return 1.0 // Direct path is most efficient
	}

	// Consider path length
	lengthFactor := 1.0 / float64(len(path))

	// Consider room connectivity (more connected rooms = better flow)
	connectivityFactor := 0.0
	for _, room := range path[1:len(path)-1] {
		connections := 0
		for _, tunnel := range c.Tunnels {
			if tunnel.From == room || tunnel.To == room {
				connections++
			}
		}
		connectivityFactor += float64(connections)
	}
	connectivityFactor /= float64(len(path) - 2)
	connectivityFactor = math.Min(connectivityFactor/4.0, 1.0) // Normalize

	return (lengthFactor*0.6 + connectivityFactor*0.4)
}

// distributeAnts distributes ants optimally across available paths
func distributeAnts(numAnts int, pathStates []PathState) []*Ant {
	ants := make([]*Ant, numAnts)
	
	// Calculate initial distribution
	remainingAnts := numAnts
	antIndex := 0

	for remainingAnts > 0 {
		bestPath := -1
		bestTime := math.MaxFloat64

		for i, state := range pathStates {
			// Calculate estimated time for this path
			timeEstimate := calculateTimeEstimate(state, remainingAnts)
			if timeEstimate < bestTime {
				bestTime = timeEstimate
				bestPath = i
			}
		}

		if bestPath == -1 {
			break
		}

		// Calculate optimal number of ants for this path
		antsForPath := calculateOptimalAntsForPath(pathStates[bestPath], remainingAnts, bestTime)
		
		// Assign ants to this path
		for i := 0; i < antsForPath; i++ {
			delay := calculateAntDelay(pathStates[bestPath], i)
			ants[antIndex] = &Ant{
				ID:        antIndex + 1,
				Path:      pathStates[bestPath].path,
				PathIndex: bestPath,
				Position:  -1,
				Delay:    delay,
			}
			antIndex++
			remainingAnts--
		}

		// Update path state
		pathStates[bestPath].antsCount += antsForPath
		pathStates[bestPath].lastAntEnd += float64(antsForPath) * 0.8
	}

	return ants
}

// calculateTimeEstimate estimates completion time for a path
func calculateTimeEstimate(state PathState, remainingAnts int) float64 {
	pathLength := float64(len(state.path) - 1)
	currentLoad := float64(state.antsCount)
	
	// Base time is path length
	baseTime := pathLength
	
	// Add time for ant interference
	interferenceTime := (currentLoad * 0.8) * (1.0 - state.efficiency)
	
	// Consider remaining capacity
	capacityFactor := math.Max(0, 1.0-currentLoad/pathLength)
	
	return baseTime + interferenceTime + (1.0-capacityFactor)*float64(remainingAnts)*0.5
}

// calculateOptimalAntsForPath calculates how many ants should use a path
func calculateOptimalAntsForPath(state PathState, remainingAnts int, targetTime float64) int {
	pathLength := float64(len(state.path) - 1)
	maxAnts := int(math.Ceil(targetTime / (pathLength * state.efficiency)))
	return min(maxAnts, remainingAnts)
}

// calculateAntDelay calculates how many turns an ant should wait before starting
func calculateAntDelay(state PathState, antNumber int) int {
	if antNumber == 0 {
		return 0
	}
	
	// Calculate delay based on path efficiency and current congestion
	baseDelay := int(float64(antNumber) * (1.0 - state.efficiency) * 1.5)
	return min(baseDelay, len(state.path)-2)
}

// simulateAntMovements simulates the actual movement of ants
func simulateAntMovements(ants []*Ant, c *colony.Colony, pathStates []PathState) []string {
	var moves []string
	turn := 0
	
	for {
		turnMoves := make([]string, 0)
		moveMade := false
		roomOccupancy := make(map[string]bool)

		// Sort ants by priority
		sortAntsByPriority(ants, c.End, turn)

		// Try to move each ant
		for _, ant := range ants {
			if ant.Position == len(ant.Path)-1 {
				continue // Ant has reached the end
			}

			// Check if ant should start moving yet
			if ant.Position == -1 && ant.Delay > turn {
				continue
			}

			nextRoom := getNextRoom(ant)
			if !canMoveToRoom(nextRoom, roomOccupancy, c) {
				continue
			}

			// Make the move
			if nextRoom != c.Start {
				turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", ant.ID, nextRoom))
			}

			updateAntPosition(ant, nextRoom, roomOccupancy)
			moveMade = true
		}

		if !moveMade {
			break
		}

		if len(turnMoves) > 0 {
			moves = append(moves, strings.Join(turnMoves, " "))
		}
		
		turn++
	}

	return moves
}

// sortAntsByPriority sorts ants based on their priority for movement
func sortAntsByPriority(ants []*Ant, endRoom string, currentTurn int) {
	sort.Slice(ants, func(i, j int) bool {
		return getAntPriority(ants[i], currentTurn) > getAntPriority(ants[j], currentTurn)
	})
}

// getAntPriority calculates movement priority for an ant
func getAntPriority(ant *Ant, currentTurn int) float64 {
	if ant.Position == len(ant.Path)-1 {
		return -1 // Already at end
	}

	// Calculate base priority based on progress
	progress := float64(ant.Position+1) / float64(len(ant.Path))
	
	// Prioritize ants that have already started moving
	if ant.Position > -1 {
		progress += 0.5
	}

	// Consider waiting time for ants that haven't started
	if ant.Position == -1 && ant.Delay <= currentTurn {
		progress += 0.3
	}

	return progress
}

// getNextRoom returns the next room for an ant
func getNextRoom(ant *Ant) string {
	if ant.Position == -1 {
		return ant.Path[0]
	}
	return ant.Path[ant.Position+1]
}

// canMoveToRoom checks if an ant can move to the specified room
func canMoveToRoom(room string, occupancy map[string]bool, c *colony.Colony) bool {
	return room == c.Start || room == c.End || !occupancy[room]
}

// updateAntPosition updates the ant's position and marks room occupancy
func updateAntPosition(ant *Ant, nextRoom string, occupancy map[string]bool) {
	if ant.Position == -1 {
		ant.Position = 0
	} else {
		ant.Position++
	}
	occupancy[nextRoom] = true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
