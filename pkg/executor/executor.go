package executor

import "github.com/rs/zerolog/log"

var executors = map[string]Executor{
	"alpine": &alpineExecutor{},
	"pip":    &pipExecutor{},
	"go":     &goExecutor{},
	"arch":   &archExecutor{},
	"npm":    &npmExecutor{},
	"ubuntu": &ubuntuExecutor{},
	"debian": &debianExecutor{},
	"rust":   &rustExecutor{},
}

type Executor interface {
	// PreCheck checks if the executor could work
	PreCheck() bool

	Run() error
}

func PreCheck() (repos []string) {
	repos = []string{}
	for name, e := range executors {
		if e.PreCheck() {
			repos = append(repos, name)
		}
	}
	return
}

func Run(repo string, mirror string) error {
	supports := []string{}
	for name := range executors {
		supports = append(supports, name)
	}

	if e, ok := executors[repo]; ok {
		return e.Run()
	} else {
		log.Warn().Str("repo", repo).Strs("supports", supports).Msg("no executor found")
	}
	return nil
}
