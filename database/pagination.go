package database

import (
	"fmt"
	"net/url"
	"time"

	"github.com/MehmoodNadeemKhan1/Page_Pagination_API/helper"
	"gorm.io/gorm"
)

type PaginatedItem interface {
	GetID() string
	GetCreatedAt() time.Time
}

func GetPaginatedQuery(query *gorm.DB, pointNext bool, cursor string, sortOrder string) (*gorm.DB, bool, error) {
	if cursor != "" {
		// Step 1️⃣ — URL-unescape only (no padding manipulation)
		unescapedCursor, err := url.QueryUnescape(cursor)
		if err != nil {
			return nil, pointNext, fmt.Errorf("failed to unescape cursor: %w", err)
		}

		// Step 2️⃣ — Decode Base64 cursor safely (URL-safe)
		decodedCursor, err := helper.DecodeCursor(unescapedCursor)
		if err != nil {
			return nil, pointNext, fmt.Errorf("failed to decode base64 cursor: %w", err)
		}

		// Step 3️⃣ — Get correct operator and order
		pointNext = decodedCursor["points_next"] == true
		operator, order := getPaginatedOperator(pointNext, sortOrder)

		// Step 4️⃣ — Apply WHERE filter
		whereStr := fmt.Sprintf("(created_at %s ? OR (created_at = ? AND id %s ?))", operator, operator)
		query = query.Where(whereStr, decodedCursor["created_at"], decodedCursor["created_at"], decodedCursor["id"])

		// Step 5️⃣ — Update sorting if needed
		if order != "" {
			sortOrder = order
		}
	}

	// Step 6️⃣ — Apply final ordering
	query = query.Order("created_at " + sortOrder)
	return query, pointNext, nil
}

func getPaginatedOperator(pointsNext bool, sortOrder string) (string, string) {
	if pointsNext && sortOrder == "asc" {
		return ">", ""
	}
	if pointsNext && sortOrder == "desc" {
		return "<", ""
	}
	if !pointsNext && sortOrder == "asc" {
		return "<", "desc"
	}
	if !pointsNext && sortOrder == "desc" {
		return ">", "asc"
	}
	return "", ""
}

func CalculatePagination[T PaginatedItem](isFirstPage bool, hasPagination bool, limit int, items []T, pointNext bool) helper.PaginationInfo {
	nextCur := helper.Cursor{}
	prevCur := helper.Cursor{}

	if len(items) == 0 {
		return helper.PaginationInfo{}
	}

	// 🧱 Case 1: First Page (no prev)
	if isFirstPage {
		if hasPagination {
			nextCur = helper.CreateCursor(items[limit-1].GetID(), items[limit-1].GetCreatedAt(), true)
		}
		// prev is empty
		return helper.GeneratePager(nextCur, prevCur)
	}

	// 🧱 Case 2: Navigating Forward (Next clicked)
	if pointNext {
		// Always create previous cursor (since not first page)
		prevCur = helper.CreateCursor(items[0].GetID(), items[0].GetCreatedAt(), false)

		if hasPagination {
			// More pages ahead → show next
			nextCur = helper.CreateCursor(items[limit-1].GetID(), items[limit-1].GetCreatedAt(), true)
		} else {
			// Last page → no next
			nextCur = helper.Cursor{}
		}

		return helper.GeneratePager(nextCur, prevCur)
	}

	// 🧱 Case 3: Navigating Backward (Prev clicked)
	if !pointNext {
		// Always create next cursor because if you’re coming back from last page → you’re now in the middle
		nextCur = helper.CreateCursor(items[limit-1].GetID(), items[limit-1].GetCreatedAt(), true)

		// But if we reach first page (no more prev)
		if hasPagination {
			prevCur = helper.CreateCursor(items[0].GetID(), items[0].GetCreatedAt(), false)
		} else {
			// This means we reached the very first page again
			prevCur = helper.Cursor{}
		}

		return helper.GeneratePager(nextCur, prevCur)
	}

	return helper.GeneratePager(nextCur, prevCur)
}
