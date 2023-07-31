package main

import "math"

// Node represents a node in the A* algorithm.
type Node struct {
	X, Y       int
	F, G, H    float64
	ParentNode *Node
}

// AStar finds a path from the start point to the goal point on the grid using the A* algorithm.
func AStar(grid *Grid, startX, startY, goalX, goalY int) []*Node {
	startNode := &Node{X: startX, Y: startY}
	goalNode := &Node{X: goalX, Y: goalY}

	openList := []*Node{startNode}
	closedList := map[Node]bool{}
	directions := [4][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	// Heuristic function: Manhattan distance (can be changed to Euclidean distance if needed).
	heuristic := func(current, goal *Node) float64 {
		return math.Abs(float64(current.X-goal.X)) + math.Abs(float64(current.Y-goal.Y))
	}

	// Check if a node is in the closed list.
	inClosedList := func(node *Node) bool {
		_, ok := closedList[*node]
		return ok
	}

	// Check if a node is in the open list.
	inOpenList := func(node *Node) bool {
		for _, n := range openList {
			if n.X == node.X && n.Y == node.Y {
				return true
			}
		}
		return false
	}

	// Get the node with the lowest F value from the open list.
	getLowestFNode := func() *Node {
		lowestF := math.MaxFloat64
		var lowestNode *Node
		for _, node := range openList {
			if node.F < lowestF {
				lowestF = node.F
				lowestNode = node
			}
		}
		return lowestNode
	}

	for len(openList) > 0 {
		currentNode := getLowestFNode()

		if currentNode.X == goalNode.X && currentNode.Y == goalNode.Y {
			// Build and return the path.
			path := []*Node{}
			current := currentNode
			for current != nil {
				path = append([]*Node{current}, path...)
				current = current.ParentNode
			}
			return path
		}

		// Move the current node from the open list to the closed list.
		openList = removeNodeFromList(openList, currentNode)
		closedList[*currentNode] = true

		// Explore adjacent nodes.
		for _, dir := range directions {
			nextX := currentNode.X + dir[0]
			nextY := currentNode.Y + dir[1]

			// Check if the next position is valid and not in the closed list.
			if nextX >= 0 && nextX < int(grid.Width) && nextY >= 0 && nextY < int(grid.Height) {
				nextNode := &Node{X: nextX, Y: nextY, ParentNode: currentNode}

				if !grid.Data[nextY][nextX].Walkable || inClosedList(nextNode) {
					continue
				}

				// Calculate the G, H, and F values for the next node.
				nextNode.G = currentNode.G + 1
				nextNode.H = heuristic(nextNode, goalNode)
				nextNode.F = nextNode.G + nextNode.H

				if !inOpenList(nextNode) {
					openList = append(openList, nextNode)
				}
			}
		}
	}

	// If no path is found, return an empty path.
	return []*Node{}
}

// Remove a node from a list of nodes.
func removeNodeFromList(list []*Node, node *Node) []*Node {
	for i, n := range list {
		if n == node {
			list[i] = list[len(list)-1]
			return list[:len(list)-1]
		}
	}
	return list
}
