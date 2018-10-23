package core

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/textileio/textile-go/pb"
	"github.com/textileio/textile-go/repo"
	mh "gx/ipfs/QmPnFwZ2JXKnXgMw8CdBPxn7FWh6LLdjUjxV1fKHuJnkr8/go-multihash"
	"gx/ipfs/QmdVrMn1LhB4ybb8hMVaMLXnA8XRSewMnK6YqXKXoTcRvN/go-libp2p-peer"
)

// joinInitial creates an outgoing join block for an emtpy thread
func (t *Thread) joinInitial() (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	msg, err := t.buildJoin(t.node().Identity.Pretty())
	if err != nil {
		return nil, err
	}

	// commit to ipfs
	res, err := t.commitBlock(msg, pb.ThreadBlock_JOIN, nil)
	if err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(res, repo.JoinBlock, nil); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	log.Debugf("added JOIN to %s: %s", t.Id, res.hash.B58String())

	// all done
	return res.hash, nil
}

// join creates an outgoing join block
func (t *Thread) join(inviterId peer.ID) (mh.Multihash, error) {
	t.mux.Lock()
	defer t.mux.Unlock()

	// build block
	msg, err := t.buildJoin(inviterId.Pretty())
	if err != nil {
		return nil, err
	}

	// commit to ipfs
	res, err := t.commitBlock(msg, pb.ThreadBlock_JOIN, nil)
	if err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(res, repo.JoinBlock, nil); err != nil {
		return nil, err
	}

	// update head
	if err := t.updateHead(res.hash); err != nil {
		return nil, err
	}

	// add new peer
	if inviterId.Pretty() != t.node().Identity.Pretty() {
		newPeer := &repo.ThreadPeer{
			Id:       inviterId.Pretty(),
			ThreadId: t.Id,
		}
		if err := t.datastore.ThreadPeers().Add(newPeer); err != nil {
			log.Errorf("error adding peer: %s", err)
		}
	}

	// post it
	if err := t.post(res, t.Peers()); err != nil {
		return nil, err
	}

	log.Debugf("added JOIN to %s: %s", t.Id, res.hash.B58String())

	// all done
	return res.hash, nil
}

// handleJoinBlock handles an incoming join block
func (t *Thread) handleJoinBlock(hash mh.Multihash, block *pb.ThreadBlock) (*pb.ThreadJoin, error) {
	msg := new(pb.ThreadJoin)
	if err := ptypes.UnmarshalAny(block.Payload, msg); err != nil {
		return nil, err
	}

	// index it locally
	if err := t.indexBlock(&commitResult{hash: hash, header: block.Header}, repo.JoinBlock, nil); err != nil {
		return nil, err
	}

	// collect author as an unwelcomed peer
	pid, err := peer.IDB58Decode(block.Header.Author)
	if err != nil {
		return nil, err
	}
	t.addOrUpdatePeer(pid, msg.Username, msg.Inboxes)

	return msg, nil
}

// buildJoin builds up a join block
func (t *Thread) buildJoin(inviterId string) (*pb.ThreadJoin, error) {
	msg := &pb.ThreadJoin{
		Inviter: inviterId,
	}
	username, err := t.datastore.Profile().GetUsername()
	if err != nil {
		return nil, err
	}
	if username != nil {
		msg.Username = *username
	}
	for _, ses := range t.datastore.CafeSessions().List() {
		msg.Inboxes = append(msg.Inboxes, ses.CafeId)
	}
	return msg, nil
}
