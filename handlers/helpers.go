package handlers

import "database/sql"

func nullOrNil(s sql.NullString) any {
	if s.Valid {
		return s.String
	}
	return nil
}

func nullOrNilBool(nb sql.NullBool) any {
	if nb.Valid {
		return nb.Bool
	}
	return nil
}
