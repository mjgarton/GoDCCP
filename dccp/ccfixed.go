// Copyright 2011 GoDCCP Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package dccp

import (
	"os"
	"time"
)

type CCFixed struct {}

func (CCFixed) NewSender() dccp.SenderCongestionControl { 
	return newFixedRateSenderControl(1e9 / sendsPerSecond) 
}

func (CCFixed) NewReceiver() dccp.ReceiverCongestionControl {
	return newFixedRateReceiverControl()
}

// ---> Fixed-rate HC-Sender Congestion Control

type fixedRateSenderControl struct {
	Mutex
	every  int64 // Strobe every every nanoseconds
	strobeRead  chan int
	strobeWrite chan int
}

func newFixedRateSenderControl(every int64) *fixedRateSenderControl {
	strobe := make(chan int)
	return &fixedRateSenderControl{ every: every, strobeRead: strobe, strobeWrite: strobe }
}

func (scc *fixedRateSenderControl) Open() {
	go func() {
		for {
			scc.Lock()
			if scc.strobeWrite == nil {
				scc.Unlock()
				break
			}
			scc.strobeWrite <- 1
			scc.Unlock()
			time.Sleep(scc.every)
		}
	}()
}

const CCID_FIXED = 0xf

func (scc *fixedRateSenderControl) GetID() byte { return CCID_FIXED }

func (scc *fixedRateSenderControl) GetCCMPS() int32 { return 1e9 }

func (scc *fixedRateSenderControl) GetRTT() int64 { return RTT_DEFAULT }

func (scc *fixedRateSenderControl) OnWrite(htype byte, x bool, seqno, ackno int64, now int64) (ccval byte, options []*Option) { return 0, nil }

func (scc *fixedRateSenderControl) OnRead(fb *FeedbackHeader) os.Error { return nil }

func (scc *fixedRateSenderControl) OnIdle(now int64) os.Error { return nil }

func (scc *fixedRateSenderControl) Strobe() {
	<-scc.strobeRead 
}

func (scc *fixedRateSenderControl) Close() { 
	scc.Lock()
	defer scc.Unlock()
	if scc.strobeWrite != nil {
		close(scc.strobeWrite) 
		scc.strobeWrite = nil
	}
}

// ---> Fixed-rate HC-Receiver Congestion Control

type fixedRateReceiverControl struct {}

func newFixedRateReceiverControl() *fixedRateReceiverControl {
	return &fixedRateReceiverControl{}
}

func (rcc *fixedRateReceiverControl) Open() {}

func (rcc *fixedRateReceiverControl) GetID() byte { return CCID_FIXED }

func (rcc *fixedRateReceiverControl) OnWrite(htype byte, x bool, seqno, ackno int64, now int64) (options []*Option) { return nil }

func (rcc *fixedRateReceiverControl) OnRead(ff *FeedforwardHeader) os.Error { return nil }

func (rcc *fixedRateReceiverControl) OnIdle(now int64) os.Error { return nil }

func (rcc *fixedRateReceiverControl) Close() {}
