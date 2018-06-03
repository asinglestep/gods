package bitset

import (
	"bytes"
	"math/bits"
	"strconv"

	"github.com/asinglestep/gods/utils"
)

const (
	WORD_SIZE  = 64
	WORD_SHIFT = 6

	WORK_MASK uint64 = 0xffffffffffffffff
)

// BitSet BitSet
type BitSet struct {
	words     []uint64 // word数组
	wordInUse uint     // 使用的word数
}

// NewBitSet NewBitSet
func NewBitSet(nbits int) *BitSet {
	bs := &BitSet{}
	bs.words = make([]uint64, needWords(nbits))
	bs.wordInUse = 0

	return bs
}

// Set bitIndex位置1
func (b *BitSet) Set(bitIndex uint) {
	wordIndex := wordIndex(bitIndex)
	wordsRequired := wordIndex + 1

	if b.wordInUse < wordsRequired {
		b.grown(wordsRequired)
		b.wordInUse = wordsRequired
	}

	b.words[wordIndex] |= (1 << (bitIndex & (WORD_SIZE - 1)))
}

// Get bitIndex位的值是否为1
func (b *BitSet) Get(bitIndex uint) bool {
	wordIndex := wordIndex(bitIndex)
	if wordIndex >= b.wordInUse {
		return false
	}

	return b.words[wordIndex]&(1<<(bitIndex&(WORD_SIZE-1))) != 0
}

// Clear 清除bitIndex位的值
func (b *BitSet) Clear(bitIndex uint) {
	wordIndex := wordIndex(bitIndex)
	if wordIndex >= b.wordInUse {
		return
	}

	b.words[wordIndex] &^= 1 << (bitIndex & (WORD_SIZE - 1))
}

// Flip 将bitIndex指定的位置取反
func (b *BitSet) Flip(bitIndex uint) {
	wordIndex := wordIndex(bitIndex)
	wordsRequired := wordIndex + 1

	if b.wordInUse < wordsRequired {
		b.grown(wordsRequired)
		b.wordInUse = wordsRequired
	}

	b.words[wordIndex] ^= 1 << (bitIndex & (WORD_SIZE - 1))
}

// PrevSetBit 找到bitIndex位之前的第一个设置为1的位置
func (b *BitSet) PrevSetBit(bitIndex uint) int {
	var word uint64

	wordIndex := wordIndex(bitIndex)
	if wordIndex >= b.wordInUse {
		wordIndex = b.wordInUse - 1
		word = b.words[wordIndex]
	} else {
		// 找到bitIndex所在的word，将bitIndex位之后的位置清0
		word = b.words[wordIndex] & (WORK_MASK >> (WORD_SIZE - 1 - bitIndex&(WORD_SIZE-1)))
	}

	for {
		if word != 0 {
			// 当前word至少上有一个不为1的位
			return int(wordIndex+1)*WORD_SIZE - 1 - bits.LeadingZeros64(word)
		}

		if wordIndex == 0 {
			return -1
		}

		// 找上一个word
		wordIndex--
		word = b.words[wordIndex]
	}
}

// PrevClearBit 找到bitIndex位之前的第一个设置为0的位置
func (b *BitSet) PrevClearBit(bitIndex uint) int {
	var word uint64

	wordIndex := wordIndex(bitIndex)
	if wordIndex >= b.wordInUse {
		wordIndex = b.wordInUse - 1
		word = ^b.words[wordIndex]
	} else {
		// 找到bitIndex所在的word，取反，将bitIndex位之后的位置清0
		word = (^b.words[wordIndex]) & (WORK_MASK >> (WORD_SIZE - 1 - bitIndex&(WORD_SIZE-1)))
	}

	for {
		if word != 0 {
			// 当前word至少上有一个不为1的位
			return int(wordIndex+1)*WORD_SIZE - 1 - bits.LeadingZeros64(word)
		}

		if wordIndex == 0 {
			return -1
		}

		// 找上一个word
		wordIndex--
		word = ^b.words[wordIndex]
	}
}

// NextSetBit 找到从bitIndex位开始的第一个设置为1的位置
func (b *BitSet) NextSetBit(bitIndex uint) int {
	wordIndex := wordIndex(bitIndex)
	if wordIndex >= b.wordInUse {
		return -1
	}

	// 找到bitIndex所在的word，将bitIndex位之前的位置清0
	word := b.words[wordIndex] & (WORK_MASK << (bitIndex & (WORD_SIZE - 1)))
	for {
		if word != 0 {
			// 当前word至少上有一个不为1的位
			return int(wordIndex)*WORD_SIZE + bits.TrailingZeros64(word)
		}

		// 找下一个word
		wordIndex++
		if wordIndex == b.wordInUse {
			return -1
		}

		word = b.words[wordIndex]
	}
}

// NextClearBit 找到从bitIndex位开始的第一个设置为0的位置
func (b *BitSet) NextClearBit(bitIndex uint) int {
	wordIndex := wordIndex(bitIndex)
	if wordIndex >= b.wordInUse {
		return -1
	}

	// 找到bitIndex所在的word，取反，将bitIndex位之前的位置清0
	word := (^b.words[wordIndex]) & (WORK_MASK << (bitIndex & (WORD_SIZE - 1)))
	for {
		if word != 0 {
			// 当前word至少上有一个不为1的位
			return int(wordIndex)*WORD_SIZE + bits.TrailingZeros64(word)
		}

		// 找下一个word
		wordIndex++
		if wordIndex == b.wordInUse {
			return -1
		}

		word = ^b.words[wordIndex]
	}
}

// And 与
func (b *BitSet) And(s *BitSet) {
	for b.wordInUse > s.wordInUse {
		b.wordInUse--
		b.words[b.wordInUse] = 0
	}

	for i := uint(0); i < b.wordInUse; i++ {
		b.words[i] &= s.words[i]
	}
}

// Or 或
func (b *BitSet) Or(s *BitSet) {
	if b.wordInUse < s.wordInUse {
		b.grown(s.wordInUse)
		b.wordInUse = s.wordInUse
	}

	for i := uint(0); i < b.wordInUse; i++ {
		b.words[i] |= s.words[i]
	}
}

// Xor 异或
func (b *BitSet) Xor(s *BitSet) {
	if b.wordInUse < s.wordInUse {
		b.grown(s.wordInUse)
		b.wordInUse = s.wordInUse
	}

	for i := uint(0); i < b.wordInUse; i++ {
		b.words[i] ^= s.words[i]
	}
}

// AndNot 清除b在对应的s中已经设置为1的位
func (b *BitSet) AndNot(s *BitSet) {
	iMin := utils.Min(b.wordInUse, s.wordInUse)

	for i := uint(0); i < iMin; i++ {
		b.words[i] &^= s.words[i]
	}
}

// String String
func (b *BitSet) String() string {
	buffer := bytes.NewBufferString("{")

	i := b.NextSetBit(0)
	if i == -1 {
		buffer.WriteString("}")
		return buffer.String()
	}

	for {
		buffer.WriteString(strconv.FormatInt(int64(i), 10))
		i = b.NextSetBit(uint(i + 1))
		if i == -1 {
			buffer.WriteString("}")
			return buffer.String()
		}

		buffer.WriteString(", ")
	}
}

// grown 扩容
func (b *BitSet) grown(wordsRequired uint) {
	nWords := uint(len(b.words))

	if nWords < wordsRequired {
		size := utils.Max(2*nWords, wordsRequired)
		newWords := make([]uint64, size)
		copy(newWords, b.words)
		b.words = newWords
	}
}

// needWords 计算需要多少个word
func needWords(nbits int) int {
	if nbits < 0 {
		return 0
	}

	return (nbits + WORD_SIZE - 1) >> WORD_SHIFT
}

// wordIndex 找到word的索引
func wordIndex(bitIndex uint) uint {
	return bitIndex >> WORD_SHIFT
}
