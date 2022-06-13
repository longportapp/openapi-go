package http

import (
	"context"
	nhttp "net/http"
	"strconv"
	"time"

	"github.com/longbridgeapp/openapi-go/internal/util"
	"github.com/longbridgeapp/openapi-go/signer"
)
var sign = &signer.Signer{}

func signature(req *nhttp.Request, secret string, body []byte) error {
	req.Header.Add("X-Timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	req.Header.Add("X-Api-Signuature", "HMAC-SHA256 SignedHeaders=authorization;x-api-key;x-timestamp")
	signstr, _, _, err := sign.Sign(context.Background(), util.UnsafeStringToBytes(secret), req, body)
	if err != nil {
		return err
	}
	req.Header.Set("X-Api-Signuature", signstr)
	return nil
}
