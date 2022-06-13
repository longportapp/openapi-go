package quote

import quotev1 "github.com/longbridgeapp/openapi-protobufs/gen/go/quote"


type SubFlag uint8

const (
	SUBFLAG_QUOTE SubFlag = 0x1
	SUBFLAG_DEPTH SubFlag = 0x2
	SUBFLAG_BROKER SubFlag = 0x4
	SUBFLAG_TRADE SubFlag = 0x8
)

func toSubTypes(flags []SubFlag) []quotev1.SubType{
    subTypes := make([]quotev1.SubType, 0,len(flags))
	for flag, _ := range subTypes {
		subTypes = append(subTypes, quotev1.SubType(flag))
	}
	return subTypes
}

func toSubFlags(flags []quotev1.SubType) []SubFlag{
    subTypes := make([]SubFlag, 0,len(flags))
	for flag, _ := range subTypes {
		subTypes = append(subTypes, SubFlag(flag))
	}
	return subTypes
}
