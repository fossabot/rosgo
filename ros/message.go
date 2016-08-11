package ros

import (
	"encoding/binary"
	"fmt"
	"io"
)

// TODO(ppg): Investigate a generic Serialize and Deserialize for any Message
// that uses reflection (and maybe tags) like JSON, binary, etc.

type MessageType interface {
	Text() string
	MD5Sum() string
	Name() string
	NewMessage() Message
}

type Message interface {
	Serialize(io.Writer) error
	Deserialize(io.Reader) error
}

func SerializeMessageField(w io.Writer, fieldType string, pData interface{}) error {
	// Switch based on fieldType
	switch fieldType {
	case "string": // built-in
		pS, ok := pData.(*string)
		if !ok {
			return fmt.Errorf("expected *string, got %T", pData)
		}

		// Write size little endian
		err := binary.Write(w, binary.LittleEndian, uint32(len([]byte(*pS))))
		if err != nil {
			return fmt.Errorf("could not write string length: %s", err)
		}
		// Write data in-order
		// TODO(ppg): Look at io.Copy here instead?
		_, err = w.Write([]byte(*pS))
		if err != nil {
			return fmt.Errorf("could not write string: %s", err)
		}
		return nil

	case "time": // built-in
		pT, ok := pData.(*Time)
		if !ok {
			return fmt.Errorf("expected *ros.Time, go %T", pData)
		}

		// Write Sec and NSec little endian
		err := binary.Write(w, binary.LittleEndian, pT.Sec)
		if err != nil {
			return fmt.Errorf("could not write time#Sec: %s", err)
		}
		err = binary.Write(w, binary.LittleEndian, pT.NSec)
		if err != nil {
			return fmt.Errorf("could not write time#NSec: %s", err)
		}
		return nil

	case "duration": // built-in
		pD, ok := pData.(*Duration)
		if !ok {
			return fmt.Errorf("expected *ros.Duration, go %T", pData)
		}

		// Write Sec and NSec little endian
		err := binary.Write(w, binary.LittleEndian, pD.Sec)
		if err != nil {
			return fmt.Errorf("could not write duration#Sec: %s", err)
		}
		err = binary.Write(w, binary.LittleEndian, pD.NSec)
		if err != nil {
			return fmt.Errorf("could not write duration#NSec: %s", err)
		}
		return nil

	case "bool", "byte", "char", "float32", "float64", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64": // built-in
		// Write default type litte endian
		// TODO(ppg): Decide if we want to do some type assertion checking here too
		if err := binary.Write(w, binary.LittleEndian, pData); err != nil {
			return fmt.Errorf("could not write %T: %s", pData, err)
		}
		return nil

	default: // not built-in, ensure data is a Message and hand off
		serializer, ok := pData.(Message)
		if !ok {
			return fmt.Errorf("cannot serialize non-message type %T", pData)
		}
		return serializer.Serialize(w)
	}
}

func DeserializeMessageField(r io.Reader, fieldType string, pData interface{}) error {
	switch fieldType {
	case "string": // built-in
		pS, ok := pData.(*string)
		if !ok {
			return fmt.Errorf("expected *string, go %T", pData)
		}

		// Read size little endian
		var size uint32
		if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
			return fmt.Errorf("could not read string length: %s", err)
		}
		// Read data in-order
		// TODO(ppg): Look at io.Copy here instead?
		// FIXME(ppg): Shouldn't this be in-order like serialize?
		data := make([]byte, int(size))
		if err := binary.Read(r, binary.LittleEndian, data); err != nil {
			return fmt.Errorf("could not read string: %s", err)
		}
		*pS = string(data)
		return nil

	case "time": // built-in
		pT, ok := pData.(*Time)
		if !ok {
			return fmt.Errorf("expected *ros.Time, go %T", pData)
		}

		// Read Sec and NSec little endian
		err := binary.Read(r, binary.LittleEndian, &pT.Sec)
		if err != nil {
			return fmt.Errorf("could not read time#Sec: %s", err)
		}
		err = binary.Read(r, binary.LittleEndian, &pT.NSec)
		if err != nil {
			return fmt.Errorf("could not read time#NSec: %s", err)
		}
		return nil

	case "duration": // built-in
		pD, ok := pData.(*Duration)
		if !ok {
			return fmt.Errorf("expected *ros.Duration, go %T", pData)
		}

		// Read Sec and NSec little endian
		err := binary.Read(r, binary.LittleEndian, &pD.Sec)
		if err != nil {
			return fmt.Errorf("could not read duration#Sec: %s", err)
		}
		err = binary.Read(r, binary.LittleEndian, &pD.NSec)
		if err != nil {
			return fmt.Errorf("could not read duration#NSec: %s", err)
		}
		return nil

	case "bool", "byte", "char", "float32", "float64", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64": // built-in
		// Read default type litte endian
		// TODO(ppg): Decide if we want to do some type assertion checking here too
		if err := binary.Read(r, binary.LittleEndian, pData); err != nil {
			return fmt.Errorf("could not read %T: %s", pData, err)
		}
		return nil

	default: // not built-in, ensure data is a Message and hand off
		deserializer, ok := pData.(Message)
		if !ok {
			return fmt.Errorf("cannot deserialize non-message type %T", pData)
		}
		return deserializer.Deserialize(r)
	}
}
