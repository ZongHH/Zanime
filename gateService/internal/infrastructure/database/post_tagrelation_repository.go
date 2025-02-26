package database

import (
	"context"
	"database/sql"
	"gateService/internal/domain/entity"
	"strings"
)

type PostTagRelationRepositoryImpl struct {
	db *sql.DB
}

func NewPostTagRelationRepositoryImpl(db *sql.DB) *PostTagRelationRepositoryImpl {
	return &PostTagRelationRepositoryImpl{db: db}
}

func (r *PostTagRelationRepositoryImpl) CreatePostTagRelation(ctx context.Context, postTagRelations []*entity.PostTagRelation) error {
	if len(postTagRelations) == 0 {
		return nil
	}

	query := `
		INSERT INTO post_tag_relations (post_id, tag_id) 
		VALUES `
	vals := []interface{}{}
	placeholders := make([]string, 0, len(postTagRelations))

	for _, relation := range postTagRelations {
		placeholders = append(placeholders, "(?, ?)")
		vals = append(vals, relation.PostID, relation.TagID)
	}

	query += strings.Join(placeholders, ",")
	_, err := r.db.ExecContext(ctx, query, vals...)
	return err
}

func (r *PostTagRelationRepositoryImpl) UpdatePostTagRelation(ctx context.Context, postTagRelation *entity.PostTagRelation) error {
	query := "UPDATE post_tag_relations SET post_id = ?, tag_id = ? WHERE post_id = ? AND tag_id = ?"
	_, err := r.db.ExecContext(ctx, query, postTagRelation.PostID, postTagRelation.TagID, postTagRelation.PostID, postTagRelation.TagID)
	return err
}

func (r *PostTagRelationRepositoryImpl) DeletePostTagRelation(ctx context.Context, postTagRelation *entity.PostTagRelation) error {
	query := "DELETE FROM post_tag_relations WHERE post_id = ? AND tag_id = ?"
	_, err := r.db.ExecContext(ctx, query, postTagRelation.PostID, postTagRelation.TagID)
	return err
}

func (r *PostTagRelationRepositoryImpl) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

func (r *PostTagRelationRepositoryImpl) CreatePostTagRelationTx(ctx context.Context, tx *sql.Tx, postTagRelations []*entity.PostTagRelation) error {
	if len(postTagRelations) == 0 {
		return nil
	}

	query := `
		INSERT INTO post_tag_relations (post_id, tag_id) 
		VALUES `
	vals := []interface{}{}
	placeholders := make([]string, 0, len(postTagRelations))

	for _, relation := range postTagRelations {
		placeholders = append(placeholders, "(?, ?)")
		vals = append(vals, relation.PostID, relation.TagID)
	}

	query += strings.Join(placeholders, ",")
	_, err := tx.ExecContext(ctx, query, vals...)
	return err
}

func (r *PostTagRelationRepositoryImpl) UpdatePostTagRelationTx(ctx context.Context, tx *sql.Tx, postTagRelation *entity.PostTagRelation) error {
	query := "UPDATE post_tag_relations SET post_id = ?, tag_id = ? WHERE post_id = ? AND tag_id = ?"
	_, err := tx.ExecContext(ctx, query, postTagRelation.PostID, postTagRelation.TagID, postTagRelation.PostID, postTagRelation.TagID)
	return err
}

func (r *PostTagRelationRepositoryImpl) DeletePostTagRelationTx(ctx context.Context, tx *sql.Tx, postTagRelation *entity.PostTagRelation) error {
	query := "DELETE FROM post_tag_relations WHERE post_id = ? AND tag_id = ?"
	_, err := tx.ExecContext(ctx, query, postTagRelation.PostID, postTagRelation.TagID)
	return err
}

func (r *PostTagRelationRepositoryImpl) GetPostTagRelationByID(ctx context.Context, postTagRelationID int64) (*entity.PostTagRelation, error) {
	query := "SELECT * FROM post_tag_relations WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, postTagRelationID)
	var postTagRelation entity.PostTagRelation
	err := row.Scan(&postTagRelation.PostID, &postTagRelation.TagID)
	return &postTagRelation, err
}
