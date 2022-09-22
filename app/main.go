package app

type GeneratorI interface {
	Generate(regex string, count int) ([]string, error)
}

type generator struct{}

func NewGenerator() GeneratorI {
	return generator{}
}

func (g generator) Generate(regex string, count int) ([]string, error) {
	var res []string

	tokenizer := NewTokenizer()
	parser := NewParser()
	tokens, err := tokenizer.Tokenize(regex)
	if err != nil {
		return nil, err
	}

	rootNode, err := parser.Parse(tokens, NewRootNode())
	if err != nil {
		return nil, err
	}

	for i := 0; i < count; i++ {
		value, err := rootNode.Generate()
		if err != nil {
			return nil, err
		}
		res = append(res, value)
	}
	return res, nil
}
