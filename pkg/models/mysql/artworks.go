package mysql

import (
	"database/sql"
	"errors"

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

func (m *ArtworkModel) GetArtForRender() (*models.Artwork, error) {
    query := `SELECT id, created FROM artworks
	WHERE rendered = FALSE AND rendering = FALSE
    ORDER BY created ASC LIMIT 1`


	row := m.DB.QueryRow(query)
	artwork := &models.Artwork{}
	err := row.Scan(&artwork.ID, &artwork.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	//If got the artwork successfully, then mark it as rendering
	updateQuery := `UPDATE artworks SET
					rendering = ?
					WHERE id = ?
					`

	_, err = m.DB.Exec(updateQuery, true, artwork.ID)
	if err != nil {
		return nil, err
	}

	return artwork, nil
}
