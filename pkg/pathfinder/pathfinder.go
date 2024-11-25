package pathfinder

import (
	"lem-in/pkg/colony"
	"math"
	"sort"
)

// Path represents a sequence of room names
type Path []string

// PathScore represents a path with its efficiency score
type PathScore struct {
	path       Path
	score      float64
	bottleneck float64
}

// FindPaths finds all possible paths from start to end using DFS and optimizes them
func FindPaths(c *colony.Colony) []Path {
	// First, find all possible paths
	visited := make(map[string]bool)
	var paths []Path
	dfs(c, c.Start, visited, Path{c.Start}, &paths)

	// Then optimize them based on the number of ants
	return optimizePaths(paths, c)
}

// calculatePathScore computes a comprehensive efficiency score for a path
func calculatePathScore(path Path, c *colony.Colony, otherPaths []Path) float64 {
	pathLen := float64(len(path))
	
	// Base score inversely proportional to path length
	baseScore := 100.0 / pathLen

	// Calculate path independence score (how unique this path is)
	independenceScore := calculateIndependenceScore(path, otherPaths)
	
	// Calculate bottleneck score (how many shared rooms with other paths)
	bottleneckScore := calculateBottleneckScore(path, c)
	
	// Calculate position score (how well-positioned the path is relative to start/end)
	positionScore := calculatePositionScore(path, c)

	// Weighted combination of all scores
	return (baseScore * 0.4) + (independenceScore * 0.3) + (bottleneckScore * 0.2) + (positionScore * 0.1)
}

// calculateIndependenceScore measures how independent this path is from others
func calculateIndependenceScore(path Path, otherPaths []Path) float64 {
	if len(otherPaths) == 0 {
		return 100.0
	}

	totalOverlap := 0
	for _, other := range otherPaths {
		if other[0] == path[0] && other[len(other)-1] == path[len(path)-1] {
			continue // Skip comparing with itself
		}
		
		// Count shared rooms (excluding start/end)
		sharedRooms := make(map[string]bool)
		for _, room := range path[1:len(path)-1] {
			sharedRooms[room] = true
		}
		
		overlap := 0
		for _, room := range other[1:len(other)-1] {
			if sharedRooms[room] {
				overlap++
			}
		}
		totalOverlap += overlap
	}

	// Higher score for less overlap
	return 100.0 / (1.0 + float64(totalOverlap))
}

// calculateBottleneckScore evaluates potential bottlenecks in the path
func calculateBottleneckScore(path Path, c *colony.Colony) float64 {
	if len(path) <= 2 {
		return 100.0 // Direct path
	}

	// Count connections for each room in the path
	bottleneckFactor := 0.0
	for _, room := range path[1:len(path)-1] {
		connections := 0
		for _, tunnel := range c.Tunnels {
			if tunnel.From == room || tunnel.To == room {
				connections++
			}
		}
		// More connections = less bottleneck
		bottleneckFactor += float64(connections)
	}

	return (bottleneckFactor / float64(len(path)-2)) * 20.0 // Scale to 0-100
}

// calculatePositionScore evaluates the path's position relative to start/end
func calculatePositionScore(path Path, c *colony.Colony) float64 {
	// Calculate average distance from optimal straight line
	startRoom := c.Rooms[c.Start]
	endRoom := c.Rooms[c.End]
	
	// Calculate ideal straight line
	dx := float64(endRoom.X - startRoom.X)
	dy := float64(endRoom.Y - startRoom.Y)
	length := math.Sqrt(dx*dx + dy*dy)
	
	if length == 0 {
		return 100.0
	}

	// Calculate average deviation from straight line
	totalDeviation := 0.0
	for _, roomName := range path[1:len(path)-1] {
		room := c.Rooms[roomName]
		
		// Calculate distance from point to line
		deviation := math.Abs(float64(room.X-startRoom.X)*dy - float64(room.Y-startRoom.Y)*dx) / length
		totalDeviation += deviation
	}

	avgDeviation := totalDeviation / float64(len(path)-2)
	return 100.0 / (1.0 + avgDeviation)
}

// optimizePaths optimizes path selection based on comprehensive scoring
func optimizePaths(paths []Path, c *colony.Colony) []Path {
	if len(paths) == 0 {
		return paths
	}

	// Calculate initial scores for all paths
	pathScores := make([]PathScore, len(paths))
	for i, path := range paths {
		otherPaths := append(paths[:i], paths[i+1:]...)
		pathScores[i] = PathScore{
			path:       path,
			score:      calculatePathScore(path, c, otherPaths),
			bottleneck: calculateBottleneckScore(path, c),
		}
	}

	// Sort by score in descending order
	sort.Slice(pathScores, func(i, j int) bool {
		return pathScores[i].score > pathScores[j].score
	})

	// Select optimal combination of paths
	var optimized []Path
	numAnts := c.NumAnts
	targetPaths := int(math.Sqrt(float64(numAnts))) + 1 // Dynamic path count based on ant count
	
	for i, ps := range pathScores {
		if i >= targetPaths && len(optimized) >= 2 {
			break
		}

		// Check if this path adds value to our selection
		isUseful := true
		totalOverlap := 0
		for _, existingPath := range optimized {
			overlap := countSharedRooms(ps.path, existingPath)
			totalOverlap += overlap
			
			// Skip if too much overlap with existing paths
			maxAllowedOverlap := (len(ps.path) + len(existingPath)) / 8
			if overlap > maxAllowedOverlap {
				isUseful = false
				break
			}
		}

		// Add path if it's useful and doesn't create too many bottlenecks
		if isUseful && (totalOverlap <= len(optimized) || len(optimized) < 2) {
			optimized = append(optimized, ps.path)
		}
	}

	// Ensure we have at least one path
	if len(optimized) == 0 {
		optimized = append(optimized, pathScores[0].path)
	}

	return optimized
}

// countSharedRooms counts rooms shared between two paths (excluding start/end)
func countSharedRooms(path1, path2 Path) int {
	shared := make(map[string]bool)
	for _, room := range path1[1:len(path1)-1] {
		shared[room] = true
	}

	count := 0
	for _, room := range path2[1:len(path2)-1] {
		if shared[room] {
			count++
		}
	}
	return count
}

func dfs(c *colony.Colony, current string, visited map[string]bool, path Path, paths *[]Path) {
	if current == c.End {
		pathCopy := make(Path, len(path))
		copy(pathCopy, path)
		*paths = append(*paths, pathCopy)
		return
	}

	visited[current] = true
	defer delete(visited, current)

	// Get and sort next rooms by their potential
	nextRooms := getNextRooms(c, current, visited)
	sortRoomsByPotential(nextRooms, c, c.End)

	for _, next := range nextRooms {
		dfs(c, next, visited, append(path, next), paths)
	}
}

// getNextRooms gets all possible next rooms
func getNextRooms(c *colony.Colony, current string, visited map[string]bool) []string {
	var nextRooms []string
	for _, tunnel := range c.Tunnels {
		var next string
		if tunnel.From == current {
			next = tunnel.To
		} else if tunnel.To == current {
			next = tunnel.From
		}

		if next != "" && !visited[next] {
			nextRooms = append(nextRooms, next)
		}
	}
	return nextRooms
}

// sortRoomsByPotential sorts rooms by their potential for reaching the end
func sortRoomsByPotential(rooms []string, c *colony.Colony, end string) {
	endRoom := c.Rooms[end]
	sort.Slice(rooms, func(i, j int) bool {
		roomI := c.Rooms[rooms[i]]
		roomJ := c.Rooms[rooms[j]]
		
		// Calculate Manhattan distance
		distI := abs(roomI.X-endRoom.X) + abs(roomI.Y-endRoom.Y)
		distJ := abs(roomJ.X-endRoom.X) + abs(roomJ.Y-endRoom.Y)
		
		// Also consider number of connections
		connectionsI := countConnections(c, rooms[i])
		connectionsJ := countConnections(c, rooms[j])
		
		// Weighted score combining distance and connections
		scoreI := float64(distI) / (1.0 + float64(connectionsI))
		scoreJ := float64(distJ) / (1.0 + float64(connectionsJ))
		
		return scoreI < scoreJ
	})
}

// countConnections counts the number of connections a room has
func countConnections(c *colony.Colony, room string) int {
	count := 0
	for _, tunnel := range c.Tunnels {
		if tunnel.From == room || tunnel.To == room {
			count++
		}
	}
	return count
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
