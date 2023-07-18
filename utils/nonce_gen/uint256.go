package nonce_gen

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

func ShaHash(b []byte, out []byte) {
	s := sha256.New()
	s.Write(b[:])
	tmp := s.Sum(nil)
	s.Reset()
	s.Write(tmp)
	copy(out[:], s.Sum(nil))
}

const Uint256IdxLen = 8 // The bigger it is, the more memory is needed, but lower chance of a collision

type Uint256 struct {
	Hash [32]byte
}

func NewUint256(h []byte) (res *Uint256) {
	res = new(Uint256)
	copy(res.Hash[:], h)
	return
}

// Get from MSB hexstring
func NewUint256FromString(s string) (res *Uint256) {
	d, e := hex.DecodeString(s)
	if e != nil {
		//println("NewUint256FromString", s, e.Error())
		return
	}
	if len(d) != 32 {
		//println("NewUint256FromString", s, "not 32 bytes long")
		return
	}
	res = new(Uint256)
	for i := 0; i < 32; i++ {
		res.Hash[31-i] = d[i]
	}
	return
}

func NewSha2Hash(data []byte) (res *Uint256) {
	res = new(Uint256)
	ShaHash(data, res.Hash[:])
	return
}

func (u *Uint256) Bytes() []byte {
	return u.Hash[:]
}

func (u *Uint256) String() (s string) {
	for i := 0; i < 32; i++ {
		s += fmt.Sprintf("%02x", u.Hash[31-i])
	}
	return
}

func (u *Uint256) Equal(o *Uint256) bool {
	return bytes.Equal(u.Hash[:], o.Hash[:])
}

func (u *Uint256) Calc(data []byte) {
	ShaHash(data, u.Hash[:])
}

func BIdx(hash []byte) (o [8]byte) {
	copy(o[:], hash[:8])
	return
}

func (u *Uint256) BIdx() (o [8]byte) {
	o = BIdx(u.Hash[:])
	return
}

func (u *Uint256) BigInt() *big.Int {
	var buf [32]byte
	for i := range buf {
		buf[i] = u.Hash[31-i]
	}
	return new(big.Int).SetBytes(buf[:])
}

func GenNonce(s string) string {
	r := NewSha2Hash([]byte(s))
	return r.BigInt().String()

}
