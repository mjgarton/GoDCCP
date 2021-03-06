// Copyright 2011-2013 GoDCCP Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package dccp

func (c *Conn) readHeader() (h *Header, err error) {
	h, err = c.hc.Read()
	if err != nil {
		if err != ErrTimeout {
			c.amb.E(EventDrop, "Bad header", h)
		}
		return nil, err
	}
	// We don't support non-extended (short) SeqNo's 
	if !h.X {
		return nil, ErrUnsupported
	}
	return h, nil
}

// idleLoop polls the congestion control OnIdle method at regular intervals
// of approximately one RTT.
func (c *Conn) idleLoop() {
	for {
		c.pollCongestionControl()

		c.Lock()
		c.syncWithCongestionControl()
		rtt := c.socket.GetRTT()
		state := c.socket.GetState()
		c.Unlock()

		if state == CLOSED {
			break
		}
		// This emit prints very often. Use when really necessary
		//c.amb.E(EventIdle, "")
		c.env.Sleep(max64(RoundtripMin, min64(rtt, RoundtripDefault)))
	}
}

func (c *Conn) readLoop() {
	for {
		c.Lock()
		state := c.socket.GetState()
		rtt := c.socket.GetRTT()
		c.Unlock()
		if state == CLOSED {
			break
		}

		// Adjust read timeout
		if err := c.hc.SetReadExpire(5 * rtt); err != nil {
			c.amb.E(EventError, "SetReadExpire")
			c.abortQuietly()
			return
		}

		// Read next header
		h, err := c.readHeader()
		if err != nil {
			_, ok := err.(ProtoError)
			if ok {
				// Drop packets that are unsupported. Intended for forward compatibility.
				continue
			} else if err == ErrTimeout {
				// In the even of timeout, poll the congestion controls
				c.pollCongestionControl()
				continue
			} else {
				// Die if the underlying link is broken
				c.abortQuietly()
				return
			}
		}
		c.amb.E(EventRead, "", h)

		c.Lock()
		c.syncWithCongestionControl()
		if c.step2_ProcessTIMEWAIT(h) != nil {
			goto Done
		}
		if c.step3_ProcessLISTEN(h) != nil {
			goto Done
		}
		if c.step4_PrepSeqNoREQUEST(h) != nil {
			goto Done
		}
		if c.step5_PrepSeqNoForSync(h) != nil {
			goto Done
		}
		if c.step6_CheckSeqNo(h) != nil {
			goto Done
		}
		if c.step7_CheckUnexpectedTypes(h) != nil {
			goto Done
		}
		if c.step8_OptionsAndMarkAckbl(h) != nil {
			goto Done
		}
		if c.step9_ProcessReset(h) != nil {
			goto Done
		}
		if c.step10_ProcessREQUEST2(h) != nil {
			goto Done
		}
		if c.step11_ProcessRESPOND(h) != nil {
			goto Done
		}
		if c.step12_ProcessPARTOPEN(h) != nil {
			goto Done
		}
		if c.step13_ProcessCloseReq(h) != nil {
			goto Done
		}
		if c.step14_ProcessClose(h) != nil {
			goto Done
		}
		if c.step15_ProcessSync(h) != nil {
			goto Done
		}
		if c.step16_ProcessData(h) != nil {
			goto Done
		}
	Done:
		c.Unlock()
	}
	c.amb.E(EventInfo, "Read loop EXIT")
}

func (c *Conn) pollCongestionControl() {
	now := c.env.Now()
	if e := c.scc.OnIdle(now); e != nil {
		if re, ok := e.(CongestionReset); ok {
			c.abortWith(re.ResetCode())
			return
		}
		if e == CongestionAck {
			c.Lock()
			c.inject(c.generateAck())
			c.Unlock()
			return
		}
		c.amb.E(EventError, "Sender CC unknown idle error")
	}
	if e := c.rcc.OnIdle(now); e != nil {
		if re, ok := e.(CongestionReset); ok {
			c.abortWith(re.ResetCode())
			return
		}
		if e == CongestionAck {
			c.Lock()
			c.inject(c.generateAck())
			c.Unlock()
			return
		}
		c.amb.E(EventError, "Receiver CC unknown idle error")
	}
}

func (c *Conn) syncWithCongestionControl() {
	c.AssertLocked()
	c.socket.SetRTT(c.scc.GetRTT())
	c.socket.SetCCMPS(c.scc.GetCCMPS())
}

func (c *Conn) syncWithLink() {
	c.AssertLocked()
	c.socket.SetPMTU(int32(c.hc.GetMTU()))
}
