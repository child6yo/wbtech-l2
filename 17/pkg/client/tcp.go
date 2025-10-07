package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

// TCPClient - TCP соединение, способное читать сокет и писать в него.
type TCPClient struct {
	Host string
	Port string

	connectionTimeout time.Duration
	connection        net.Conn
}

// ConnectTCP создает новое TCP-соединение. Принимает хост, порт и таймаут соединения.
func ConnectTCP(host, port string, timeout time.Duration) (*TCPClient, error) {
	addr := host + ":" + port
	fmt.Println(addr)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}

	return &TCPClient{
		Host:              host,
		Port:              port,
		connection:        conn,
		connectionTimeout: timeout}, nil
}

// Relay запускает запись из сокета в reader и чтение из сокета во writer.
func (tc *TCPClient) Relay(ctx context.Context, reader io.Reader, writer io.Writer) error {
	if tc.connection == nil {
		return errors.New("not connected")
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		tc.read(ctx, writer)
	}()
	go func() {
		defer wg.Done()
		tc.write(ctx, reader)
	}()

	wg.Wait()
	tc.disconnect()

	return nil
}

func (tc *TCPClient) read(ctx context.Context, writer io.Writer) {
	scanner := bufio.NewScanner(tc.connection)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			text := scanner.Text()
			_, err := writer.Write([]byte(text + "\n"))
			if err != nil {
				log.Println(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("scan error: %v", err)

	}
}

func (tc *TCPClient) write(ctx context.Context, reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			text := scanner.Text()
			_, err := tc.connection.Write([]byte(text + "\n"))
			if err != nil {
				log.Println(err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("scan error: %v", err)

	}
}

func (tc *TCPClient) disconnect() {
	err := tc.connection.Close()
	if err != nil {
		log.Println(err)
	}
}
