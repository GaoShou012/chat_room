package frontier

import (
	"testing"
	"time"

	pegomock "github.com/petergtz/pegomock"
)

func TestSender_SendTo_Block(t *testing.T) {
	s := NewSender(0)

	a := NewMockSendable(pegomock.WithT(t))
	m1 := NewMockMessage(pegomock.WithT(t))
	m2 := NewMockMessage(pegomock.WithT(t))
	s.SendTo(a, m1)
	s.SendTo(a, m2)

	a.VerifyWasCalledEventually(pegomock.Times(1), time.Millisecond).Send(m1)
	a.VerifyWasCalledEventually(pegomock.Times(1), time.Millisecond).Send(m2)
}

func TestSender_SendTo_UnBlock(t *testing.T) {
	s := NewSender(2)

	a := NewMockSendable(pegomock.WithT(t))
	m1 := NewMockMessage(pegomock.WithT(t))
	m2 := NewMockMessage(pegomock.WithT(t))
	s.SendTo(a, m1)
	s.SendTo(a, m2)

	a.VerifyWasCalledEventually(pegomock.Times(1), time.Millisecond).Send(m1)
	a.VerifyWasCalledEventually(pegomock.Times(1), time.Millisecond).Send(m2)
}

func TestBalanceSender(t *testing.T) {
	s1 := NewMockSender(pegomock.WithT(t))
	s2 := NewMockSender(pegomock.WithT(t))
	s := NewBalanceSender(s1, s2)

	m1 := NewMockMessage(pegomock.WithT(t))
	pegomock.When(m1.GetRId()).ThenReturn("m1")
	m2 := NewMockMessage(pegomock.WithT(t))
	pegomock.When(m2.GetRId()).ThenReturn("m2")

	ch1 := NewMockSendable(pegomock.WithT(t))
	pegomock.When(ch1.Hash()).ThenReturn(0)

	s.SendTo(ch1, m1)
	s.SendTo(ch1, m2)

	s1.VerifyWasCalledOnce().SendTo(ch1, m1)
	s1.VerifyWasCalledOnce().SendTo(ch1, m2)

	s2.VerifyWasCalled(pegomock.Never()).SendTo(ch1, m1)
	s2.VerifyWasCalled(pegomock.Never()).SendTo(ch1, m2)

	ch2 := NewMockSendable(pegomock.WithT(t))
	pegomock.When(ch2.Hash()).ThenReturn(1)

	m3 := NewMockMessage(pegomock.WithT(t))
	s.SendTo(ch2, m3)
	s1.VerifyWasCalled(pegomock.Never()).SendTo(ch2, m3)
	s2.VerifyWasCalledOnce().SendTo(ch2, m3)
}
