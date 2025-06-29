package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Writer struct {
	writer io.Writer
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(reader io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(reader)}
}



func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}


func (r *Resp) readLine() (line []byte, n int, err error) {
	for{
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r'{
			break
		}
	}
	return 	line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid integer: %s", string(line))
	}
	return int(i64), n, nil
}

func (r *Resp) Read() (Value, error){
	// Read the first byte to determine the type of value
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type{
		case ARRAY:
			return r.readArray()
		case BULK:
			return r.readBulk()	
		default:
			fmt.Println("Unknown type:", string(_type))
			return Value{}, fmt.Errorf("unknown type: %s", string(_type))
	}
}

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// Read the length of the array
	len, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	v.array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.array = append(v.array, val)
	}
	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}

	bulk:=	make([]byte, len)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	r.readLine() // Read the trailing \r\n (CRLF) sequence

	return v, nil
}

func (v Value) marshal() []byte {
	// This function should convert the Value struct to a byte slice
	// according to the RESP protocol.
	// Convert Go object (struct, map, slice, etc.) into a JSON string.
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

func (v Value) marshalString() []byte {
	// Convert the Value struct to a byte slice according to the RESP protocol for arrays.
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n') // without the CRLF, there will be a problem because the RESP client will not understand the response without it.
	return bytes
}	

func (v Value) marshalBulk() []byte {
	// Convert the Value struct to a byte slice according to the RESP protocol for bulk strings.
	var bytes []byte
	bytes = append(bytes, BULK) // $ sign
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...) // for example, if the bulk string has 10 characters, it will be $10
	bytes = append(bytes, '\r', '\n') // with CRLF, $10\r\n.
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalArray() []byte {
	// Convert the Value struct to a byte slice according to the RESP protocol for arrays.
	var bytes []byte
	len := len(v.array)
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...) // for example, if the array has 10 elements, it will be *10
	bytes = append(bytes, '\r', '\n') // with CRLF, *10\r\n
	// then marshal each element in the array
	for i:= 0; i < len; i++ {
		bytes = append(bytes, v.array[i].marshal()...)
	}
	
	return bytes
}

func (v Value) marshalError() []byte {

	// Convert the Value struct to a byte slice according to the RESP protocol for errors.
	var bytes []byte
	bytes = append(bytes, ERROR) // - sign
	bytes = append(bytes, v.str...) // error message
	bytes = append(bytes, '\r', '\n') // with CRLF, -error message\r\n
	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n") // RESP protocol for null values is $-1\r\n
}

func (w *Writer) Write(v Value) error {
	var bytes = v.marshal()
	_, err := w.writer.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write value: %w", err)
	}
	return nil
}