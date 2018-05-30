package paperfishGo

func (w Writer) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}
