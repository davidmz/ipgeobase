package mmc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	BodyFormatError  = errors.New("Invalid body data")
	LineTooLongError = errors.New("Request line too long")
	CloseConnError   = errors.New("Close")
)

type Handler interface {
	ServeMemcache(*Request, *Response) error
}

type HandlerFunc func(*Request, *Response) error

func (h HandlerFunc) ServeMemcache(req *Request, res *Response) error { return h(req, res) }

type Request struct {
	Command string
	Args    []string

	conn *connection
}

type Response struct {
	conn *connection
}

func NewSession(rwc io.ReadWriteCloser, handler Handler) {
	c := &connection{rwc}
	c.run(handler)
}

func (r *Request) ReadBody(length int) ([]byte, error) { return r.conn.readRequestBody(length) }

func (r *Response) Status(status string) error {
	_, err := io.WriteString(r.conn, status+eol)
	return err
}
func (r *Response) UnknownCommandError() error          { return r.Status("ERROR") }
func (r *Response) ClientError(msg string) error        { return r.Status("CLIENT_ERROR " + msg) }
func (r *Response) ServerError(msg string) error        { return r.Status("SERVER_ERROR " + msg) }
func (r *Response) Value(key string, body []byte) error { return r.ValueFull(key, body, 0, 0) }
func (r *Response) ValueFull(key string, body []byte, flags uint32, cas uint64) error {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "VALUE %s %d %d %d", key, flags, len(body), cas)
	buf.WriteString(eol)
	buf.Write(body)
	buf.WriteString(eol)
	_, err := buf.WriteTo(r.conn)
	return err
}

/////////////////////////////

const (
	eol           = "\r\n"
	maxLineLength = 8192
)

type connection struct {
	io.ReadWriteCloser
}

func (c *connection) run(handler Handler) {
	for {
		resp := &Response{c}
		req, err := c.readRequest()
		if err != nil {
			break
		}
		err = handler.ServeMemcache(req, resp)
		if err != nil {
			break
		}
	}
	c.Close()
}

func (c *connection) readRequest() (*Request, error) {
	line, err := readLine(c, maxLineLength)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(string(line), " ")
	if len(parts) == 0 {
		parts = append(parts, "")
	}
	return &Request{
		Command: parts[0],
		Args:    parts[1:],
		conn:    c,
	}, nil
}

func (c *connection) readRequestBody(length int) ([]byte, error) {
	body := make([]byte, length+len(eol))
	if _, err := io.ReadFull(c, body); err != nil {
		return nil, err
	}
	if string(body[length:]) != eol {
		return nil, BodyFormatError
	}
	return body[:length], nil
}

// readLine читает данные, пока в них не встретится комбинация "\r\n" или
// не произойдёт ошибка чтения.
// При успешном чтении возвращается прочитанная строка (без финальных "\r\n"),
// иначе — nil и ошибка
func readLine(r io.Reader, maxLength int) ([]byte, error) {
	var (
		ptr int
		buf []byte = make([]byte, maxLength+2)
	)
	for ptr < len(buf) {
		n, err := r.Read(buf[ptr : ptr+1])
		ptr += n
		if ptr >= 2 && buf[ptr-1] == '\n' && buf[ptr-2] == '\r' {
			break
		} else if err != nil {
			return nil, err
		} else if ptr > maxLength {
			return nil, LineTooLongError
		}
	}
	return buf[0 : ptr-2], nil
}
