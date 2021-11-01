package configure

import (
	"wiki_bit/app/model/configure"
)

func GetInfo() (configure configure.Configure, err error) {

	if configure, err = configure.FindOne(); err != nil {
		return configure, err
	}

	return configure, nil
}
