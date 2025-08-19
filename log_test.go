package compose

import (
	"bytes"
	"io"
	"testing"
)

func BenchmarkPrefixWriter_Write(b *testing.B) {
	prefix := "[TEST] "
	data := []byte("Hello, World!")
	
	// Benchmark old implementation
	b.Run("OldImplementation", func(b *testing.B) {
		oldWriter := &oldPrefixWriter{
			Out:    io.Discard,
			Prefix: prefix,
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			oldWriter.Write(data)
		}
	})
	
	// Benchmark new implementation
	b.Run("NewImplementation", func(b *testing.B) {
		newWriter := NewPrefixWriter(io.Discard, prefix)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			newWriter.Write(data)
		}
	})
}

func BenchmarkPrefixWriter_WriteLarge(b *testing.B) {
	prefix := "[TEST] "
	data := bytes.Repeat([]byte("Hello, World! "), 1000) // 15KB of data
	
	b.Run("OldImplementation", func(b *testing.B) {
		oldWriter := &oldPrefixWriter{
			Out:    io.Discard,
			Prefix: prefix,
		}
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			oldWriter.Write(data)
		}
	})
	
	b.Run("NewImplementation", func(b *testing.B) {
		newWriter := NewPrefixWriter(io.Discard, prefix)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			newWriter.Write(data)
		}
	})
}

// oldPrefixWriter represents the previous implementation for benchmarking
type oldPrefixWriter struct {
	Out    io.Writer
	Prefix string
}

func (w *oldPrefixWriter) Write(p []byte) (n int, err error) {
	n, err = w.Out.Write(append([]byte(w.Prefix), p...))
	if n > len(p) {
		n = len(p)
	}
	return
}

func TestPrefixWriter_Write(t *testing.T) {
	var buf bytes.Buffer
	writer := NewPrefixWriter(&buf, "[TEST] ")
	
	testData := []byte("Hello, World!")
	expected := "[TEST] Hello, World!"
	
	n, err := writer.Write(testData)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	
	if n != len(testData) {
		t.Errorf("Expected %d bytes written, got %d", len(testData), n)
	}
	
	if buf.String() != expected {
		t.Errorf("Expected output %q, got %q", expected, buf.String())
	}
}