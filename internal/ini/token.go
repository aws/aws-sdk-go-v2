package ini

type LineToken interface {
	isLineToken()
}

type LineTokenProfile struct {
	Type string
	Name string
}

func (*LineTokenProfile) isLineToken() {}

type LineTokenProperty struct {
	Key   string
	Value string
}

func (*LineTokenProperty) isLineToken() {}

type LineTokenContinuation struct {
	Value string
}

func (*LineTokenContinuation) isLineToken() {}

type LineTokenSubProperty struct {
	Key   string
	Value string
}

func (*LineTokenSubProperty) isLineToken() {}
