package models

import (
	"context"
	"github.com/uptrace/bun"
)

type Project struct {
	ID       int64   `bun:",pk" json:"id"`
	Name     string  `json:"name"`
	Template string  `json:"template"`
	Pools    []*Pool `bun:"rel:has-many,join:id=project_id" json:"pools"`
}

func GetProjects(db *bun.DB) ([]*Project, error) {
	ctx := context.Background()
	projects := make([]*Project, 0)

	err := db.NewSelect().Model(&projects).Scan(ctx)

	return projects, err
}

func CreateProject(db *bun.DB, req *Project) (*Project, error) {
	ctx := context.Background()
	_, err := db.NewInsert().Model(req).Exec(ctx)
	if err != nil {
		return nil, err
	}

	project := &Project{}

	err = db.NewSelect().
		Model(project).
		Where("project.id = ?", req.ID).
		Scan(ctx)

	return project, err
}

func GetProject(db *bun.DB, projectId string) (*Project, error) {
	ctx := context.Background()
	projectWithPools := &Project{}

	err := db.NewSelect().
		Model(projectWithPools).
		Relation("Pools").
		Where("project.id = ?", projectId).
		Scan(ctx)

		// ColumnExpr("project.*").
		// ColumnExpr("pool.*").
		// Join("JOIN pools ON pool.project_id = project.id").

		// Limit(1).

	return projectWithPools, err
}
