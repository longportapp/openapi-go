package util

import "github.com/shopspring/decimal"

func ParseDecimal(value string) (res *decimal.Decimal, err error) {
	if value == "" {
		return nil, nil
	}
	*res, err = decimal.NewFromString(value)
	return
}
