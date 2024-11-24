package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lem-in/pkg/colony"
)

// ParseInput reads and parses the input file into a Colony structure
func ParseInput(filename string) (*colony.Colony, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ERROR: invalid data format")
	}
	defer file.Close()

	c := colony.NewColony()
	scanner := bufio.NewScanner(file)
	
	if err := parseAnts(scanner, c); err != nil {
		return nil, err
	}

	if err := parseRoomsAndTunnels(scanner, c); err != nil {
		return nil, err
	}

	if c.Start == "" || c.End == "" {
		return nil, fmt.Errorf("ERROR: invalid data format")
	}

	return c, nil
}

func parseAnts(scanner *bufio.Scanner, c *colony.Colony) error {
	if !scanner.Scan() {
		return fmt.Errorf("ERROR: invalid data format")
	}
	
	antCount, err := strconv.Atoi(scanner.Text())
	if err != nil || antCount <= 0 {
		return fmt.Errorf("ERROR: invalid data format")
	}
	
	c.NumAnts = antCount
	c.Input = append(c.Input, scanner.Text())
	return nil
}

func parseRoomsAndTunnels(scanner *bufio.Scanner, c *colony.Colony) error {
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

		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "-") {
			if err := parseTunnel(line, c); err != nil {
				return err
			}
			continue
		}

		if err := parseRoom(line, expectingStart, expectingEnd, c); err != nil {
			return err
		}

		if expectingStart {
			expectingStart = false
		} else if expectingEnd {
			expectingEnd = false
		}
	}

	return nil
}

func parseRoom(line string, isStart, isEnd bool, c *colony.Colony) error {
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

	room := &colony.Room{
		Name:    name,
		X:       x,
		Y:       y,
		IsStart: isStart,
		IsEnd:   isEnd,
	}

	if isStart {
		c.Start = name
	}
	if isEnd {
		c.End = name
	}

	c.Rooms[name] = room
	return nil
}

func parseTunnel(line string, c *colony.Colony) error {
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

	c.Tunnels = append(c.Tunnels, colony.Tunnel{From: parts[0], To: parts[1]})
	return nil
}
