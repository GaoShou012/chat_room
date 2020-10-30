package frontier

import (
	"testing"

	pegomock "github.com/petergtz/pegomock"
	pb "wchat.im/proto/message"
)

func TestAnchor(t *testing.T) {
	la := NewMockLeavable(pegomock.WithT(t))
	an := &Anchor{}

	an.SetAnchor(la, nil)
	if !an.IsAnchored(la) {
		t.Error("expect anchored")
	}
	if ok := an.DelAnchor(la); !ok {
		t.Error("expect del anchor ok")
	}
	la.VerifyWasCalledOnce().Leave(nil)
	if an.IsAnchored(la) {
		t.Error("expect not anchored")
	}

	la1 := NewMockLeavable(pegomock.WithT(t))
	la2 := NewMockLeavable(pegomock.WithT(t))
	an.SetAnchor(la1, 1)
	an.SetAnchor(la2, 2)
	an.DelAnchors()
	la1.VerifyWasCalledOnce().Leave(1)
	la2.VerifyWasCalledOnce().Leave(2)

	for _, l := range []Leavable{la1, la2} {
		if an.IsAnchored(l) {
			t.Errorf("expect %v not anchored", l)
		}
	}

	la3 := NewMockLeavable(pegomock.WithT(t))
	if an.IsAnchored(la3) {
		t.Errorf("expect IsAnchored returns false")
	}
	if an.DelAnchor(la3) {
		t.Errorf("expect DelAnchor returns false")
	}
}

func TestRoom_Join(t *testing.T) {
	s := NewMockSender(pegomock.WithT(t))
	r := NewRoom("test_room", s)

	ch := new(WSChannel)
	r.Join(ch)
	if !ch.IsAnchored(r) {
		t.Errorf("expect %v anchored", ch)
	}
	ch.DelAnchors()
	if ch.IsAnchored(r) {
		t.Errorf("expect %v not anchored", ch)
	}
	if r.Join(ch) {
		t.Error("expect room discard")
	}
}

func TestRoom_Join_Twice(t *testing.T) {
	s := NewMockSender(pegomock.WithT(t))
	r := NewRoom("test_room", s)

	ch := new(WSChannel)
	r.Join(ch)
	if !r.Join(ch) {
		t.Error("expect join success")
	}
	if r.subs.Len() != 1 {
		t.Errorf("expect member len: %d, got: %d", 1, r.subs.Len())
	}
	ch.DelAnchors()
	if r.subs.Len() != 0 {
		t.Errorf("expect member len: %d, got: %d", 0, r.subs.Len())
	}
}

func TestRoom_Route(t *testing.T) {
	s := NewMockSender(pegomock.WithT(t))
	r := NewRoom("test_room", s)

	ch1 := new(WSChannel)
	ch2 := new(WSChannel)
	ch3 := new(WSChannel)
	r.Join(ch1)
	r.Join(ch2)
	r.Join(ch3)

	m := NewMockMessage(pegomock.WithT(t))
	r.Route(m)
	ch3.DelAnchor(r)
	r.Route(m)
	s.VerifyWasCalled(pegomock.Times(2)).SendTo(ch1, m)
	s.VerifyWasCalled(pegomock.Times(2)).SendTo(ch2, m)
	s.VerifyWasCalledOnce().SendTo(ch3, m)
}

func TestBucket_JoinRoom(t *testing.T) {
	s := NewMockSender(pegomock.WithT(t))
	b := NewBucket("test_bucket", s, 32)

	ch1 := new(WSChannel)
	b.JoinRoom("1", ch1)
	if len(b.rooms) != 1 {
		t.Errorf("expect rooms: %d, got: %d", 1, len(b.rooms))
	}
	ch2 := new(WSChannel)
	b.JoinRoom("1", ch2)
	if len(b.rooms) != 1 {
		t.Errorf("expect rooms: %d, got: %d", 1, len(b.rooms))
	}
	ch3 := new(WSChannel)
	b.JoinRoom("2", ch3)
	if len(b.rooms) != 2 {
		t.Errorf("expect rooms: %d, got: %d", 2, len(b.rooms))
	}
	ch3.DelAnchors()
	if len(b.rooms) != 1 {
		t.Errorf("expect rooms: %d, got: %d", 1, len(b.rooms))
	}
	ch1.DelAnchors()
	ch2.DelAnchors()
	if len(b.rooms) != 0 {
		t.Errorf("expect rooms: %d, got: %d", 0, len(b.rooms))
	}
}

func TestBucket_JoinRoom_Twice(t *testing.T) {
	s := NewMockSender(pegomock.WithT(t))
	b := NewBucket("test_bucket", s, 32)

	ch1 := new(WSChannel)
	b.JoinRoom("1", ch1)
	r1 := b.rooms["1"]
	b.JoinRoom("1", ch1)
	r2 := b.rooms["1"]

	if r1 != r2 {
		t.Error("join room twince, room changed")
	}
}

func TestBucket_LeaveRoom(t *testing.T) {
	s := NewMockSender(pegomock.WithT(t))
	b := NewBucket("test_bucket", s, 32)

	ch1 := new(WSChannel)
	ch2 := new(WSChannel)
	ch3 := new(WSChannel)
	b.JoinRoom("1", ch1)
	b.JoinRoom("1", ch2)
	b.JoinRoom("2", ch3)
	b.LeaveRoom("1", ch1)
	if len(b.rooms) != 2 {
		t.Errorf("expect rooms: %d, got: %d", 2, len(b.rooms))
	}
	b.LeaveRoom("1", ch2)
	if len(b.rooms) != 1 {
		t.Errorf("expect rooms: %d, got: %d", 1, len(b.rooms))
	}
	b.LeaveRoom("2", ch3)
	if len(b.rooms) != 0 {
		t.Errorf("expect rooms: %d, got: %d", 0, len(b.rooms))
	}
}

func TestBucket_Leave_SameRoomName(t *testing.T) {
	s := NewMockSender(pegomock.WithT(t))
	b := NewBucket("test_bucket", s, 32)

	ch := new(WSChannel)
	b.JoinRoom("test_room", ch)
	r := NewRoom("test_room", s)
	b.Leave(r)
	if len(b.rooms) != 1 {
		t.Errorf("expect rooms: %d, got: %d", 1, len(b.rooms))
	}
	b.Leave(b.rooms["test_room"])
	if len(b.rooms) != 0 {
		t.Errorf("expect rooms: %d, got: %d", 0, len(b.rooms))
	}
}

func TestBucket_Route(t *testing.T) {
	s := NewMockSender(pegomock.WithT(t))
	b := NewBucket("test_bucket", s, 32)

	ch1 := new(WSChannel)
	ch2 := new(WSChannel)
	ch3 := new(WSChannel)
	b.JoinRoom("1", ch1)
	b.JoinRoom("1", ch2)
	b.JoinRoom("2", ch3)

	m1 := NewMockMessage(pegomock.WithT(t))
	pegomock.When(m1.GetRType()).ThenReturn(pb.RType_ROOM)
	pegomock.When(m1.GetRId()).ThenReturn("1")
	b.Route(m1)

	s.VerifyWasCalledOnce().SendTo(ch1, m1)
	s.VerifyWasCalledOnce().SendTo(ch2, m1)
	s.VerifyWasCalled(pegomock.Times(0)).SendTo(ch3, m1)

	m2 := NewMockMessage(pegomock.WithT(t))
	pegomock.When(m2.GetRType()).ThenReturn(pb.RType_ROOM)
	pegomock.When(m2.GetRId()).ThenReturn("2")
	b.Route(m2)

	s.VerifyWasCalled(pegomock.Times(0)).SendTo(ch1, m2)
	s.VerifyWasCalled(pegomock.Times(0)).SendTo(ch2, m2)
	s.VerifyWasCalledOnce().SendTo(ch3, m2)
}
