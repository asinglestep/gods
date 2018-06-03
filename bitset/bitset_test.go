package bitset

import (
	"testing"
)

func Test_Set(t *testing.T) {
	bs := NewBitSet(0)
	for i := uint(0); i < 10; i++ {
		if i&0x01 != 0 {
			bs.Set(i)
		}
	}

	if bs.String() != "{1, 3, 5, 7, 9}" {
		t.Fatalf("BitSet set error, bs.String() want \"{1, 3, 5, 7, 9}\", but got %v\n", bs.String())
	}
}

func Test_Get(t *testing.T) {
	bs := NewBitSet(0)
	bs.Set(100)

	if !bs.Get(100) {
		t.Fatal("BitSet get error, get(100): flase")
	}
}

func Test_Clear(t *testing.T) {
	bs := NewBitSet(0)
	bs.Set(100)

	if !bs.Get(100) {
		t.Fatal("BitSet get error, get(100): flase")
	}

	bs.Clear(100)
	if bs.Get(100) {
		t.Fatal("BitSet get error, get(100): true")
	}
}

func Test_Flip(t *testing.T) {
	bs := NewBitSet(0)
	bs.Flip(10)

	if bs.String() != "{10}" {
		t.Fatalf("BitSet flip(10) error, bs.String() want \"{10}\", but got %v\n", bs.String())
	}

	bs.Set(3)
	if bs.String() != "{3, 10}" {
		t.Fatalf("BitSet set(3) error, bs.String() want \"{3, 10}\", but got %v\n", bs.String())
	}

	bs.Flip(3)
	if bs.String() != "{10}" {
		t.Fatalf("BitSet flip(3) error, bs.String() want \"{10}\", but got %v\n", bs.String())
	}
}

func Test_PrevSetBit(t *testing.T) {
	bs := NewBitSet(0)
	bs.Set(100)
	if bs.PrevSetBit(100) != 100 {
		t.Fatalf("BitSet PrevSetBit error, bs.PrevSetBit(100): %v\n", bs.PrevSetBit(100))
	}

	if bs.PrevSetBit(200) != 100 {
		t.Fatalf("BitSet PrevSetBit error, bs.PrevSetBit(100): %v\n", bs.PrevSetBit(200))
	}

	bs.Set(10)
	if bs.PrevSetBit(30) != 10 {
		t.Fatalf("BitSet PrevSetBit error, bs.PrevSetBit(30): %v\n", bs.PrevSetBit(30))
	}

	if bs.PrevSetBit(70) != 10 {
		t.Fatalf("BitSet PrevSetBit error, bs.PrevSetBit(70): %v\n", bs.PrevSetBit(70))
	}
}

func Test_PrevClearBit(t *testing.T) {
	bs := NewBitSet(0)
	bs.Set(100)
	if bs.PrevClearBit(100) != 99 {
		t.Fatalf("BitSet PrevClearBit error, bs.PrevClearBit(100): %v\n", bs.PrevClearBit(100))
	}

	if bs.PrevClearBit(200) != 127 {
		t.Fatalf("BitSet PrevClearBit error, bs.PrevClearBit(200): %v\n", bs.PrevClearBit(200))
	}
}

func Test_NextSetBit(t *testing.T) {
	bs := NewBitSet(0)
	bs.Set(100)
	if bs.NextSetBit(3) != 100 {
		t.Fatalf("BitSet NextSetBit error, bs.NextSetBit(3): %v\n", bs.NextSetBit(3))
	}

	if bs.NextSetBit(100) != 100 {
		t.Fatalf("BitSet NextSetBit error, bs.NextSetBit(100): %v\n", bs.NextSetBit(100))
	}
}

func Test_NextClearBit(t *testing.T) {
	bs := NewBitSet(0)
	bs.Set(100)
	if bs.NextClearBit(100) != 101 {
		t.Fatalf("BitSet NextCleanBit error, bs.NextClearBit(100): %v\n", bs.NextClearBit(100))
	}
}

func Test_And(t *testing.T) {
	b := NewBitSet(0)
	b.Set(1)
	b.Set(2)
	b.Set(5)

	s := NewBitSet(0)
	s.Set(2)
	s.Set(3)

	b.And(s)
	if b.String() != "{2}" {
		t.Fatalf("BitSet And error, b.String() want \"{2}\", but got %v\n", b.String())
	}
}

func Test_Or(t *testing.T) {
	b := NewBitSet(0)
	b.Set(1)
	b.Set(2)
	b.Set(5)

	s := NewBitSet(0)
	s.Set(2)
	s.Set(3)

	b.Or(s)
	if b.String() != "{1, 2, 3, 5}" {
		t.Fatalf("BitSet Or error, b.String() want \"{1, 2, 3, 5}\", but got %v\n", b.String())
	}
}

func Test_Xor(t *testing.T) {
	b := NewBitSet(0)
	b.Set(1)
	b.Set(2)
	b.Set(5)

	s := NewBitSet(0)
	s.Set(2)
	s.Set(3)

	b.Xor(s)
	if b.String() != "{1, 3, 5}" {
		t.Fatalf("BitSet Xor error, b.String() want \"{1, 3, 5}\", but got %v\n", b.String())
	}
}

func Test_AndNot(t *testing.T) {
	b := NewBitSet(0)
	b.Set(1)
	b.Set(2)
	b.Set(5)

	s := NewBitSet(0)
	s.Set(2)
	s.Set(3)

	b.AndNot(s)
	if b.String() != "{1, 5}" {
		t.Fatalf("BitSet AndNot error, b.String() want \"{1, 3, 5}\", but got %v\n", b.String())
	}
}
