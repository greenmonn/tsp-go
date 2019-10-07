package graph

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path"

	"github.com/greenmonn/tsp-go/utils"
)

type Node struct {
	ID        int
	X         float64
	Y         float64
	Connected []*Node
	Degree    int
}

func NewNode(ID int, X float64, Y float64) *Node {
	return &Node{ID: ID,
		X:         X,
		Y:         Y,
		Connected: []*Node{},
		Degree:    0}
}

func (self *Node) Distance(other *Node) float64 {
	return math.Sqrt(math.Pow(self.X-other.X, 2) + math.Pow(self.Y-other.Y, 2))
}

func ReadNodesFromFile(filepath string) []*Node {
	file, err := os.Open(path.Join(utils.GetSourceRootPath(), filepath))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	nodes := make([]*Node, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id, xPos, yPos := utils.ParseLine(scanner.Text())
		if id == -1 {
			continue
		}

		nodes = append(nodes, NewNode(id, xPos, yPos))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return nodes
}

func WritePathToFile(route []*Node, name string) int {
	file, err := os.Create(path.Join(utils.GetSourceRootPath(), name+".csv"))
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)
	bytesWrote := 0

	for _, node := range route {
		n, err := writer.WriteString(fmt.Sprintf("%d\n", node.ID))

		if err != nil {
			log.Fatal(err)
		}
		writer.Flush()

		bytesWrote += n
	}

	return bytesWrote
}
