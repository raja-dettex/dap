package dap

type Options struct {
	DbName  string
	Encoder *JSONEncoder
	Decoder *JSONDecoder
}
