package cmd

type Parser struct {
	modules []string
	unzip   bool
	merge   bool
	srcZip  string
	destDir string
}

type ParserOptions struct {
	modules []string
	unzip   bool
	merge   bool
	srcZip  string
	destDir string
}
type ParserOption func(*ParserOptions)

func WithModules(modules ...string) ParserOption {
	return func(o *ParserOptions) {
		o.modules = modules
	}
}

func WithUnzip(unzip bool) ParserOption {
	return func(o *ParserOptions) {
		o.unzip = unzip
	}
}

func WithSrcZip(srcZip string) ParserOption {
	return func(o *ParserOptions) {
		o.srcZip = srcZip
	}
}

func WithDestDir(destDir string) ParserOption {
	return func(o *ParserOptions) {
		o.destDir = destDir
	}
}

func WithMerge(merge bool) ParserOption {
	return func(o *ParserOptions) {
		o.merge = merge
	}
}

func NewParser(opts ...ParserOption) *Parser {
	options := &ParserOptions{
		modules: []string{},
	}

	for _, opt := range opts {
		opt(options)
	}

	return &Parser{
		modules: options.modules,
	}
}

func (parser *Parser) Parse() error {
	if parser.unzip {
		// Unzip(parser.srcZip, parser.destDir)
	}
	return nil
}
