package format

// EmptyEnDecodeFormatNode 这个节点一般用作链表头节点，用于自动化构建链表使用
type EmptyEnDecodeFormatNode struct {
	DefaultEnDeCodeFormat
	BytesDecodeFormatValid
}

func (c *EmptyEnDecodeFormatNode) ValidDecode(encode any) (bool, error) {
	return true, nil
}
func (c *EmptyEnDecodeFormatNode) ValidEncode(encode any) (bool, error) {
	return true, nil
}
func (c *EmptyEnDecodeFormatNode) Encode(obj any) (any, error) {
	return obj, nil
}
func (c *EmptyEnDecodeFormatNode) Decode(encode any, back any) (any, error) {
	return encode, nil
}
