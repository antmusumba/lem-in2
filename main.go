package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lem-in/pkg/parser"
	"lem-in/pkg/simulator"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: invalid data format")
		return
	}

	colony, err := parser.ParseInput(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	moves := simulator.SimulateMovement(colony)
	for _, move := range moves {
		fmt.Println(move)
	}
}

type Room struct {
	Name    string
	X       int
	Y       int
	IsStart bool
	IsEnd   bool
}

type Tunnel struct {
	From string
	To   string
}

type Colony struct {
	NumAnts int
	Rooms   map[string]*Room
	Tunnels []Tunnel
	Start   string
	End     string
	Input   []string
}

func newColony() *Colony {
	return &Colony{
		Rooms:   make(map[string]*Room),
		Tunnels: make([]Tunnel, 0),
		Input:   make([]string, 0),
	}
}

func (c *Colony) parseInput(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("ERROR: invalid data format")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// First line must be number of ants
	if !scanner.Scan() {
		return fmt.Errorf("ERROR: invalid data format")
	}

	antCount, err := strconv.Atoi(scanner.Text())
	if err != nil || antCount <= 0 {
		return fmt.Errorf("ERROR: invalid data format")
	}
	c.NumAnts = antCount
	c.Input = append(c.Input, scanner.Text())

	var expectingStart, expectingEnd bool

	for scanner.Scan() {
		line := scanner.Text()
		c.Input = append(c.Input, line)

		if line == "" {
			continue
		}

		if line == "##start" {
			expectingStart = true
			continue
		}
		if line == "##end" {
			expectingEnd = true
			continue
		}

		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Check if it's a tunnel definition
		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return fmt.Errorf("ERROR: invalid data format")
			}

			if _, exists := c.Rooms[parts[0]]; !exists {
				return fmt.Errorf("ERROR: invalid data format")
			}
			if _, exists := c.Rooms[parts[1]]; !exists {
				return fmt.Errorf("ERROR: invalid data format")
			}

			c.Tunnels = append(c.Tunnels, Tunnel{From: parts[0], To: parts[1]})
			continue
		}

		// Must be a room definition
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return fmt.Errorf("ERROR: invalid data format")
		}

		name := parts[0]
		if strings.HasPrefix(name, "L") || strings.HasPrefix(name, "#") {
			return fmt.Errorf("ERROR: invalid data format")
		}

		x, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("ERROR: invalid data format")
		}

		y, err := strconv.Atoi(parts[2])
		if err != nil {
			return fmt.Errorf("ERROR: invalid data format")
		}

		room := &Room{
			Name:    name,
			X:       x,
			Y:       y,
			IsStart: expectingStart,
			IsEnd:   expectingEnd,
		}

		if expectingStart {
			c.Start = name
			expectingStart = false
		} else if expectingEnd {
			c.End = name
			expectingEnd = false
		}

		c.Rooms[name] = room
	}

	if c.Start == "" {
		return fmt.Errorf("ERROR: invalid data format")
	}
	if c.End == "" {
		return fmt.Errorf("ERROR: invalid data format")
	}

	return nil
}

type Path []string

func (c *Colony) findPaths() []Path {
	visited := make(map[string]bool)
	var paths []Path
	c.dfs(c.Start, visited, Path{c.Start}, &paths)

	// Sort paths by length
	for i := 0; i < len(paths)-1; i++ {
		for j := i + 1; j < len(paths); j++ {
			if len(paths[i]) > len(paths[j]) {
				paths[i], paths[j] = paths[j], paths[i]
			}
		}
	}

	return paths
}

func (c *Colony) dfs(current string, visited map[string]bool, path Path, paths *[]Path) {
	if current == c.End {
		pathCopy := make(Path, len(path))
		copy(pathCopy, path)
		*paths = append(*paths, pathCopy)
		return
	}

	visited[current] = true
	defer delete(visited, current)

	for _, tunnel := range c.Tunnels {
		next := ""
		if tunnel.From == current {
			next = tunnel.To
		} else if tunnel.To == current {
			next = tunnel.From
		}

		if next != "" && !visited[next] {
			c.dfs(next, visited, append(path, next), paths)
		}
	}
}

type Ant struct {
	ID       int
	Path     Path
	Position int
}

func (c *Colony) simulateAntMovement() []string {
	paths := c.findPaths()
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
