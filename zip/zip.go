package zip

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg zip -path Ext,Zip

type engine byte

const (
	sevenZipEngine engine = iota
	goZipEngine
)

type state struct {
	folder string
	target string
	engine engine
}

type Option func(*state)

func With7Zip() Option {
	return func(s *state) {
		s.engine = sevenZipEngine
	}
}

func CompressFolder(folder string, target string, opts ...Option) error {
	s := state{
		folder: folder,
		target: target,
		engine: goZipEngine,
	}

	for _, opt := range opts {
		opt(&s)
	}

	switch s.engine {
	case goZipEngine:
		return goZipFolder(s.folder, s.target)
	case sevenZipEngine:
		return sevenZipFolder(s.folder, s.target)
	default:
		return newError("unknown zip engine: ", s.engine)
	}
}
