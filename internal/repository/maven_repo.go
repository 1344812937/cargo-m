package repository

import "cargo-m/internal/model"

type MavenRepo struct {
	*BaseRepository[model.MavenArtifactModel]
}

func NewMavenRepo(dataSource *DataSource) *MavenRepo {
	instance := &MavenRepo{
		BaseRepository: &BaseRepository[model.MavenArtifactModel]{
			dataSource: dataSource,
		},
	}
	instance.InitializeRepository()
	return instance
}

func (repo *MavenRepo) FindAll() ([]model.MavenArtifactModel, error) {
	allData, tx := repo.List()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return allData, nil
}

func (repo *MavenRepo) GetByKey(key string) (*model.MavenArtifactModel, error) {
	list, err := repo.SelectList(nil, ` valid = 1 and key = ?`, key)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return &list[0], nil
	}
	return nil, nil
}

func (repo *MavenRepo) Save(data []*model.MavenArtifactModel) error {
	repo.GetConnection().CreateInBatches(data, 100)
	return nil
}
