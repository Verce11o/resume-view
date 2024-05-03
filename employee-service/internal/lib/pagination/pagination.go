package pagination

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

func DecodeCursor(encodedCursor string) (time.Time, uuid.UUID, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		return time.Time{}, uuid.Nil, err
	}

	res, err := time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}

	id, err := uuid.Parse(arrStr[1])
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}

	return res, id, nil
}

func EncodeCursor(t time.Time, uuid string) string {
	key := fmt.Sprintf("%s,%s", t.Format(time.RFC3339Nano), uuid)
	return base64.StdEncoding.EncodeToString([]byte(key))
}
