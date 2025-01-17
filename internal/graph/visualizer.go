package graph

import (
	"bytes"
	"fmt"
	"strings"
)

// Visualizer handles trust graph visualization
type Visualizer struct {
	nodes map[string]bool
	edges []Edge
}

// Edge represents a connection in the graph
type Edge struct {
	From      string
	To        string
	Predicate string
}

// NewVisualizer creates a new graph visualizer
func NewVisualizer() *Visualizer {
	return &Visualizer{
		nodes: make(map[string]bool),
		edges: []Edge{},
	}
}

// AddEdge adds a new edge to the graph
func (v *Visualizer) AddEdge(from, to, predicate string) {
	v.nodes[from] = true
	v.nodes[to] = true
	v.edges = append(v.edges, Edge{
		From:      from,
		To:        to,
		Predicate: predicate,
	})
}

// GenerateDOT returns a DOT format representation of the graph
func (v *Visualizer) GenerateDOT() string {
	var buf bytes.Buffer

	buf.WriteString("digraph TrustGraph {\n")
	buf.WriteString("  rankdir=LR;\n")
	buf.WriteString("  node [shape=box, style=rounded];\n")

	// Add nodes
	for node := range v.nodes {
		buf.WriteString(fmt.Sprintf("  %q;\n", node))
	}

	// Add edges
	for _, edge := range v.edges {
		buf.WriteString(fmt.Sprintf("  %q -> %q [label=%q];\n",
			edge.From, edge.To, edge.Predicate))
	}

	buf.WriteString("}\n")
	return buf.String()
}

// GenerateASCII returns an ASCII art representation of the graph
func (v *Visualizer) GenerateASCII() string {
	var buf bytes.Buffer

	// Sort nodes for consistent output
	nodeList := make([]string, 0, len(v.nodes))
	for node := range v.nodes {
		nodeList = append(nodeList, node)
	}

	for _, edge := range v.edges {
		buf.WriteString(fmt.Sprintf("%s --%s--> %s\n",
			padRight(edge.From, 20),
			padCenter(edge.Predicate, 10),
			edge.To))
	}

	return buf.String()
}

func padRight(s string, length int) string {
	if len(s) >= length {
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

func padCenter(s string, length int) string {
	if len(s) >= length {
		return s[:length]
	}
	padding := length - len(s)
	leftPad := padding / 2
	rightPad := padding - leftPad
	return strings.Repeat(" ", leftPad) + s + strings.Repeat(" ", rightPad)
} 