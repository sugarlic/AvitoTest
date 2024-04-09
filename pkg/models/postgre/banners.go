package postgre

import (
	"database/sql"
	"encoding/json"

	"github.com/avitoTest/pkg/models"
)

type BannerModel struct {
	DB *sql.DB
}

// Insert - Метод для создание нового банера в базе дынных.
func (m *BannerModel) Insert(banner models.Banner) error {
	stmt_banner := `INSERT INTO Banners (Content, FeatureID, IsActive, LastModified)
	VALUES ($1, $2, $3, NOW())`
	stmt_tags := `INSERT INTO BannerTags (BannerID, TagID)
	VALUES ($1, $2)`
	stmt_feature := `INSERT INTO BannerFeatures (BannerID, FeatureID)
	VALUES ($1, $2)`

	contentJSON, err := json.Marshal(banner.Content)
	if err != nil {
		panic(err.Error())
	}

	_, err = m.DB.Exec(stmt_banner, string(contentJSON), banner.Feature_id, banner.Is_active)
	if err != nil {
		return err
	}

	lastInsertID, err := m.LastInsertId()
	if err != nil {
		return err
	}

	for _, tag_id := range banner.Tags {
		_, err = m.DB.Exec(stmt_tags, lastInsertID, tag_id)
		if err != nil {
			return err
		}
	}

	_, err = m.DB.Exec(stmt_feature, lastInsertID, banner.Feature_id)
	if err != nil {
		return err
	}

	return nil
}

func (m *BannerModel) Get(tag_id, feature_id int) error {
	stmt := `SELECT b.ID, b.Content
	FROM Banners b
	JOIN BannerTags bt ON b.ID = bt.BannerID
	JOIN BannerFeatures bf ON b.ID = bf.BannerID
	WHERE bt.TagID = $1 AND bf.FeatureID = $2
	`

	_, err := m.DB.Exec(stmt, tag_id, feature_id)
	if err != nil {
		return err
	}

	return nil
}

func (m *BannerModel) LastInsertId() (int64, error) {
	var lastInsertID int64

	err := m.DB.QueryRow("SELECT LASTVAL()").Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}
