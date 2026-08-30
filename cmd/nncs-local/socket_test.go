package main

import (
	"encoding/binary"
	"net"
	"testing"
	"time"
)

func TestReplySourceMatchesProbeType(t *testing.T) {
	regular, err := listenNNCS("127.0.0.1", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer regular.Close()
	same, err := listenNNCS("127.0.0.1", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer same.Close()
	other, err := listenNNCS("127.0.0.2", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer other.Close()
	client, err := listenNNCS("127.0.0.1", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	go serveNNCS(regular, "127.0.0.1", same, other)
	for _, kind := range []uint32{1, 2, 3, 4, 5, 101, 102, 103} {
		request := make([]byte, 16)
		binary.BigEndian.PutUint32(request, kind)
		if _, err := client.WriteToUDP(request, regular.LocalAddr().(*net.UDPAddr)); err != nil {
			t.Fatal(err)
		}
		client.SetReadDeadline(time.Now().Add(time.Second))
		response := make([]byte, 64)
		n, source, err := client.ReadFromUDP(response)
		if err != nil || n != 16 {
			t.Fatalf("type %d: size=%d error=%v", kind, n, err)
		}
		expected := regular
		if kind == 2 {
			expected = other
		} else if kind == 3 || kind == 102 {
			expected = same
		}
		if source.String() != expected.LocalAddr().String() {
			t.Fatalf("type %d: wrong reply source", kind)
		}
		if binary.BigEndian.Uint32(response[:4]) != kind || binary.BigEndian.Uint32(response[4:8]) != uint32(client.LocalAddr().(*net.UDPAddr).Port) {
			t.Fatal("wrong observed endpoint")
		}
	}
}
