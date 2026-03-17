package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

func ReadCommand(r *bufio.Reader) ([]string, error) {

	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if line[0] != '*' {
		return nil, fmt.Errorf("invalid RESP array")
	}

	count, err := strconv.Atoi(line[1 : len(line)-2])
	if err != nil {
		return nil, err
	}

	args := make([]string, count)

	for i := 0; i < count; i++ {

		bulkHeader, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if bulkHeader[0] != '$' {
			return nil, fmt.Errorf("invalid bulk string")
		}

		length, err := strconv.Atoi(bulkHeader[1 : len(bulkHeader)-2])
		if err != nil {
			return nil, err
		}

		data := make([]byte, length+2)

		_, err = io.ReadFull(r, data)
		if err != nil {
			return nil, err
		}

		args[i] = string(data[:length])
	}

	return args, nil
}

func WriteSimpleString(w *bufio.Writer, s string) error {

	_, err := w.WriteString("+" + s + "\r\n")
	return err
}

func WriteError(w *bufio.Writer, msg string) error {

	_, err := w.WriteString("-" + msg + "\r\n")
	return err
}

func WriteBulkString(w *bufio.Writer, s string) error {

	_, err := w.WriteString(
		fmt.Sprintf("$%d\r\n%s\r\n", len(s), s),
	)

	return err
}

func WriteNull(w *bufio.Writer) error {

	_, err := w.WriteString("$-1\r\n")
	return err
}
