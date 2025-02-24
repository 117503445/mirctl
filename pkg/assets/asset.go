package assets

import "errors"

func GetMirror(repo string, mirror string) (map[string]string, error) {
	s := []map[string]string{
		{
			"mirror": "ustc",
			"repo":   "alpine",
			"url":    "https://mirrors.ustc.edu.cn/alpine",
		},
	}
	for _, v := range s {
		if v["mirror"] == mirror && v["repo"] == repo {
			return v, nil
		}
	}
	return nil, errors.New("mirror not found")
}
