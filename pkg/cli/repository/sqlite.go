package repository

import (
	"context"
	"database/sql"
	"github.com/koooyooo/soi-go/pkg/model"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
)

func NewSQLiteRepository(ctx context.Context, basePath, bucket string) (Repository, error) {
	r := sqliteRepository{
		basePath:      basePath,
		currentBucket: bucket,
	}
	if err := r.reload(ctx); err != nil {
		return nil, err
	}
	return &r, nil
}

type sqliteRepository struct {
	basePath      string
	currentBucket string
	db            *sql.DB
}

func (r *sqliteRepository) reload(ctx context.Context) error {
	dbFilePath := filepath.Join(r.basePath, r.currentBucket+".db")
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return err
	}
	r.db = db
	return nil
}

func (r *sqliteRepository) Init(ctx context.Context) error {
	return initSchema(r.db)
}

func initSchema(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = db.Exec(ddl)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *sqliteRepository) LoadAll(ctx context.Context, bucket string) ([]*model.SoiData, error) {
	if r.currentBucket != bucket {
		r.reload(ctx)
		r.currentBucket = bucket
	}
	stmt, err := r.db.PrepareContext(ctx, selectAllSoisQuery)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.QueryContext(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}
	var sois []*model.SoiData
	for rows.Next() {
		soi, err := loadSoi(ctx, r.db, rows)
		if err != nil {
			return nil, err
		}
		sois = append(sois, soi)
	}
	return sois, nil
}

func (r *sqliteRepository) Load(ctx context.Context, bucket string, hash string) (*model.SoiData, bool, error) {
	if r.currentBucket != bucket {
		r.reload(ctx)
		r.currentBucket = bucket
	}
	stmt, err := r.db.PrepareContext(ctx, selectSoiQuery)
	if err != nil {
		return nil, false, err
	}
	rows, err := stmt.QueryContext(ctx, bucket, hash)
	if err != nil {
		return nil, false, err
	}
	if !rows.Next() {
		return nil, false, nil
	}
	soi, err := loadSoi(ctx, r.db, rows)
	if err != nil {
		return nil, false, err
	}
	return soi, true, nil
}

func loadSoi(ctx context.Context, db *sql.DB, rs *sql.Rows) (*model.SoiData, error) {
	soi := &model.SoiData{}
	rs.Scan(
		&soi.ID,
		&soi.Hash,
		&soi.Name,
		&soi.Title,
		&soi.Path,
		&soi.URI,
		&soi.Rate,
		&soi.NumViews,
		&soi.NumReads,
		&soi.Comprehension,
		&soi.CreatedAt,
	)
	// tags
	stmt, err := db.PrepareContext(ctx, selectTagQuery)
	if err != nil {
		return nil, err
	}
	trows, err := stmt.QueryContext(ctx, soi.ID)
	for trows.Next() {
		var tag string
		if err := trows.Scan(&tag); err != nil {
			return nil, err
		}
		soi.Tags = append(soi.Tags, tag)
	}
	// kvtags
	kvtStmt, err := db.PrepareContext(ctx, selectKVTagQuery)
	if err != nil {
		return nil, err
	}
	kvtRows, err := kvtStmt.QueryContext(ctx, soi.ID)
	for kvtRows.Next() {
		kvt := model.KVTag{}
		if err := kvtRows.Scan(&kvt.Key, &kvt.Value); err != nil {
			return nil, err
		}
		soi.KVTags = append(soi.KVTags, kvt)
	}
	// usage logs
	ulStmt, err := db.PrepareContext(ctx, selectUsageLogQuery)
	if err != nil {
		return nil, err
	}
	ulRows, err := ulStmt.QueryContext(ctx, soi.ID)
	for ulRows.Next() {
		ul := model.UsageLog{}
		var soiID int
		if err := ulRows.Scan(&soiID, &ul.Type, &ul.UsedAt); err != nil {
			return nil, err
		}
		soi.UsageLogs = append(soi.UsageLogs, ul)
	}
	// og
	ogStmt, err := db.PrepareContext(ctx, selectOGQuery)
	if err != nil {
		return nil, err
	}
	ogRows, err := ogStmt.QueryContext(ctx, soi.ID)
	if err != nil {
		return nil, err
	}
	if !ogRows.Next() {
		return nil, err
	}
	//var selectOGQuery = "select title, url, type, description, site_name from ogs where soi_id = ?"
	var ogID, soiID string
	if err := ogRows.Scan(&ogID, &soiID, &soi.OGTitle, &soi.OGURL, &soi.OGType, &soi.OGDescription, &soi.OGSiteName); err != nil {
		return nil, err
	}
	// og images
	ogiStmt, err := db.PrepareContext(ctx, selectOGImgQuery)
	if err != nil {
		return nil, err
	}
	ogiRows, err := ogiStmt.QueryContext(ctx, ogID)
	if err != nil {
		return nil, err
	}
	for ogiRows.Next() {
		ogi := model.OGImage{}
		//var selectOGImgQuery = "select url, secure_url, type, width, height, alt from og_imgs where og_id = ?"
		var imgID, ogID string
		if err := ogiRows.Scan(&imgID, &ogID, &ogi.URL, &ogi.SecureURL, &ogi.Type, &ogi.Width, &ogi.Height, &ogi.Alt); err != nil {
			return nil, err
		}
		soi.OGImages = append(soi.OGImages, ogi)
	}
	return soi, nil
}

func (r *sqliteRepository) Store(ctx context.Context, bucket string, soi *model.SoiData) error {
	if r.currentBucket != bucket {
		r.reload(ctx)
		r.currentBucket = bucket
	}
	tx, err := r.db.Begin()
	defer tx.Rollback()
	if err != nil {
		return err
	}
	// store sois
	tstmt, err := tx.PrepareContext(ctx, storeSoiQuery)
	if err != nil {
		return err
	}
	rslt, err := tstmt.ExecContext(
		ctx,
		bucket,
		soi.Hash,
		soi.Name,
		soi.Title,
		soi.Path,
		soi.URI,
		soi.Rate,
		soi.NumViews,
		soi.NumReads,
		soi.Comprehension,
		soi.CreatedAt)
	if err != nil {
		return err
	}
	soiID, err := rslt.LastInsertId()
	if err != nil {
		return err
	}
	// store tags
	for _, tag := range soi.Tags {
		tstmt, err := tx.PrepareContext(ctx, storeTagQuery)
		if err != nil {
			return err
		}
		trslt, err := tstmt.ExecContext(ctx, tag)
		if err != nil {
			return err
		}
		// rel
		ststmt, err := tx.PrepareContext(ctx, storeSoiTagsQuery)
		if err != nil {
			return err
		}
		tagID, err := trslt.LastInsertId()
		if err != nil {
			return err
		}
		ststmt.ExecContext(ctx, soiID, tagID)
	}
	// store kv tags
	for _, kvTag := range soi.KVTags {
		kvtstmt, err := tx.PrepareContext(ctx, storeKVTagQuery)
		if err != nil {
			return err
		}
		kvtrslt, err := kvtstmt.ExecContext(ctx, kvTag.Key, kvTag.Value)
		if err != nil {
			return err
		}
		tagID, err := kvtrslt.LastInsertId()
		if err != nil {
			return err
		}
		// rel
		skvtstmt, err := tx.PrepareContext(ctx, storeSoiKVTagsQuery)
		if err != nil {
			return err
		}
		skvtstmt.ExecContext(ctx, soiID, tagID)
	}
	// store usage logs
	for _, usageLog := range soi.UsageLogs {
		ulstmt, err := tx.PrepareContext(ctx, storeUsageLogQuery)
		if err != nil {
			return err
		}
		ulstmt.ExecContext(ctx, soiID, usageLog.Type, usageLog.UsedAt)
	}
	// store ogs
	ogstmt, err := tx.PrepareContext(ctx, storeOGQuery)
	if err != nil {
		return err
	}
	ogstmtRslt, err := ogstmt.ExecContext(ctx, soiID, soi.OGTitle, soi.OGURL, soi.OGType, soi.OGDescription, soi.OGSiteName)
	if err != nil {
		return err
	}
	ogID, err := ogstmtRslt.LastInsertId()
	if err != nil {
		return err
	}
	for _, img := range soi.OGImages {
		imgstmt, err := tx.PrepareContext(ctx, storeOGImgQuery)
		if err != nil {
			return err
		}
		imgstmt.ExecContext(ctx, ogID, img.URL, img.SecureURL, img.Type, img.Width, img.Height, img.Alt)
	}
	return tx.Commit()
}

func (r *sqliteRepository) Exists(ctx context.Context, bucket string, hash string) (bool, error) {
	if r.currentBucket != bucket {
		r.reload(ctx)
		r.currentBucket = bucket
	}
	_, ok, err := r.Load(ctx, bucket, hash)
	if err != nil {
		return ok, err
	}
	return ok, nil
}

func (r *sqliteRepository) Remove(ctx context.Context, bucket string, hash string) error {
	if r.currentBucket != bucket {
		r.reload(ctx)
		r.currentBucket = bucket
	}
	//TODO implement me
	panic("implement me")
}
