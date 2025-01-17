package graph

import (
	"github.com/sirupsen/logrus"
	"axia/internal/crypto"
	"axia/internal/state"
)

type Node struct {
	ID          string
	Data        interface{}
	State       *state.StateManager
	Proof       *crypto.Proof
	logger      *logrus.Logger
}

type Edge struct {
	From        *Node
	To          *Node
	Weight      float64
	State       *state.StateManager
	Proof       *crypto.Proof
	logger      *logrus.Logger
}

type Graph struct {
	Nodes       []*Node
	Edges       []*Edge
	proofGen    *crypto.ProofGenerator
	logger      *logrus.Logger
}

func NewGraph(logger *logrus.Logger) *Graph {
	return &Graph{
		Nodes:    make([]*Node, 0),
		Edges:    make([]*Edge, 0),
		proofGen: crypto.NewProofGenerator(),
		logger:   logger,
	}
}

func (g *Graph) AddNode(data interface{}) (*Node, error) {
	g.logger.WithField("data", data).Info("Adding new node")
	
	node := &Node{
		ID:     uuid.New().String(),
		Data:   data,
		State:  state.NewNodeStateManager(g.logger),
		logger: g.logger,
	}
	
	proof, err := g.proofGen.GenerateProof(node)
	if err != nil {
		g.logger.WithError(err).Error("Failed to generate proof for node")
		return nil, err
	}
	
	node.Proof = proof
	g.Nodes = append(g.Nodes, node)
	
	g.logger.WithFields(logrus.Fields{
		"node_id": node.ID,
		"proof":   proof.Hash,
	}).Info("Node added successfully")
	
	return node, nil
} 