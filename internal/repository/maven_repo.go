package repository

import "cargo-m/internal/model"

type MavenRepo struct{}

func NewMavenRepo() *MavenRepo {
	return &MavenRepo{}
}

func (repo *MavenRepo) FindAll() ([]model.MavenArtifactModel, error) {
	var allData []model.MavenArtifactModel
	tx := Db.Model(&model.MavenArtifactModel{}).Where(` valid = 1`).Find(&allData)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return allData, nil
}

func (repo *MavenRepo) GetByKey(key string) (*model.MavenArtifactModel, error) {
	list := []*model.MavenArtifactModel{}
	tx := Db.Model(&model.MavenArtifactModel{}).Where(` valid = 1 and key = ?`, key).Find(&list)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

func (repo *MavenRepo) Save(data []*model.MavenArtifactModel) error {
	Db.Model(&model.MavenArtifactModel{}).CreateInBatches(data, 100)
	return nil
}
