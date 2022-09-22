package app

func Generate(regex string, count int) ([]string, error) {
	var res []string

	tokenizer := NewTokenizer()
	parser := NewParser()
	tokens, err := tokenizer.Tokenize(regex)
	if err != nil {
		return nil, err
	}

	rootNode, err := parser.Parse(tokens)
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
