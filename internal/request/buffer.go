package request

import (
	"bytes"
	"io"
)

const buffer_size = 8
const size_of_rn = 2

type buffer struct {
	Data   []byte
	Cursor int
}

func NewBuffer() buffer {
	return buffer{
		Data:   make([]byte, buffer_size),
		Cursor: 0,
	}
}

func (b *buffer) clear() {
	new_buffer := make([]byte, len(b.Data))
	b.Data = new_buffer
	b.Cursor = 0
}

func (b *buffer) readLine(reader io.Reader) (string, error) {
	not_found := -1

	/*
		fmt.Println()
		fmt.Println("locating rn index")
	*/

	idx := bytes.Index(b.Data[:b.Cursor], []byte("\r\n"))
	//fmt.Printf("Index = %d\n", idx)
	for idx == not_found {

		/*
			fmt.Println()
			fmt.Println("Start Iteration...")
			fmt.Println("Checking buffer size")
		*/

		if len(b.Data) <= b.Cursor {
			/*
				fmt.Printf("Buffer too small... len: %d, cursor: %d\n", len(b.Data), b.Cursor)
				fmt.Println("Enlarging buffer size")
			*/

			larger_buffer := make([]byte, len(b.Data)*2)
			copy(larger_buffer, b.Data)
			b.Data = larger_buffer

			/*
				fmt.Println("Buffer enlarged")
				fmt.Printf("New buffer size: %d\n", len(b.Data))
				fmt.Printf("New buffer: %v\n", b.Data)
			*/
		}

		/*
			fmt.Printf("Current Data: %v\n", b.Data)
			fmt.Printf("Current Cursor position: %d\n", b.Cursor)
			fmt.Println("Reading new data...")
		*/

		bytes_read, err := reader.Read(b.Data[b.Cursor:])
		b.Cursor += bytes_read

		/*
			fmt.Printf("New Cursor position: %d\n", b.Cursor)
			fmt.Printf("New Data: %v\n", b.Data)
		*/

		if err == io.EOF {
			//fmt.Printf("EOF...\n")
			line := string(b.Data)
			//fmt.Printf("Stringified: %s\n", string(b.Data))
			b.clear()
			return line, err
		}

		if err != nil {
			return "", err
		}

		idx = bytes.Index(b.Data[:b.Cursor], []byte("\r\n"))
		//fmt.Printf("Checking for rn... index: %d\n", idx)
	}

	//fmt.Printf("Relevant data: %v\n", b.Data[:idx])
	line := string(b.Data[:idx])
	bytes_read := len(b.Data[:idx]) + 2
	//fmt.Printf("Bytes Read: %d\n", bytes_read)

	new_buffer := make([]byte, len(b.Data))
	copy(new_buffer, b.Data[idx+size_of_rn:])
	b.Data = new_buffer
	b.Cursor -= bytes_read
	//fmt.Printf("Remaining buffer: %v Cursor: %d\n", b.Data, b.Cursor)
	//fmt.Printf("Line: %s\n", line)

	return line, nil
}

func (b *buffer) readRemaining(r io.Reader) ([]byte, error) {
	//fmt.Println("Time to read...")

	is_eof := false
	for !is_eof {
		if len(b.Data) <= b.Cursor {

			/*
				fmt.Printf("Buffer too small... len: %d, cursor: %d\n", len(b.Data), b.Cursor)
				fmt.Println("Enlarging buffer size")
			*/

			larger_buffer := make([]byte, len(b.Data)*2)
			copy(larger_buffer, b.Data)
			b.Data = larger_buffer

			/*
				fmt.Println("Buffer enlarged")
				fmt.Printf("New buffer size: %d\n", len(b.Data))
				fmt.Printf("New buffer: %v\n", b.Data)
			*/
		}

		/*
			fmt.Println("Starting data iteration")
			fmt.Printf("Current Data: %v\n", b.Data)
			fmt.Printf("Current Cursor position: %d\n", b.Cursor)
			fmt.Println("Reading new data...")
		*/

		bytes_read, err := r.Read(b.Data[b.Cursor:])
		b.Cursor += bytes_read

		/*
			fmt.Printf("New Cursor position: %d\n", b.Cursor)
			fmt.Printf("New Data: %v\n", b.Data)
		*/

		if err == io.EOF {
			is_eof = true
			break
		}

		if err != nil {
			return []byte{}, err
		}
	}

	return b.Data[:b.Cursor], nil
}
