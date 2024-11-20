package main

import (
	"fmt"
	"sort"
)

type Graph struct {
	vertices map[string][]string
}

func NewGraph() *Graph {
	return &Graph{vertices: make(map[string][]string)}
}

func (g *Graph) AddEdge(start, end string) {
	g.vertices[start] = append(g.vertices[start], end)
	g.vertices[end] = append(g.vertices[end], start) // For undirected graph
}

// FindAllPaths finds all paths from start to end
func (g *Graph) FindAllPaths(start, end string) [][]string {
	var paths [][]string
	var dfs func(current string, visited map[string]bool, path []string)

	dfs = func(current string, visited map[string]bool, path []string) {
		if current == end {
			// Add the completed path
			paths = append(paths, append([]string{}, path...))
			return
		}

		visited[current] = true

		for _, neighbor := range g.vertices[current] {
			if !visited[neighbor] {
				dfs(neighbor, visited, append(path, neighbor))
			}
		}

		visited[current] = false
	}

	dfs(start, make(map[string]bool), []string{start})
	return paths
}

func SimulateAnts(paths [][]string, ants int) {
	// Sort paths by length (shortest first)
	sort.Slice(paths, func(i, j int) bool {
		return len(paths[i]) < len(paths[j])
	})

	// Distribute ants across paths
	assignments := make([][]int, len(paths))
	for i := 0; i < ants; i++ {
		assignments[i%len(paths)] = append(assignments[i%len(paths)], i+1)
	}

	// Simulate movement
	step := 0
	for {
		step++
		fmt.Printf("\nStep %d:\n", step)
		moving := false

		for i, path := range paths {
			for j, ant := range assignments[i] {
				pos := step - j - 1 // Calculate the position of the ant along the path

				// Check if the position is valid
				if pos >= 0 && pos < len(path) {
					moving = true
					fmt.Printf("Ant %d moves to %s\n", ant, path[pos])
				}
			}
		}

		// Stop if no ants are moving
		if !moving {
			break
		}
	}
}

func main() {
	graph := NewGraph()
	graph.AddEdge("1", "3")
	graph.AddEdge("1", "2")
	graph.AddEdge("3", "4")
	graph.AddEdge("2", "4")
	graph.AddEdge("4", "5")
	graph.AddEdge("5", "6")
	graph.AddEdge("6", "7")

	paths := graph.FindAllPaths("1", "7")
	fmt.Println("Paths from start to end:", paths)

	ants := 6
	SimulateAnts(paths, ants)
}
