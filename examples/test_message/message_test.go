package test_message

import (
	"bytes"
	"fmt"
	"testing"

	example_msgs "github.com/ppg/rosgo/examples/msg"
	"github.com/ppg/rosgo/msgs/std_msgs"
	"github.com/ppg/rosgo/ros"
)

//go:generate ros-gen-go msg --in=../msg/AllFieldTypes.msg
//go:generate ros-gen-go msg --in=../msg/Hello.msg

func TestInitialize(t *testing.T) {
	msg := example_msgs.AllFieldTypes{}
	fmt.Println(msg)
	fmt.Println(msg.H)

	if msg.B != 0 {
		t.Error(msg.B)
	}

	if msg.I8 != 0 {
		t.Error(msg.I8)
	}

	if msg.I16 != 0 {
		t.Error(msg.I16)
	}

	if msg.I32 != 0 {
		t.Error(msg.I32)
	}

	if msg.I64 != 0 {
		t.Error(msg.I64)
	}

	if msg.U8 != 0 {
		t.Error(msg.U8)
	}

	if msg.U16 != 0 {
		t.Error(msg.U16)
	}

	if msg.U32 != 0 {
		t.Error(msg.U32)
	}

	if msg.U64 != 0 {
		t.Error(msg.U64)
	}

	if msg.F32 != 0.0 {
		t.Error(msg.F32)
	}

	if msg.F64 != 0.0 {
		t.Error(msg.F64)
	}

	if msg.T.Sec != 0 || msg.T.NSec != 0 {
		t.Error(msg.T)
	}

	if msg.D.Sec != 0 || msg.D.NSec != 0 {
		t.Error(msg.D)
	}

	if msg.S != "" {
		t.Error(msg.S)
	}

	if msg.C.R != 0.0 || msg.C.G != 0.0 || msg.C.B != 0.0 || msg.C.A != 0 {
		t.Error(msg.C)
	}

	if len(msg.DynAry) != 0 {
		t.Error(msg.DynAry)
	}

	if len(msg.FixAry) != 2 || msg.FixAry[0] != 0 || msg.FixAry[1] != 0 {
		t.Error(msg.FixAry)
	}

}

func CheckBytes(t *testing.T, expected, actual []byte) {
	if !bytes.Equal(expected, actual) {
		t.Logf("expected: %+v", expected)
		t.Logf("  actual: %+v", actual)
		if len(expected) != len(actual) {
			t.Errorf("mismatched length: expected=%d, got=%d", len(expected), len(actual))
		} else {
			for i := 0; i < len(expected); i++ {
				if expected[i] != actual[i] {
					t.Errorf("result[%3d] is expected to be %02X but %02X", i, expected[i], actual[i])
				} else {
					t.Logf("%02X", expected[i])
				}
			}
		}
	}
}

func TestSerializeHeader(t *testing.T) {
	var msg std_msgs.Header
	msg.Seq = 0x89ABCDEF
	msg.Stamp = ros.NewTime(0x89ABCDEF, 0x01234567)
	msg.FrameID = "frame_id"
	var buf bytes.Buffer
	err := msg.Serialize(&buf)
	if err != nil {
		t.Error(err)
	}
	result := buf.Bytes()
	expected := []byte{
		// Header.Seq
		0xEF, 0xCD, 0xAB, 0x89,
		// Header.Stamp
		0xEF, 0xCD, 0xAB, 0x89,
		0x67, 0x45, 0x23, 0x01,
		// Header.FrameId
		0x08, 0x00, 0x00, 0x00,
		0x66, 0x72, 0x61, 0x6D, 0x65, 0x5F, 0x69, 0x64,
	}
	CheckBytes(t, expected, result)
}

func TestSerializeInt16(t *testing.T) {
	var msg std_msgs.Int16
	msg.Data = 0x0123
	var buf bytes.Buffer
	err := msg.Serialize(&buf)
	if err != nil {
		t.Error(err)
	}
	result := buf.Bytes()
	expected := []byte{
		0x23, 0x01,
	}
	CheckBytes(t, expected, result)
}

func TestSerializeInt32(t *testing.T) {
	var msg std_msgs.Int32
	msg.Data = 0x01234567
	var buf bytes.Buffer
	err := msg.Serialize(&buf)
	if err != nil {
		t.Error(err)
	}
	result := buf.Bytes()
	expected := []byte{
		0x67, 0x45, 0x23, 0x01,
	}
	CheckBytes(t, expected, result)
}

func getTestData() []byte {
	return []byte{
		// Header.Seq
		0xEF, 0xCD, 0xAB, 0x89,
		// Header.Stamp
		0xEF, 0xCD, 0xAB, 0x89,
		0x67, 0x45, 0x23, 0x01,
		// Header.FrameId
		0x08, 0x00, 0x00, 0x00,
		0x66, 0x72, 0x61, 0x6D, 0x65, 0x5F, 0x69, 0x64,
		// B
		0x01,
		// I8
		0x01,
		// I16
		0x23, 0x01,
		// I32
		0x67, 0x45, 0x23, 0x01,
		// I64
		0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01,
		// U8
		0x01,
		// U16
		0x23, 0x01,
		// U32
		0x67, 0x45, 0x23, 0x01,
		// U64
		0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01,
		// F32
		0xDB, 0x0F, 0x49, 0x40,
		// F64
		0x18, 0x2D, 0x44, 0x54, 0xFB, 0x21, 0x09, 0x40,
		// T
		0xEF, 0xCD, 0xAB, 0x89,
		0x67, 0x45, 0x23, 0x01,
		// D
		0xEF, 0xCD, 0xAB, 0x89,
		0x67, 0x45, 0x23, 0x01,
		// S
		0x0D, 0x00, 0x00, 0x00,
		0x48, 0x65, 0x6C, 0x6C, 0x6F, 0x2C, 0x20, 0x77, 0x6F, 0x72, 0x6C, 0x64, 0x21,
		// C
		0x00, 0x00, 0x80, 0x3F,
		0x00, 0x00, 0x00, 0x3F,
		0x00, 0x00, 0x80, 0x3E,
		0x00, 0x00, 0x00, 0x3E,
		// DynAry
		0x02, 0x00, 0x00, 0x00,
		0x67, 0x45, 0x23, 0x01,
		0xEF, 0xCD, 0xAB, 0x89,
		// FixAry
		0x02, 0x00, 0x00, 0x00,
		0x67, 0x45, 0x23, 0x01,
		0xEF, 0xCD, 0xAB, 0x89,
	}
}

func TestSerialize(t *testing.T) {
	var msg example_msgs.AllFieldTypes

	msg.H.Seq = 0x89ABCDEF
	msg.H.Stamp = ros.NewTime(0x89ABCDEF, 0x01234567)
	msg.H.FrameID = "frame_id"
	msg.B = 0x01
	msg.I8 = 0x01
	msg.I16 = 0x0123
	msg.I32 = 0x01234567
	msg.I64 = 0x0123456789ABCDEF
	msg.U8 = 0x01
	msg.U16 = 0x0123
	msg.U32 = 0x01234567
	msg.U64 = 0x0123456789ABCDEF
	msg.F32 = 3.141592653589793238462643383
	msg.F64 = 3.1415926535897932384626433832795028842
	msg.T = ros.NewTime(0x89ABCDEF, 0x01234567)
	msg.D = ros.NewDuration(0x89ABCDEF, 0x01234567)
	msg.S = "Hello, world!"
	msg.C = std_msgs.ColorRGBA{R: 1.0, G: 0.5, B: 0.25, A: 0.125}

	msg.DynAry = append(msg.DynAry, 0x01234567)
	msg.DynAry = append(msg.DynAry, 0x89ABCDEF)
	msg.FixAry[0] = 0x01234567
	msg.FixAry[1] = 0x89ABCDEF

	var buf bytes.Buffer
	err := msg.Serialize(&buf)
	if err != nil {
		t.Error(err)
	}

	result := buf.Bytes()
	expected := getTestData()
	CheckBytes(t, expected, result)
}

func TestDeserialize(t *testing.T) {
	source := getTestData()
	reader := bytes.NewReader(source)
	var msg example_msgs.AllFieldTypes
	err := msg.Deserialize(reader)
	if err != nil {
		t.Errorf("failed to deserialize message: %s", err)
	}

	if msg.H.Seq != 0x89ABCDEF {
		t.Errorf("msg.H.Seq incorrect; got=%+v", msg.H.Seq)
	}
	if msg.H.Stamp.Sec != 0x89ABCDEF || msg.H.Stamp.NSec != 0x01234567 {
		t.Errorf("msg.H.Stamp incorrect; got=%+v", msg.H.Stamp)
	}
	if msg.H.FrameID != "frame_id" {
		t.Errorf("msg.H.FrameID incorrect; got=%+v", msg.H.FrameID)
	}
	if msg.B != 0x01 {
		t.Errorf("msg.B incorrect; got=%+v", msg.B)
	}
	if msg.I8 != 0x01 {
		t.Errorf("msg.I8 incorrect; got=%+v", msg.I8)
	}
	if msg.I16 != 0x0123 {
		t.Errorf("msg.I16 incorrect; got=%+v", msg.I16)
	}
	if msg.I32 != 0x01234567 {
		t.Errorf("msg.I32 incorrect; got=%+v", msg.I32)
	}
	if msg.I64 != 0x0123456789ABCDEF {
		t.Errorf("msg.I64 incorrect; got=%+v", msg.I64)
	}
	if msg.U8 != 0x01 {
		t.Errorf("msg.U8 incorrect; got=%+v", msg.U8)
	}
	if msg.U16 != 0x0123 {
		t.Errorf("msg.U16 incorrect; got=%+v", msg.U16)
	}
	if msg.U32 != 0x01234567 {
		t.Errorf("msg.U32 incorrect; got=%+v", msg.U32)
	}
	if msg.U64 != 0x0123456789ABCDEF {
		t.Errorf("msg.U64 incorrect; got=%+v", msg.U64)
	}
	if msg.F32 != 3.141592653589793238462643383 {
		t.Errorf("msg.F32 incorrect; got=%+v", msg.F32)
	}
	if msg.F64 != 3.1415926535897932384626433832795028842 {
		t.Errorf("msg.F64 incorrect; got=%+v", msg.F64)
	}
	if msg.T.Sec != 0x89ABCDEF || msg.T.NSec != 0x01234567 {
		t.Errorf("msg.T incorrect; got=%+v", msg.T)
	}
	if msg.D.Sec != 0x89ABCDEF || msg.D.NSec != 0x01234567 {
		t.Errorf("msg.D incorrect; got=%+v", msg.D)
	}
	if msg.S != "Hello, world!" {
		t.Errorf("msg.S incorrect; got=%+v", msg.S)
	}
	if msg.C.R != 1.0 || msg.C.G != 0.5 || msg.C.B != 0.25 || msg.C.A != 0.125 {
		t.Errorf("msg.C incorrect; got=%+v", msg.C)
	}
	if msg.DynAry[0] != 0x01234567 || msg.DynAry[1] != 0x89ABCDEF {
		t.Errorf("msg.DynAry incorrect: got=%+v", msg.DynAry)
	}
	if msg.FixAry[0] != 0x01234567 || msg.FixAry[1] != 0x89ABCDEF {
		t.Errorf("msg.FixAry incorrect: got=%+v", msg.FixAry)
	}
	if reader.Len() != 0 {
		t.Errorf("reader has data remaining: %d", reader.Len())
	}
}
