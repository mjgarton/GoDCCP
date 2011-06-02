// Copyright 2010 GoDCCP Authors. All rights reserved.
// Use of this source code is governed by a 
// license that can be found in the LICENSE file.

package dccp

func (h *Header) HasAckNo() bool { return getAckNoSubheaderSize(h.Type, h.X) > 0 }

func NewHeaderSkeleton(htype byte, sourcePort, destPort uint16) *Header {
	return &Header{
		SourcePort: sourcePort,
		DestPort:   destPort,
		Type:       htype,
		X:          true,
	}
}

// NewResetHeader() creates a new Reset header
func NewResetHeader(resetCode byte, sourcePort, destPort uint16) *Header {
	return &Header{
		SourcePort: sourcePort,
		DestPort:   destPort,
		Type:       Reset,
		X:          true,
		ResetCode:  resetCode,
	}
}

// NewCloseHeader() creates a new Close header
func NewCloseHeader(sourcePort, destPort uint16) *Header {
	return &Header{
		SourcePort: sourcePort,
		DestPort:   destPort,
		Type:       Close,
		X:          true,
	}
}

// NewAckHeader() creates a new Ack header
func NewAckHeader(sourcePort, destPort uint16) *Header {
	return &Header{
		SourcePort: sourcePort,
		DestPort:   destPort,
		Type:       Ack,
		X:          true,
	}
}

// NewDataHeader() creates a new Data header
func NewDataHeader(data []byte, sourcePort, destPort uint16) *Header {
	return &Header{
		SourcePort: sourcePort,
		DestPort:   destPort,
		Type:       Data,
		X:          true,
		Data:       data,
	}
}

// NewDataAckHeader() creates a new DataAck header
func NewDataAckHeader(data []byte, sourcePort, destPort uint16) *Header {
	return &Header{
		SourcePort: sourcePort,
		DestPort:   destPort,
		Type:       DataAck,
		X:          true,
		Data:       data,
	}
}

// NewSyncHeader() creates a new Sync header
func NewSyncHeader(sourcePort, destPort uint16) *Header {
	return &Header{
		SourcePort: sourcePort,
		DestPort:   destPort,
		Type:       Sync,
		X:          true,
	}
}

// NewSyncAckHeader() creates a new Sync header
func NewSyncAckHeader(sourcePort, destPort uint16) *Header {
	return &Header{
		SourcePort: sourcePort,
		DestPort:   destPort,
		Type:       SyncAck,
		X:          true,
	}
}

// NewResponseHeader() creates a new Response header
func NewResponseHeader(serviceCode uint32, sourcePort, destPort uint16) *Header {
	return &Header{
		SourcePort:  sourcePort,
		DestPort:    destPort,
		Type:        Response,
		X:           true,
		ServiceCode: serviceCode,
	}
}
