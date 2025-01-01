package utils

import (
	"errors"
	"strconv"
)

func ToInt(text string) (int, error) {
	number, err := strconv.Atoi(text)
	if err != nil {
		return 0, errors.New("não é permitido inserir caracteres, por favor insira um número")
	}
	return number, nil
}
