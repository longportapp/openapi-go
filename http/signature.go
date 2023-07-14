package http

import (
	"context"
	nhttp "net/http"
	"strconv"
	"time"

	"github.com/longbridgeapp/openapi-go/internal/signer"
	"github.com/longbridgeapp/openapi-go/internal/util"
)

const (
	headerTimestamp = "x-timestamp"
)

var sign = &signer.Signer{}

func signature(req *nhttp.Request, secret string, body []byte) error {
	if v := req.Header.Get(headerTimestamp); v == "" {
		req.Header.Add(headerTimestamp, strconv.FormatInt(time.Now().UnixMilli(), 10))
	}

	req.Header.Add("x-api-signature", "HMAC-SHA256 SignedHeaders=authorization;x-api-key;x-timestamp")
	signstr, _, _, err := sign.Sign(context.Background(), util.UnsafeStringToBytes(secret), req, body)
	if err != nil {
		return err
	}
	req.Header.Set("x-api-signature", "HMAC-SHA256 SignedHeaders=authorization;x-api-key;x-timestamp, Signature="+signstr)
	return nil
}
