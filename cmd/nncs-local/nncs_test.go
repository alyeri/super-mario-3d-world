package main

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestMakeNNCSResponse(t *testing.T) {
	response := makeNNCSResponse(102, &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 54321}, "127.0.0.2")
	if got := binary.BigEndian.Uint32(response[0:4]); got != 102 {
		t.Fatalf("type=%d, want 102", got)
	}
	if got := binary.BigEndian.Uint32(response[4:8]); got != 54321 {
		t.Fatalf("port=%d, want 54321", got)
	}
	if got := net.IP(response[8:12]).String(); got != "127.0.0.1" {
		t.Fatalf("observed IP=%s, want 127.0.0.1", got)
	}
	if got := net.IP(response[12:16]).String(); got != "127.0.0.2" {
		t.Fatalf("server IP=%s, want 127.0.0.2", got)
	}
}
