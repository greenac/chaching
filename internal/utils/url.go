package utils

import (
	genErr "github.com/greenac/chaching/internal/error"
)

func JoinUrl(base string, add string) (string, genErr.IGenError) {
	if len(base) == 0 {
		ge := genErr.GenError{}
		return "", ge.AddMsg("base url can not be empty string")
	}

	if base[len(base)-1] != '/' {
		if len(add) > 0 && add[0] != '/' {
			base += "/"
		}
	} else if len(add) > 0 && add[0] == '/' {
		add = add[1:]
	}

	return base + add, nil
}
