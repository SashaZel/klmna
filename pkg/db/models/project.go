package models

import (
	"github.com/go-pg/pg/v10"
)

type Project struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Template string `json:"template"`
}

func GetProjects(db *pg.DB) ([]*Project, error) {
	projects := make([]*Project, 0)

	err := db.Model(&projects).Select()

	return projects, err
}

func CreateProject(db *pg.DB, req *Project) (*Project, error) {
	_, err := db.Model(req).Insert()
	if err != nil {
		return nil, err
	}

	project := &Project{}

	err = db.Model(project).Where("project.id = ?", req.ID).Select()

	return project, err
}

func GetProject(db *pg.DB, projectId string) (*Project, error) {
	project := &Project{}

	err := db.Model(project).Where("project.id = ?", projectId).Select()

	return project, err
}
