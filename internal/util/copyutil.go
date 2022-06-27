package util

import (
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

var opt = copier.Option{
	IgnoreEmpty: true,
	DeepCopy:    false,
	Converters: []copier.TypeConverter{
		{
			SrcType: copier.String,
			DstType: decimal.Decimal{},
			Fn: func(src interface{}) (interface{}, error) {
				value, ok := src.(string)

				if !ok {
					return nil, errors.New("convert string to decimal, but src type not matching")
				}

				return decimal.NewFromString(value)
			},
		},
	},
}

func Copy(toValue interface{}, fromValue interface{}) error {
	return copier.CopyWithOption(toValue, fromValue, opt)
}
