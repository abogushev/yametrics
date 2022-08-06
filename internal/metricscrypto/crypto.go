package metricscrypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"yametrics/internal/protocol"
)

func GetMetricSign(m protocol.Metrics, key string) string {
	switch m.MType {
	case protocol.COUNTER:
		return getSign(fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta), key)
	case protocol.GAUGE:
		return getSign(fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value), key)
	default:
		return ""
	}
}

func getSign(key, msg string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(msg))

	return fmt.Sprintf("%x", h.Sum(nil))
}
