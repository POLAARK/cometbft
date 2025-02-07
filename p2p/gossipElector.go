package p2p

import (
	"math/rand"
	"os"
	"strconv"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/p2p/nodekey"
	"github.com/cometbft/cometbft/types"
)

// PeerWithStake represents a peer and its associated stake.
type peerWithStake struct {
	peer  Peer
	stake int64
}

// GossipElector manages peer selection for gossiping transactions.
type GossipElector struct {
	validatorSet *types.ValidatorSet
	peers        *PeerSet // Ensure peers are accessible
	logger log.Logger
}

// NewGossipElector initializes a GossipElector with a validator set and peer set.
func NewGossipElector(vs *types.ValidatorSet, ps *PeerSet) *GossipElector {
	return &GossipElector{
		validatorSet: vs,
		peers:        ps,
		logger: log.NewLogger(os.Stdout),
	}
}

// SelectPeersForPropagation selects peers for gossiping transactions based on stake.
// SelectPeersForPropagation selects peers for gossiping transactions based on stake.
func (g *GossipElector) SelectPeersForPropagation() []Peer {
	// Load threshold percentage from environment variable.
	thresholdStr := os.Getenv("PROPAGATION_STAKE_THRESHOLD")
	if thresholdStr == "" {
		thresholdStr = "30" // Default to 30% if not set
	}
	thresholdPerc, err := strconv.Atoi(thresholdStr)
	if err != nil {
		thresholdPerc = 30
	}

	// Retrieve total voting power from validator set.
	totalVotingPower := g.validatorSet.TotalVotingPower()
	thresholdAbsolute := int64(totalVotingPower * int64(thresholdPerc) / 100)

	// Build a list pairing each peer with its stake.
	var peerStakes []peerWithStake

	// Collect eligible peers based on validator set.
	for _, v := range g.validatorSet.Validators {
		id := nodekey.PubKeyToID(v.PubKey)

		// Log the ID of the current validator and its voting power

		// Retrieve peer from the peer set.
		peer := g.peers.lookup[id]
		g.logger.Info("Peer matched !", "peer", peer)
		if peer != nil {
			// Log when a match is found

			peerStakes = append(peerStakes, peerWithStake{
				peer:  peer.peer,
				stake: v.VotingPower,
			})
		} else {
			// Log when no peer is found for the current validator
			g.logger.Info("No peer found for validator", "validatorID", id)
		}
	}

	// Log the total number of peers found after the loop
	g.logger.Info("Total peers found", "numPeers", len(peerStakes), "peers", g.peers)

	// Shuffle peers randomly.
	rand.Shuffle(len(peerStakes), func(i, j int) {
		peerStakes[i], peerStakes[j] = peerStakes[j], peerStakes[i]
	})

	// Select peers until cumulative stake exceeds threshold.
	var selected []Peer
	var cumulative int64
	for _, ps := range peerStakes {
		if cumulative >= thresholdAbsolute { // Corrected comparison
			break
		}
		selected = append(selected, ps.peer)
		cumulative += ps.stake
	}

	// Log selected peers for debugging
	g.logger.Info("selected peers", "peers", selected)
	return selected
}


func (g *GossipElector) AddPeerToSet(peer Peer) {
	g.peers.mtx.Lock()
	defer g.peers.mtx.Unlock()

	id := peer.ID()

	if _, exists := g.peers.lookup[id]; exists {
		return
	}

	g.peers.lookup[id] = &peerSetItem{peer: peer}
	g.peers.list = append(g.peers.list, peer)

}