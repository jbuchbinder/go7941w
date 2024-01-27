package go7941w

import (
	"fmt"
	"log"
	"time"

	"github.com/tarm/serial"
)

const (
	GO7941W_CMD_READ_UID                 = 0x10
	GO7941W_CMD_WRITE_UID                = 0x11
	GO7941W_CMD_READ_SECTOR              = 0x12
	GO7941W_CMD_WRITE_SECTOR             = 0x13
	GO7941W_CMD_CHANGE_PASSWORD          = 0x14
	GO7941W_CMD_READ_ID                  = 0x15
	GO7941W_CMD_WRITE_T5577              = 0x16
	GO7941W_CMD_READ_DATA_ALL            = 0x17
	GO7941W_CMD_READ_TAG_CARD_ULTRALIGHT = 0x19
	GO7941W_CMD_READ_TAG_CARD_SECTOR     = 0x20
	GO7941W_CMD_WRITE_TAG_CARD_SECTOR    = 0x21
)

type Go7941W struct {
	Port    string
	Timeout int

	s *serial.Port
}

func (g *Go7941W) Open() error {
	c := &serial.Config{
		Name:        g.Port,
		Baud:        115200,
		Size:        8,
		Parity:      serial.ParityNone,
		ReadTimeout: time.Second * time.Duration(g.Timeout),
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		return err
	}
	g.s = s
	return nil
}

func (g *Go7941W) ReadUID() error {
	// 0x10 read UID number
	err := g.sendCommand(0x00, GO7941W_CMD_READ_SECTOR, []byte{0x00})
	if err != nil {
		return err
	}

	_, err = g.readResponse(6)
	return err
}

func (g *Go7941W) WriteUID(uid []byte) error {
	// 0x11 write UID number "4 bytes", use the default window code fffffffffff
	return g.sendCommand(0x00, GO7941W_CMD_WRITE_UID, uid)
}

func (g *Go7941W) ReadSector(sector byte, block byte, group byte, password []byte) ([]byte, error) {
	// 0x12 read specified sector

	ret := make([]byte, 6)

	if len(password) != 6 {
		return ret, fmt.Errorf("password len = %d", len(password))
	}

	data := []byte{}
	data = append(data, sector)
	data = append(data, block)
	data = append(data, group)
	data = append(data, password...)

	err := g.sendCommand(0x00, GO7941W_CMD_READ_SECTOR, data)
	if err != nil {
		return []byte{}, err
	}
	return g.readResponse(6)
}

func (g *Go7941W) WriteSector(sector byte, block byte, group byte, password []byte, sectorData []byte) error {
	// 0x13 write specified sector

	if len(password) != 6 {
		return fmt.Errorf("password len = %d", len(password))
	}

	data := []byte{}
	data = append(data, sector)
	data = append(data, block)
	data = append(data, group)
	data = append(data, password...)
	data = append(data, sectorData...)

	return g.sendCommand(0x00, GO7941W_CMD_CHANGE_PASSWORD, data)
}

func (g *Go7941W) ChangePassword(sector byte, group byte, oldPass []byte, newPass []byte) error {
	// 0x14 Change group A or group B password

	if len(oldPass) != 6 || len(newPass) != 6 {
		return fmt.Errorf("oldpass len = %d, newpass len = %d", len(oldPass), len(newPass))
	}

	data := []byte{}
	data = append(data, sector)
	data = append(data, group)
	data = append(data, oldPass...)
	data = append(data, newPass...)

	return g.sendCommand(0x00, GO7941W_CMD_CHANGE_PASSWORD, data)
}

func (g *Go7941W) ReadID() ([]byte, error) {
	// 0x15 Read ID card number

	data := []byte{}

	err := g.sendCommand(0x00, GO7941W_CMD_READ_ID, data)
	if err != nil {
		return []byte{}, err
	}
	return g.readResponse(32)
}

func (g *Go7941W) WriteT5577ID(id []byte) error {
	// 0x16 write T5577 card number
	data := []byte{}
	data = append(data, byte(len(id)))
	data = append(data, id...)

	return g.sendCommand(0, GO7941W_CMD_WRITE_T5577, data)
}

func (g *Go7941W) ReadAll(group byte, password []byte) ([]byte, error) {
	// 0x17 Read the data of all blocks in all local areas (M1-1K card)
	if group != 0x0A && group != 0x0B {
		return []byte{}, fmt.Errorf("group must be 0x0A or 0x0B")
	}
	if len(password) != 6 {
		return []byte{}, fmt.Errorf("password length != 6 (length was %d)", len(password))
	}

	data := []byte{}
	data = append(data, group)
	data = append(data, password...)

	err := g.sendCommand(0x00, GO7941W_CMD_READ_DATA_ALL, data)
	if err != nil {
		return []byte{}, err
	}

	return g.readResponse(128)
}

func (g *Go7941W) ReadUltracardID() ([]byte, error) {
	// 0x19 Read tag card UID (Ultralight card)
	err := g.sendCommand(0x00, GO7941W_CMD_READ_TAG_CARD_ULTRALIGHT, []byte{})
	if err != nil {
		return []byte{}, err
	}

	return g.readResponse(128)
}

func (g *Go7941W) WriteTagSector(sector byte, data []byte) error {
	// 0x21 write tag card sector
	cmd := []byte{}
	//0xab, 0xba, 0x00, 0x21, 0x05, 0x54, 0x00, 0x82, 0x10, 0x9e, 0x7c}
	log.Printf("%#v", cmd)

	// TODO: FIXME: IMPLEMENT

	return g.sendCommand(0, GO7941W_CMD_WRITE_TAG_CARD_SECTOR, cmd)
}

func (g *Go7941W) sendCommand(address byte, cmd byte, data []byte) error {
	buf := []byte{}
	buf = append(buf, 0xab, 0xba)
	buf = append(buf, address)
	buf = append(buf, cmd)
	buf = append(buf, byte(len(data)))
	if len(data) > 0 {
		buf = append(buf, data...)
	}
	buf = append(buf, g.xorChecksum(buf[2:]))
	return g.sendBytes(buf)
}

func (g *Go7941W) sendBytes(b []byte) error {
	_, err := g.s.Write(b)
	return err
}

func (g *Go7941W) xorChecksum(b []byte) byte {
	o := byte(0x00)
	for _, x := range b {
		o ^= x
	}
	return o
}

func (g *Go7941W) readResponse(maxLength int) ([]byte, error) {
	// Allow processing time
	time.Sleep(500 * time.Millisecond)

	out := make([]byte, maxLength)

	n, err := g.s.Read(out)
	if err != nil {
		return out, err
	}
	if n == 0 {
		return out, fmt.Errorf("0 bytes received")
	}

	return g.parseResponse(out)
}

func (g *Go7941W) parseResponse(in []byte) ([]byte, error) {
	l := len(in)
	if l < 4 {
		return []byte{}, fmt.Errorf("invalid response: %#v", in)
	}
	if g.xorChecksum(in[3:l-1]) != in[l-1] {
		return []byte{}, fmt.Errorf("XOR check failed: %#v", in)
	}
	if in[0] != 0xCD || in[1] != 0xDC {
		return []byte{}, fmt.Errorf("Invalid header: %#v", in)
	}
	if in[3] == 0x81 {
		dataLength := in[4]
		log.Printf("datalength = %d", dataLength)
		data := in[5 : 5+dataLength]
		return data, nil
	} else if in[3] == 0x80 {
		return []byte{}, fmt.Errorf("Op failed: %#v", in)
	} else {
		return []byte{}, fmt.Errorf("Unexpected return value: %#v", in)
	}
}

func (g *Go7941W) Close() error {
	if g.s != nil {
		g.s.Flush()
		return g.s.Close()
	}
	return nil
}
