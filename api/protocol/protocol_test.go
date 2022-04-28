package protocol

import (
	"runtime"
	"testing"
)

func BenchmarkTest(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := &Proto{
			Version:  1,
			Op:       OpHeartbeatReply,
			Checksum: 1,
			Seq:      1,
			Data:     []byte("hello world"),
		}
		data := make([]byte, len(p.Data)+HeaderLen)
		err := p.SerializeTo(data)
		if err != nil {
			b.Fatal(err)
		}
		//p1, err := DecodeFromBytes(data)
		//if err != nil {
		//	b.Fatal(err)
		//}
		//if p1.Version != p.Version || p1.Op != p.Op || p1.Checksum != p.Checksum || p1.Seq != p.Seq || string(p1.Data) != string(p.Data) {
		//	b.Fatalf("not equal %+v", p1)
		//}
	}
	runtime.GC()
}

func TestProto_SerializeTo(t *testing.T) {
	type fields struct {
		Version  uint16
		Op       uint16
		Checksum uint16
		Seq      uint16
		Data     []byte
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "TestProto_DecodeFromBytes",
			fields: fields{
				Version:  1,
				Op:       1,
				Checksum: 1,
				Seq:      1,
				Data:     []byte{1},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Proto{
				Version:  tt.fields.Version,
				Op:       tt.fields.Op,
				Checksum: tt.fields.Checksum,
				Seq:      tt.fields.Seq,
				Data:     tt.fields.Data,
			}
			data := make([]byte, HeaderLen)
			if err := p.SerializeTo(data); (err != nil) != tt.wantErr {
				t.Errorf("Proto.DecodeFromBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
