package service

import (
	"encoding/binary"
)

type AudioBuffer struct {
	buffer []byte
	pos    int
}

func (rw *AudioBuffer) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func (rw *AudioBuffer) SetPos(pos int) {
	rw.pos = pos
}

func (rw *AudioBuffer) GetBuffer() []byte {
	return rw.buffer
}

func (rw *AudioBuffer) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (rw *AudioBuffer) WriteFloat(p []float32) (int, error) {
	return len(p), binary.Write(rw, binary.LittleEndian, p)
}

func (rw *AudioBuffer) ReadFloat(out []float32) (int, error) {
	return len(out), binary.Read(rw, binary.LittleEndian, out)
}
