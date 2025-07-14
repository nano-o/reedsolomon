package reedsolomon

import (
	"bytes"
	"math/rand"
	"testing"
	"time"
)

// TestHighParityLossRecovery verifies that encoding with more parity
// than data shards can recover even when two thirds of shards are lost.
func TestHighParityLossRecovery(t *testing.T) {
	parallelIfNotShort(t)
	dataSize := 1 << 20 // 1MB
	data := make([]byte, dataSize)
	fillRandom(data)

	enc, err := New(10, 20, testOptions()...)
	if err != nil {
		t.Fatal(err)
	}

	shards, err := enc.Split(data)
	if err != nil {
		t.Fatal(err)
	}

	if err := enc.Encode(shards); err != nil {
		t.Fatal(err)
	}

	ok, err := enc.Verify(shards)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("verification failed")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	loss := rng.Perm(len(shards))[:20]
	for _, idx := range loss {
		shards[idx] = nil
	}

	if err := enc.Reconstruct(shards); err != nil {
		t.Fatal(err)
	}

	ok, err = enc.Verify(shards)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("verification failed after reconstruct")
	}

	buf := new(bytes.Buffer)
	if err := enc.Join(buf, shards, len(data)); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), data) {
		t.Fatal("recovered data does not match original")
	}
}
