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
func (m *BannerModel) Insert(banner models.Banner) (int, error) {
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
		return 0, err
	}

	lastInsertID, err := m.LastInsertId()
	if err != nil {
		return 0, err
	}

	for _, tag_id := range banner.Tags {
		_, err = m.DB.Exec(stmt_tags, lastInsertID, tag_id)
		if err != nil {
			return 0, err
		}
	}

	_, err = m.DB.Exec(stmt_feature, lastInsertID, banner.FeatureId)
	if err != nil {
		return 0, err
	}

	return int(lastInsertID), nil
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

// добавить оффсет
func (m *BannerModel) GetList(limit int, token string) ([]*models.Banner, error) {
	stmt := `SELECT id, Content, FeatureID, IsActive, LastModified, CreatedAt 
	FROM banners ORDER BY id DESC LIMIT $1`
	stmt_tag := `SELECT tagid
	FROM bannertags
	WHERE bannerid = $1`

	rows, err := m.DB.Query(stmt, limit)
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

		rows_tags, err := m.DB.Query(stmt_tag, s.ID)
		if err != nil {
			return nil, err
		}

		// получение тэгов
		var tag_id int
		for rows_tags.Next() {
			err = rows_tags.Scan(&tag_id)
			if err != nil {
				return nil, err
			}
			s.Tags = append(s.Tags, tag_id)
		}
		if err = rows_tags.Err(); err != nil {
			return nil, err
		}

		if s.IsActive || token == "admin_token" {
			banners = append(banners, s)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return banners, nil
}

func (m *BannerModel) Update(banner models.Banner) error {
	stmt := `UPDATE banners
	SET featureid = $2,
	isactive = $3,
	content = $4
	WHERE id = $1;`
	stmt_tag := `UPDATE bannertags
	SET bannerid = $2
	WHERE tagid = $1;`

	_, err := m.DB.Exec(stmt, banner.ID, banner.FeatureId, banner.IsActive, banner.Content)
	if err != nil {
		return err
	}

	for _, tag := range banner.Tags {
		_, err = m.DB.Exec(stmt_tag, tag, banner.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *BannerModel) Delete(id int) error {
	stmt_banners := `DELETE FROM banners
	WHERE id = $1;`
	stmt_tags := `DELETE FROM bannertags
	WHERE bannerid = $1;`
	stmt_features := `DELETE FROM bannerfeatures
	WHERE id = $1;`

	_, err := m.DB.Exec(stmt_banners, id)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec(stmt_tags, id)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec(stmt_features, id)
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
