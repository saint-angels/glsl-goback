package mysql

import (
	"database/sql"
	// "errors"

	"saint-angels/shaderbox/pkg/models"
)

type ArtworkModel struct {
	DB *sql.DB
}

func (m *ArtworkModel) Insert() (int, error) {
	query := `INSERT INTO artworks (created)
				VALUES(UTC_TIMESTAMP())`
	result, err := m.DB.Exec(query)
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *ArtworkModel) GetOldestUnrendered() (*models.Artwork, error) {
    query := `SELECT id, created FROM artworks
    ORDER BY created ASC LIMIT 1
	WHERE rendered = FALSE`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	//Should come after err check as you can't .Close() nil
	defer rows.Close()

	artwork := &models.Artwork{}
	for rows.Next() {
		err = rows.Scan(&artwork.ID, &artwork.Created)
		if err != nil {
			return nil, err
		}
		break
	}
	//Get possible errors during the iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return artwork, nil
}
