package links

import (
	"fmt"
	"math/big"
)

func UrlToId(b62 string) (string, error) {
	var i big.Int
	_, ok := i.SetString(b62, 62)
	if !ok {
		return "", fmt.Errorf("Can't parse %s as base62.", b62)
	}
	return i.Text(16), nil

}

func IdToUrl(id string) (string, error) {
	var i big.Int
	_, ok := i.SetString(id, 16)
	if !ok {
		return "", fmt.Errorf("Can't parse %s as hex.", id)
	}
	return i.Text(62), nil
}
