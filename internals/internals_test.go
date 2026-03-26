package internals

import "testing"

// ============================================================
// Slice 测试
// ============================================================

func TestDemoSliceStructure(t *testing.T) {
	DemoSliceStructure()
}

func TestDemoSliceSharing(t *testing.T) {
	DemoSliceSharing()
}

func TestDemoSliceGrowth(t *testing.T) {
	DemoSliceGrowth()
}

func TestDemoSliceCopy(t *testing.T) {
	DemoSliceCopy()
}

func TestDemoSliceAppend(t *testing.T) {
	DemoSliceAppend()
}

func TestDemoSliceExpression(t *testing.T) {
	DemoSliceExpression()
}

func TestDemoSliceNil(t *testing.T) {
	DemoSliceNil()
}

func TestDemoSliceMemoryLeak(t *testing.T) {
	DemoSliceMemoryLeak()
}

// ============================================================
// Map 测试
// ============================================================

func TestDemoMapStructure(t *testing.T) {
	DemoMapStructure()
}

func TestDemoMapHash(t *testing.T) {
	DemoMapHash()
}

func TestDemoMapLookup(t *testing.T) {
	DemoMapLookup()
}

func TestDemoMapInsert(t *testing.T) {
	DemoMapInsert()
}

func TestDemoMapDelete(t *testing.T) {
	DemoMapDelete()
}

func TestDemoMapGrow(t *testing.T) {
	DemoMapGrow()
}

func TestDemoMapIterate(t *testing.T) {
	DemoMapIterate()
}

func TestDemoMapConcurrency(t *testing.T) {
	DemoMapConcurrency()
}

func TestDemoMapKeyTypes(t *testing.T) {
	DemoMapKeyTypes()
}

func TestDemoMapMemory(t *testing.T) {
	DemoMapMemory()
}

// ============================================================
// Channel 测试
// ============================================================

func TestDemoChannelStructure(t *testing.T) {
	DemoChannelStructure()
}

func TestDemoChannelSend(t *testing.T) {
	DemoChannelSend()
}

func TestDemoChannelReceive(t *testing.T) {
	DemoChannelReceive()
}

func TestDemoChannelClose(t *testing.T) {
	DemoChannelClose()
}

func TestDemoChannelBuffered(t *testing.T) {
	DemoChannelBuffered()
}

func TestDemoChannelUnbuffered(t *testing.T) {
	DemoChannelUnbuffered()
}

func TestDemoChannelSelect(t *testing.T) {
	DemoChannelSelect()
}

func TestDemoChannelNil(t *testing.T) {
	DemoChannelNil()
}

func TestDemoChannelCommonPatterns(t *testing.T) {
	DemoChannelCommonPatterns()
}

// ============================================================
// Interface 测试
// ============================================================

func TestDemoInterfaceStructure(t *testing.T) {
	DemoInterfaceStructure()
}

func TestDemoInterfaceType(t *testing.T) {
	DemoInterfaceType()
}

func TestDemoInterfaceNil(t *testing.T) {
	DemoInterfaceNil()
}

func TestDemoInterfaceAssert(t *testing.T) {
	DemoInterfaceAssert()
}

func TestDemoInterfaceSwitch(t *testing.T) {
	DemoInterfaceSwitch()
}

func TestDemoInterfaceMethod(t *testing.T) {
	DemoInterfaceMethod()
}

func TestDemoInterfaceConversion(t *testing.T) {
	DemoInterfaceConversion()
}

func TestDemoInterfaceBoxing(t *testing.T) {
	DemoInterfaceBoxing()
}

func TestDemoInterfaceCompare(t *testing.T) {
	DemoInterfaceCompare()
}

func TestDemoInterfaceEscape(t *testing.T) {
	DemoInterfaceEscape()
}

// ============================================================
// String 测试
// ============================================================

func TestDemoStringStructure(t *testing.T) {
	DemoStringStructure()
}

func TestDemoStringImmutable(t *testing.T) {
	DemoStringImmutable()
}

func TestDemoStringBytes(t *testing.T) {
	DemoStringBytes()
}

func TestDemoStringConcat(t *testing.T) {
	DemoStringConcat()
}

func TestDemoStringRune(t *testing.T) {
	DemoStringRune()
}

func TestDemoStringIntern(t *testing.T) {
	DemoStringIntern()
}

func TestDemoStringCompare(t *testing.T) {
	DemoStringCompare()
}

func TestDemoStringSubstr(t *testing.T) {
	DemoStringSubstr()
}

func TestDemoStringSwitch(t *testing.T) {
	DemoStringSwitch()
}

func TestDemoStringUTF8(t *testing.T) {
	DemoStringUTF8()
}

func TestDemoStringPool(t *testing.T) {
	DemoStringPool()
}

func TestDemoStringReflection(t *testing.T) {
	DemoStringReflection()
}
