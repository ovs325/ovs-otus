package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

func TestServerDisconnect_(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer func() { require.NoError(t, l.Close()) }()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("5s")
		require.NoError(t, err)

		client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
		require.NoError(t, client.Connect())
		defer func() { require.NoError(t, client.Close()) }()

		in.WriteString("I'm telnet client\n")
		err = client.Send()
		require.NoError(t, err)

		err = client.Receive()
		require.NoError(t, err)
	}()

	go func() { // Сервер
		defer wg.Done()

		conn, err := l.Accept()
		require.NoError(t, err)
		require.NotNil(t, conn)

		request := make([]byte, 1024)
		n, err := conn.Read(request)
		require.NoError(t, err)
		require.Equal(t, "I'm telnet client\n", string(request)[:n])

		_, err = conn.Write([]byte("Hello from NC\n"))
		require.NoError(t, err)

		_, err = conn.Write([]byte("Bye, client!\n"))
		require.NoError(t, err)

		require.NoError(t, conn.Close())
	}()

	wg.Wait()
}

func TestClientDisconnect(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:")
	require.NoError(t, err)
	defer func() { require.NoError(t, l.Close()) }()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		in := &bytes.Buffer{}
		in.WriteString("I\nwill be\nback!\n")
		out := &bytes.Buffer{}

		timeout, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
		require.NoError(t, client.Connect())
		defer func() { require.NoError(t, client.Close()) }()

		err = client.Send()
		require.NoError(t, err)
		require.Equal(t, "", out.String())
	}()

	go func() {
		defer wg.Done()

		conn, err := l.Accept()
		require.NoError(t, err)
		require.NotNil(t, conn)

		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					require.NoError(t, conn.Close())
					return
				}
				require.NoError(t, err)
			}
			fmt.Print(string(buf[:n]))
		}
	}()
	wg.Wait()
}
