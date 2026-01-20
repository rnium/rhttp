package codecs

import (
	"encoding/base64"
	"fmt"
)

func ToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func ToBase64Data(contentType string, data []byte) string {
	return fmt.Sprintf(
		"data:%s;base64,%s",
		contentType,
		ToBase64(data),
	)
}
