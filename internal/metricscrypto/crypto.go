package metricscrypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"yametrics/internal/protocol"
)

func GetMetricSign(m protocol.Metrics, key string) string {
	var data string
	switch m.MType {
	case protocol.COUNTER:
		data = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	case protocol.GAUGE:
		data = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	default:
		panic("key is not defined")
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}
