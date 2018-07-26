package worker

import (
	"bytes"
	"testing"
)

func Test_crbuffer_Write(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name string
		args args
		buf  []byte
	}{
		{"simple", args{[]byte("a")}, []byte("a")},
		{"ln", args{[]byte("a\nb")}, []byte("a\nb")},
		{"cr", args{[]byte("a\nb\r1")}, []byte("a\n1")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &crbuffer{}
			_, err := c.Write(tt.args.p)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !bytes.Equal(c.buf, tt.buf) {
				t.Errorf("error processing buffer. got: %q\nexpected: %q", c.buf, tt.buf)
			}
		})
	}
}
