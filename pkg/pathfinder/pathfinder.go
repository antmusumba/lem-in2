package pathfinder

import "lem-in/pkg/colony"

// Path represents a sequence of room names
type Path []string

// FindPaths finds all possible paths from start to end using DFS
func FindPaths(c *colony.Colony) []Path {
	visited := make(map[string]bool)
	var paths []Path
	dfs(c, c.Start, visited, Path{c.Start}, &paths)
	
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

func dfs(c *colony.Colony, current string, visited map[string]bool, path Path, paths *[]Path) {
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
			dfs(c, next, visited, append(path, next), paths)
		}
	}
}
