package util

import (
	"time"

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
		{
			SrcType: time.Time{},
			DstType: copier.String,
			Fn: func(src interface{}) (interface{}, error) {
				value, ok := src.(time.Time)

				if !ok {
					return nil, errors.New("convert time to string, but src type not matching")
				}
				return FormatDate(&value), nil
			},
		},
	},
}

func Copy(toValue interface{}, fromValue interface{}) error {
	return copier.CopyWithOption(toValue, fromValue, opt)
}
