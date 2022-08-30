package util

import "github.com/shopspring/decimal"

func ParseDecimal(value string) (res *decimal.Decimal, err error) {
	if value == "" {
		return nil, nil
	}
	*res, err = decimal.NewFromString(value)
	return
}

func Percent(d *decimal.Decimal) *decimal.Decimal {
	if d == nil {
		return nil
	}
	per := decimal.NewFromInt(100)
	res := d.Div(per)
	return &res
}
