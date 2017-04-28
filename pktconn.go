// pktconn.go - wrapper around io.ReadWriteCloser that implements gopacket's 'device'.
//
// To the extent possible under law, Ivan Markin waived all copyright
// and related or neighboring rights to pktconn, using the creative
// commons "cc0" public domain dedication. See LICENSE or
// <http://creativecommons.org/publicdomain/zero/1.0/> for full details.

package pktconn

import (
	"io"
	"time"

	"github.com/google/gopacket"
)

type PacketConn struct {
	rwc      io.ReadWriteCloser
	b        []byte
	reinject bool
}

func New(rwc io.ReadWriteCloser, bufsize int) (*PacketConn, error) {
	copier := &PacketConn{
		rwc: rwc,
		b:   make([]byte, bufsize),
	}
	return copier, nil
}

func (c *PacketConn) Close() error {
	return c.rwc.Close()
}

func (c *PacketConn) ReadPacketData() ([]byte, gopacket.CaptureInfo, error) {
	ci := gopacket.CaptureInfo{
		Timestamp:     time.Now(),
		CaptureLength: 0,
		Length:        0,
	}
	n, err := c.rwc.Read(c.b)
	if err != nil {
		return c.b[:0], ci, err
	}
	ci.CaptureLength = n
	ci.Length = n

	return c.b[:n], ci, err
}

func (c *PacketConn) WritePacketData(data []byte) error {
	_, err := c.rwc.Write(data)
	return err
}
