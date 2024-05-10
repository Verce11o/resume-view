package pagination

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

func DecodeCursor(encodedCursor string) (time.Time, uuid.UUID, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return time.Time{}, [16]byte{}, err
	}

	arrStr := strings.Split(string(byt), ",")
	if len(arrStr) != 2 {
		return time.Time{}, [16]byte{}, err
	}

	res, err := time.Parse(time.RFC3339Nano, arrStr[0])
	if err != nil {
		return time.Time{}, [16]byte{}, err
	}

	viewID, err := uuid.Parse(arrStr[1])
	if err != nil {
		return time.Time{}, [16]byte{}, err
	}

	return res, viewID, nil
}

func EncodeCursor(t time.Time, uuid string) string {
	key := fmt.Sprintf("%s,%s", t.Format(time.RFC3339Nano), uuid)
	return base64.StdEncoding.EncodeToString([]byte(key))
}
