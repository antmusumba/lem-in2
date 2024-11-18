package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Room represents a room in the ant farm.
type Room struct {
	Name string
	X    int
	Y    int
	Ants int // Number of ants currently in the room
}

// Tunnel represents a connection between two rooms.
type Tunnel struct {
	From *Room
	To   *Room
}

// Colony represents the entire ant farm structure.
type Colony struct {
	Rooms   map[string]*Room
	Tunnels []*Tunnel
	Start   *Room
	End     *Room
}

// readInput reads the input file and constructs the colony structure.
func readInput(filename string) (*Colony, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	colony := &Colony{
		Rooms:   make(map[string]*Room),
		Tunnels: make([]*Tunnel, 0),
	}

	var currentSection string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		if line == "##start" {
			currentSection = "start"
			continue
		} else if line == "##end" {
			currentSection = "end"
			continue
		}

		switch currentSection {
		case "start":
			room := parseRoom(line)
			if room != nil {
				colony.Start = room
				colony.Rooms[room.Name] = room
			}

		case "end":
			room := parseRoom(line)
			if room != nil {
				colony.End = room
				colony.Rooms[room.Name] = room
			}

		default:
			if strings.Contains(line, "-") { // Tunnel definition
				tunnelParts := strings.Split(line, "-")
				fromRoom := colony.Rooms[tunnelParts[0]]
				toRoom := colony.Rooms[tunnelParts[1]]
				if fromRoom == nil || toRoom == nil {
					return nil, fmt.Errorf("invalid tunnel between %s and %s", tunnelParts[0], tunnelParts[1])
				}
				colony.Tunnels = append(colony.Tunnels, &Tunnel{From: fromRoom, To: toRoom})
			} else { // Room definition for regular rooms (not start or end)
				room := parseRoom(line)
				if room != nil {
					colony.Rooms[room.Name] = room
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return colony, nil
}

// parseCoordinates parses room coordinates from a string.
func parseCoordinates(coords string) (int, int, error) {
	parts := strings.Fields(coords)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid coordinates format: %s", coords)
	}

	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid x-coordinate: %s", parts[0])
	}
	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid y-coordinate: %s", parts[1])
	}
	return x, y, nil
}

// parseRoom parses a room definition from a line.
func parseRoom(line string) *Room {
    parts := strings.Fields(line)
    if len(parts) < 3 { // Check for at least three parts (name x y)
        fmt.Printf("Invalid room format: %s (expected format: name x y)\n", line)
        return nil // Return nil if invalid format
    }

    name := parts[0] // This will be a numeric string like "1", "2", etc.
    xStr := parts[1]
    yStr := parts[2]

    x, y, _ := parseCoordinates(fmt.Sprintf("%s %s", xStr, yStr))
    return &Room{Name: name, X: x, Y: y}
}

// findShortestPath finds the shortest path from start to end using BFS.
func findShortestPath(colony *Colony) [][]*Room {
	var paths [][]*Room

	queue := [][]*Room{{colony.Start}}
	visited := make(map[string]bool)

	for len(queue) > 0 {
	    path := queue[0]
	    queue = queue[1:]

	    currentRoom := path[len(path)-1]
	    if currentRoom == colony.End {
	        paths = append(paths, path)
	        continue
	    }

	    visited[currentRoom.Name] = true

	    for _, tunnel := range colony.Tunnels {
	        if tunnel.From == currentRoom && !visited[tunnel.To.Name] {
	            newPath := append([]*Room{}, path...)
	            newPath = append(newPath, tunnel.To)
	            queue = append(queue, newPath)
	        }
	    }
    }

	return paths // Return all found paths.
}

// moveAnts moves ants along the shortest paths found.
func moveAnts(colony *Colony) []string {
    paths := findShortestPath(colony)

    // Check if any paths were found.
    if len(paths) == 0 || len(paths[0]) <= 1 { // Ensure there's at least one valid path with more than one room.
        return []string{"ERROR: No valid path found."}
    }

    moves := []string{}
    antCount := 3 // Example number of ants; this can be dynamically set.

    for i := 0; i < antCount; i++ {
        moves = append(moves, fmt.Sprintf("L%d-%s", i+1, paths[0][1].Name)) // Move each ant along the first valid path.
        for j := 2; j < len(paths[0]); j++ { 
            moves = append(moves, fmt.Sprintf("L%d-%s", i+1, paths[0][j].Name))
        }
    }

    return moves // Return all ant movements.
}

// printResults prints the results of the simulation.
func printResults(colony *Colony) {
	fmt.Println(len(colony.Rooms)) // Number of rooms

	for _, room := range colony.Rooms {
	    fmt.Printf("%s %d %d\n", room.Name, room.X, room.Y)
    }

	for _, tunnel := range colony.Tunnels {
	    fmt.Printf("%s-%s\n", tunnel.From.Name, tunnel.To.Name)
    }

	moves := moveAnts(colony)
	for _, move := range moves {
	    fmt.Println(move)
    }
}

func main() {
	if len(os.Args) < 2 {
	    fmt.Println("Usage: go run . <input_file>")
	    return
    }

	colony, err := readInput(os.Args[1])
	if err != nil {
	    fmt.Println("ERROR:", err)
	    return
    }

	printResults(colony)
}
