package helper

import (
	"encoding/base64"
	"encoding/json"
	"time"
)

type Cursor map[string]interface{}

type PaginationInfo struct {
	NextCursor string `json:"next_cursor"`
	PrevCursor string `json:"prev_cursor"`
}

func CreateCursor(id string, createdAt time.Time, pointedNext bool) Cursor {
	return Cursor{
		"id":          id,
		"created_at":  createdAt,
		"points_next": pointedNext,
	}
}

func encodeCursor(cursor Cursor) string {
	if len(cursor) == 0 {
		return ""
	}
	serializedCursor, err := json.Marshal(cursor)
	if err != nil {
		return ""
	}

	// ✅ Use URL-safe Base64 encoding (no +, /, or =)
	encodedCursor := base64.RawURLEncoding.EncodeToString(serializedCursor)
	return encodedCursor
}

func GeneratePager(next Cursor, prev Cursor) PaginationInfo {
	return PaginationInfo{
		NextCursor: encodeCursor(next),
		PrevCursor: encodeCursor(prev),
	}
}

func DecodeCursor(cursor string) (Cursor, error) {
	decodedCursor, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return nil, err
	}

	var cur Cursor
	if err := json.Unmarshal(decodedCursor, &cur); err != nil {
		return nil, err
	}

	return cur, nil
}
