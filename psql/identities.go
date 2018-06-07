package psql

import (
	"github.com/wvh/uuid"
)

// CreateIdentity creates a user and links the external id to our own uuid.
func (db *DB) CreateAndGetIdentity(id string) (uuid.UUID, error) {
	var newUuid uuid.UUID

	tx, err := db.Begin()
	if err != nil {
		return newUuid, err
	}
	defer tx.Rollback()

	provisionalUuid, err := uuid.NewUUID()
	if err != nil {
		return newUuid, err
	}

	err = tx.QueryRow(`
		WITH existing AS (
			SELECT uid FROM identities WHERE extid=$1
		), inserted AS (
			INSERT INTO users (uid, extid) VALUES ($2, $1)
			ON CONFLICT (extid) DO NOTHING
			RETURNING uid
		)
		SELECT uid
		FROM existing
		UNION ALL
		SELECT uid
		FROM inserted`,
		id, provisionalUuid).Scan(newUuid)

	if err != nil {
		return newUuid, handleError(err)
	}

	return newUuid, tx.Commit()
}
