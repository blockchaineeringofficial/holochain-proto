package holochain

import (
	. "github.com/holochain/holochain-proto/hash"
	peer "github.com/libp2p/go-libp2p-peer"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDelName(t *testing.T) {
	Convey("delete action should have the right name", t, func() {
		// https://github.com/holochain/holochain-proto/issues/715
		// a := NewDelAction(DelEntry{Hash: ""})
		a := ActionDel{entry: DelEntry{Hash: ""}}
		So(a.Name(), ShouldEqual, "del")
	})
}

func TestAPIFnDelName(t *testing.T) {
	Convey("delete action function should have the right name", t, func() {
		// https://github.com/holochain/holochain-proto/issues/715
		// a := NewDelAction(DelEntry{Hash: ""})
		a := ActionDel{entry: DelEntry{Hash: ""}}
		fn := &APIFnDel{action: a}
		So(fn.Name(), ShouldEqual, "del")
	})
}

func TestActionDelete(t *testing.T) {
	nodesCount := 3
	mt := setupMultiNodeTesting(nodesCount)
	defer mt.cleanupMultiNodeTesting()

	h := mt.nodes[0]
	ringConnect(t, mt.ctx, mt.nodes, nodesCount)

	profileHash := commit(h, "profile", `{"firstName":"Zippy","lastName":"Pinhead"}`)
	entry := DelEntry{Hash: profileHash, Message: "expired"}
	a := &ActionDel{entry: entry}
	response, err := h.commitAndShare(a, NullHash())
	if err != nil {
		panic(err)
	}
	deleteHash := response.(Hash)

	Convey("when deleting a hash the del entry itself should be published to the DHT", t, func() {
		req := GetReq{H: deleteHash, GetMask: GetMaskEntry}
		_, err := callGet(h, req, &GetOptions{GetMask: req.GetMask})
		So(err, ShouldBeNil)

		h2 := mt.nodes[2]
		_, err = callGet(h2, req, &GetOptions{GetMask: req.GetMask})
		So(err, ShouldBeNil)
	})
}

func TestSysValidateDel(t *testing.T) {
	d, _, h := PrepareTestChain("test")
	defer CleanupTestChain(h, d)

	hash := commit(h, "evenNumbers", "2")
	//	_, def, _ := h.GetEntryDef("evenNumbers")

	Convey("it should check that entry isn't linking ", t, func() {
		a := NewDelAction(DelEntry{Hash: hash})
		_, ratingsDef, _ := h.GetEntryDef("rating")
		err := a.SysValidation(h, ratingsDef, nil, []peer.ID{h.nodeID})
		So(err, ShouldBeError)
	})
}

func TestSysDel(t *testing.T) {
	d, _, h := PrepareTestChain("test")
	defer CleanupTestChain(h, d)
	var err error

	Convey("deleting should fail for all sys entry types except delete", t, func() {
		a := NewDelAction(DelEntry{})
		_, err = h.ValidateAction(a, DNAEntryType, nil, []peer.ID{h.nodeID})
		So(err, ShouldEqual, ErrEntryDefInvalid)

		_, err = h.ValidateAction(a, KeyEntryType, nil, []peer.ID{h.nodeID})
		So(err, ShouldEqual, ErrEntryDefInvalid)

		_, err = h.ValidateAction(a, AgentEntryType, nil, []peer.ID{h.nodeID})
		So(err, ShouldEqual, ErrEntryDefInvalid)

		_, err = h.ValidateAction(a, HeadersEntryType, nil, []peer.ID{h.nodeID})
		So(err, ShouldEqual, ErrEntryDefInvalid)
	})
}
