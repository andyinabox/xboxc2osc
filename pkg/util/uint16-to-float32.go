package util

func Uint16ToFloat32(v uint16) float32 {
	return float32(v) / float32(65536-1) // -1 because of zero index
}
