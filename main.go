package main

import (
	"fmt"
	"os"
	"strings"

	"lem2/utils"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . input.txt")
		return
	}
	file := os.Args[1]
	mySlices, _ := utils.ReadInput(file)

	// Variables to store start, end nodes, and edges
	var startNode, endNode string
	edges := make([]string, 0)
	readingNodes := true

	// Separate nodes and edges
	for _, line := range mySlices {
		line = strings.TrimSpace(line)

		// Skip the start and end markers and set the appropriate flags
		if line == "##start" {
			// Set reading flag to true after encountering ##start
			readingNodes = true
			continue
		} else if line == "##end" {
			// Set reading flag to false after encountering ##end
			readingNodes = false
			continue
		}

		// If we're reading nodes, set start and end nodes based on markers
		if readingNodes {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				if startNode == "" {
					startNode = parts[0] // Set the start node (first encountered after ##start)
				}
				endNode = parts[0] // Continuously update the end node (last encountered after ##end)
			}
		} else {
			// After ##end, collect edges
			if len(line) > 0 && strings.Contains(line, "-") {
				edges = append(edges, line)
			}
		}
	}

	// Check if start and end nodes were found
	if startNode == "" || endNode == "" {
		fmt.Println("Start or end node is missing.")
		return
	}

	// Create a new graph
	graph := NewGraph()

	// Add edges to the graph
	for _, edge := range edges {
		parts := strings.Split(edge, "-")
		if len(parts) == 2 {
			// Add edges to the graph (undirected)
			graph.AddEdge(parts[0], parts[1])
		}
	}

	// Print the graph to check adjacency list
	fmt.Println("Graph Adjacency List:")
	for node, neighbors := range graph.vertices {
		fmt.Printf("%s: %v\n", node, neighbors)
	}

	// Perform BFS starting from the start node
	fmt.Println("\nStarting BFS from node:", startNode)
	graph.BFS(startNode)
}

// Graph represents a graph using an adjacency list with string keys and values.
type Graph struct {
	vertices map[string][]string
}

// NewGraph creates a new graph.
func NewGraph() *Graph {
	return &Graph{vertices: make(map[string][]string)}
}

// AddEdge adds an edge between two nodes.
func (g *Graph) AddEdge(start, end string) {
	g.vertices[start] = append(g.vertices[start], end)
	g.vertices[end] = append(g.vertices[end], start) // For undirected graph
}

// BFS performs a breadth-first search starting from the given node.
func (g *Graph) BFS(start string) {
	visited := make(map[string]bool)
	queue := []string{start}

	visited[start] = true

	fmt.Println("BFS starting from node:", start)

	// Iterate through the queue
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:] // Dequeue the node

		fmt.Printf("Visited: %s\n", node)

		// Enqueue all unvisited neighbors
		for _, neighbor := range g.vertices[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
}
