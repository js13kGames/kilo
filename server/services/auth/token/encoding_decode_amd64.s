#include "textflag.h"
#include "funcdata.h"

#define ROLL_VEC_REG            X5
#define UNDERSCORE_MASK_REG     X6
#define UNDERSCORE_SHIFT_REG    X7
#define SHUFFLE_VEC_REG         X8
#define SHUFFLE_CONSTANT_0_REG  X9
#define SHUFFLE_CONSTANT_1_REG  X10

DATA underscoreMask<>+0x00(SB)/8, $0x5F5F5F5F5F5F5F5F  // Corresponds to ASCII underscore, "_"
DATA underscoreMask<>+0x08(SB)/8, $0x5F5F5F5F5F5F5F5F
DATA underscoreShift<>+0x00(SB)/8, $0x2121212121212121 // 33
DATA underscoreShift<>+0x08(SB)/8, $0x2121212121212121

GLOBL underscoreMask<>(SB), (NOPTR+RODATA), $16
GLOBL underscoreShift<>(SB), (NOPTR+RODATA), $16

DATA shuffleConstant0<>+0x00(SB)/8, $0x0140014001400140
DATA shuffleConstant0<>+0x08(SB)/8, $0x0140014001400140
DATA shuffleConstant1<>+0x00(SB)/8, $0x0001100000011000
DATA shuffleConstant1<>+0x08(SB)/8, $0x0001100000011000

GLOBL shuffleConstant0<>(SB), (NOPTR+RODATA), $16
GLOBL shuffleConstant1<>(SB), (NOPTR+RODATA), $16

// shuffleVec corresponds to:
// [16]int8{
//   	2,  1,  0,
//   	6,  5,  4,
//   	10, 9,  8,
//   	14, 13, 12,
//   	-1, -1, -1, -1,
//   }
//
DATA shuffleVec<>+0(SB)/8, $0x090A040506000102
DATA shuffleVec<>+8(SB)/8, $0xFFFFFFFF0C0D0E08
GLOBL shuffleVec<>(SB), (NOPTR+RODATA), $16

// rollVec (for the URL-safe base64 alphabet) corresponds to:
// [16]int8{
//   	0, 0, 17, 4,
//   	-65, -65, -71, -71,
//   	0, 0, 0, 0,
//   	0, 0, 0, 0,
//   }
DATA rollVec<>+0(SB)/8, $0xB9B9BFBF04110000
DATA rollVec<>+8(SB)/8, $0x0000000000000000
GLOBL rollVec<>(SB), (NOPTR+RODATA), $16

// Note that we're accepting a slice as src but only taking the first 8 bytes
// from the slice header.
// func decodeSSE3(dst *Token, src []byte)
TEXT ·decodeSSE3(SB), NOSPLIT, $0-16
	XORQ SI, SI
	XORQ DI, DI

    // Input.
	MOVQ dst+0(FP), R10
	MOVQ src+8(FP), BX

    // Constants.
	MOVOA underscoreMask<>+0(SB), UNDERSCORE_MASK_REG
	MOVOA underscoreShift<>+0(SB), UNDERSCORE_SHIFT_REG

	MOVOA rollVec<>+0(SB), ROLL_VEC_REG
	MOVOA shuffleVec<>+0(SB), SHUFFLE_VEC_REG
	MOVOA shuffleConstant0<>+0(SB), SHUFFLE_CONSTANT_0_REG
	MOVOA shuffleConstant1<>+0(SB), SHUFFLE_CONSTANT_1_REG

dec:
	MOVOA  0(BX)(SI*1), X0

	MOVOA X0, X1
    PSRLL $4, X1
	PAND  UNDERSCORE_MASK_REG, X1

    MOVOA   X0, X2
    PCMPEQB UNDERSCORE_MASK_REG, X2
    PAND    UNDERSCORE_SHIFT_REG, X2

	MOVOA  ROLL_VEC_REG, X3
	PSHUFB X1, X3
    PADDSB X3, X0
    PADDSB X2, X0

    PMADDUBSW SHUFFLE_CONSTANT_0_REG, X0
    PMADDWL   SHUFFLE_CONSTANT_1_REG, X0            // PMADDWD in AT&T.
    PSHUFB    SHUFFLE_VEC_REG, X0

    // Write to dst.
	MOVOU X0, (R10)(DI*1)
	ADDQ $12, DI
    ADDQ $16, SI

    CMPQ SI, $32
    JLT  dec

	RET
