package format

// EmptyEnDecodeFormatNode 这个节点一般用作链表头节点，用于自动化构建链表使用
type EmptyEnDecodeFormatNode struct {
	DefaultEnDeCodeFormat
	BytesDecodeFormatValid
}

func (c *EmptyEnDecodeFormatNode) ValidDecode(encode any) (bool, error) {
	return c.BytesDecodeFormatValid.ValidDecode(encode)
}
func (c *EmptyEnDecodeFormatNode) Encode(obj any) (any, error) {
	return nil, nil
}
func (c *EmptyEnDecodeFormatNode) Decode(encode any, back any) (any, error) {
	return nil, nil
}
