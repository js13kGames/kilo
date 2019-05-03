#include "textflag.h"
#include "funcdata.h"

#define LOOKUP_VEC_REG X5
#define SHUFFLE_VEC_REG X6
#define INDEX_1_MASK_REG X7
#define INDEX_1_SHIFT_REG X8
#define INDEX_2_MASK_REG X9
#define INDEX_2_SHIFT_REG X10
#define SUB51_REG X11
#define CMP26_REG X12

// lookupVec (for the URL-safe base64 alphabet) corresponds to:
// [16]int8{
//	65, 71, -4, -4,
//  -4, -4, -4, -4,
//	-4, -4, -4, -4,
//  -17, 32, 0,  0,
// }
DATA lookupVec<>+0(SB)/8, $0xFCFCFCFCFCFC4741
DATA lookupVec<>+8(SB)/8, $0xEF200000FCFCFCFC
GLOBL lookupVec<>(SB), (NOPTR+RODATA), $16

DATA shuffleVec<>+0(SB)/8, $0x0405030401020001
DATA shuffleVec<>+8(SB)/8, $0x0A0B090A07080607
GLOBL shuffleVec<>(SB), (NOPTR+RODATA), $16

DATA index1Mask<>+0(SB)/8, $0x0fc0fc000fc0fc00
DATA index1Mask<>+8(SB)/8, $0x0fc0fc000fc0fc00
GLOBL index1Mask<>(SB), (NOPTR+RODATA), $16

DATA index1Shift<>+0(SB)/8, $0x0400004004000040
DATA index1Shift<>+8(SB)/8, $0x0400004004000040
GLOBL index1Shift<>(SB), (NOPTR+RODATA), $16

DATA index2Mask<>+0(SB)/8, $0x003f03f0003f03f0
DATA index2Mask<>+8(SB)/8, $0x003f03f0003f03f0
GLOBL index2Mask<>(SB), (NOPTR+RODATA), $16

DATA index2Shift<>+0(SB)/8, $0x0100001001000010
DATA index2Shift<>+8(SB)/8, $0x0100001001000010
GLOBL index2Shift<>(SB), (NOPTR+RODATA), $16

DATA sub51Mask<>+0(SB)/8, $0x3333333333333333
DATA sub51Mask<>+8(SB)/8, $0x3333333333333333
GLOBL sub51Mask<>(SB), (NOPTR+RODATA), $16

DATA cmp26Mask<>+0(SB)/8, $0x1919191919191919
DATA cmp26Mask<>+8(SB)/8, $0x1919191919191919
GLOBL cmp26Mask<>(SB), (NOPTR+RODATA), $16

// func encodeSSE3(dst []byte, src *Token)
TEXT ·encodeSSE3(SB), NOSPLIT, $0-32
	XORQ SI, SI
	XORQ DI, DI

    // Input.
	MOVQ dst+0(FP), R10
	MOVQ src+24(FP), BX

    // Constants.
	MOVOA shuffleVec<>+0(SB), SHUFFLE_VEC_REG
	MOVOA lookupVec<>+0(SB), LOOKUP_VEC_REG
	MOVOA index1Mask<>+0(SB), INDEX_1_MASK_REG
	MOVOA index1Shift<>+0(SB), INDEX_1_SHIFT_REG
	MOVOA index2Mask<>+0(SB), INDEX_2_MASK_REG
	MOVOA index2Shift<>+0(SB), INDEX_2_SHIFT_REG
	MOVOA sub51Mask<>+0(SB), SUB51_REG
	MOVOA cmp26Mask<>+0(SB), CMP26_REG

enc:
	MOVQ   0(BX)(SI*1), X0  // Quad into X0
	MOVD   8(BX)(SI*1), X1  // Double into X1
	PSLLDQ $8, X1           // Shift X1 left 8 times.
	POR    X1, X0           // X0 is now ---- | 3332 | 2211 | 1000

	PSHUFB SHUFFLE_VEC_REG, X0

	MOVOA   X0, X1
	PAND    INDEX_1_MASK_REG, X0
	PAND    INDEX_2_MASK_REG, X1
	PMULHUW INDEX_1_SHIFT_REG, X0
	PMULLW  INDEX_2_SHIFT_REG, X1
	POR     X1, X0

	MOVOA   X0, X1
	MOVOA   X0, X2
	PSUBUSB SUB51_REG, X1
	PCMPGTB CMP26_REG, X2
	PSUBSB  X2, X1

	MOVOA  LOOKUP_VEC_REG, X3
	PSHUFB X1, X3
	PADDSB X3, X0

    // Write to dst.
	MOVOU X0, (R10)(DI*1)
	ADDQ $16, DI
	ADDQ $12, SI

	CMPQ SI, $24
	JLT  enc

	RET
