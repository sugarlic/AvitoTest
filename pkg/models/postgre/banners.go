package postgre

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/avitoTest/pkg/models"
)

type BannerModel struct {
	DB *sql.DB
}

// Insert - Метод для создание нового банера в базе дынных.
func (m *BannerModel) Insert(banner models.Banner) error {
	stmt_banner := `INSERT INTO Banners (Content, FeatureID, IsActive, LastModified, CreatedAt)
	VALUES ($1, $2, $3, NOW(), NOW())`
	stmt_tags := `INSERT INTO BannerTags (BannerID, TagID)
	VALUES ($1, $2)`
	stmt_feature := `INSERT INTO BannerFeatures (BannerID, FeatureID)
	VALUES ($1, $2)`

	contentJSON, err := json.Marshal(banner.Content)
	if err != nil {
		panic(err.Error())
	}

	_, err = m.DB.Exec(stmt_banner, string(contentJSON), banner.FeatureId, banner.IsActive)
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

	_, err = m.DB.Exec(stmt_feature, lastInsertID, banner.FeatureId)
	if err != nil {
		return err
	}

	return nil
}

func (m *BannerModel) Get(tag_id, feature_id int) (*models.Banner, error) {
	stmt := `SELECT b.ID, b.Content
	FROM Banners b
	JOIN BannerTags bt ON b.ID = bt.BannerID
	JOIN BannerFeatures bf ON b.ID = bf.BannerID
	WHERE bt.TagID = $1 AND bf.FeatureID = $2
	`

	row := m.DB.QueryRow(stmt, tag_id, feature_id)

	s := &models.Banner{}

	err := row.Scan(&s.ID, &s.Content)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

// добавить получение массива тэгов
func (m *BannerModel) GetList() ([]*models.Banner, error) {
	stmt := `SELECT id, Content, FeatureID, IsActive, LastModified, CreatedAt 
	FROM banners ORDER BY id DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var banners []*models.Banner

	for rows.Next() {
		s := &models.Banner{}
		err = rows.Scan(&s.ID, &s.FeatureId, &s.Content, &s.IsActive, &s.CreatedAt, &s.LastModified)
		if err != nil {
			return nil, err
		}
		banners = append(banners, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

func (m *BannerModel) Delete(id int) error {
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
