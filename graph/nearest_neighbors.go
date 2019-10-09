package graph

import "math"

func SetNearestNeighbors(K int) {
	graph.NearestNeighbors = nearestNeighbors(K)
}

func NearestNeighbors(id int) []*Node {
	return graph.NearestNeighbors[id]
}

func nearestNeighbors(K int) map[int][]*Node {
	// Save nearest K neighbors in advance
	N := GetNodesCount()
	nodes := GetNodes()

	nearests := make(map[int][]*Node)

	for i := 0; i < N; i++ {
		nearests[nodes[i].ID] = findNearests(nodes, K, i)
	}

	return nearests
}

func findNearests(nodes []*Node, count int, i int) []*Node {
	nearests := make([]*Node, count)
	values := make([]float64, count)

	for i := 0; i < count; i++ {
		values[i] = math.MaxFloat64
	}

	N := GetNodesCount()
	D := GetDistanceByIndex

	for j := 0; j < N; j++ {
		distance := D(i, j)
		for k := 0; k < count; k++ {
			if distance < values[k] {
				for rank := count - 1; rank > k; rank-- {
					nearests[rank] = nearests[rank-1]
					values[rank] = values[rank-1]
				}

				nearests[k] = nodes[j]
				values[k] = distance

				break
			}
		}
	}

	return nearests
}
