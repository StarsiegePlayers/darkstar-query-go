package protocol

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StatusBytesTestSuite struct {
	suite.Suite
	byte  StatusByte
	bytes *StatusByteStruct
}

func (t *StatusBytesTestSuite) SetupTest() {
	t.byte = 0
	t.bytes = new(StatusByteStruct)
}

func (t *StatusBytesTestSuite) TestUnmarshalByte_Even() {
	question := StatusByte(0xAA)

	t.bytes.Reserved2 = true
	t.bytes.WON = true
	t.bytes.Started = true
	t.bytes.Dedicated = true

	t.byte = question
	output := t.byte.Struct()

	t.Assert().Equal(*t.bytes, output)
}

func (t *StatusBytesTestSuite) TestUnmarshalByte_Odd() {
	question := StatusByte(0x55)

	t.bytes.Reserved1 = true
	t.bytes.Dynamix = true
	t.bytes.AllowOldClients = true
	t.bytes.Protected = true

	t.byte = question
	output := t.byte.Struct()

	t.Assert().Equal(*t.bytes, output)
}

func (t *StatusBytesTestSuite) TestUnmarshalByte_All() {
	question := StatusByte(0xFF)

	t.bytes.Reserved1 = true
	t.bytes.Dynamix = true
	t.bytes.AllowOldClients = true
	t.bytes.Protected = true
	t.bytes.Reserved2 = true
	t.bytes.WON = true
	t.bytes.Started = true
	t.bytes.Dedicated = true

	t.byte = question
	output := t.byte.Struct()

	t.Assert().Equal(*t.bytes, output)
}

func (t *StatusBytesTestSuite) TestUnmarshalByte_None() {
	question := StatusByte(0x00)

	t.byte = question
	output := t.byte.Struct()

	t.Assert().Equal(*t.bytes, output)
}

func (t *StatusBytesTestSuite) TestMarshalByte_Even() {
	question := byte(0xAA)

	t.bytes.Reserved2 = true
	t.bytes.WON = true
	t.bytes.Started = true
	t.bytes.Dedicated = true

	output := t.bytes.MarshalBinary()

	t.Assert().Equal(question, output)
}

func (t *StatusBytesTestSuite) TestMarshalByte_Odd() {
	question := byte(0x55)

	t.bytes.Reserved1 = true
	t.bytes.Dynamix = true
	t.bytes.AllowOldClients = true
	t.bytes.Protected = true
	output := t.bytes.MarshalBinary()

	t.Assert().Equal(question, output)
}

func (t *StatusBytesTestSuite) TestMarshalByte_All() {
	question := byte(0xFF)

	t.bytes.Reserved1 = true
	t.bytes.Dynamix = true
	t.bytes.AllowOldClients = true
	t.bytes.Protected = true
	t.bytes.Reserved2 = true
	t.bytes.WON = true
	t.bytes.Started = true
	t.bytes.Dedicated = true
	output := t.bytes.MarshalBinary()

	t.Assert().Equal(question, output)
}

func (t *StatusBytesTestSuite) TestMarshalByte_None() {
	question := byte(0x00)

	output := t.bytes.MarshalBinary()

	t.Assert().Equal(question, output)
}

func TestStatusByte_StatusBytesTestSuite(t *testing.T) {
	suite.Run(t, new(StatusBytesTestSuite))
}
