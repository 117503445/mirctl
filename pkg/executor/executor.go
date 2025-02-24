package executor

var executors = map[string]Executor{
	"alpine": &alpineExecutor{},
}

type Executor interface {
	Run() error
}

func Run(repo string, mirror string) error {
	if e, ok := executors[repo]; ok {
		return e.Run()
	}
	return nil
}
