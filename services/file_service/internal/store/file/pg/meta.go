package pg

import (
	"github.com/alexey-dobry/fileshare/services/file_service/internal/domain/model"
)

func (r *Repository) Create(file *model.File) error {
	return r.db.Create(&file).Error
}

func (r *Repository) GetByID(id string) (*model.File, error) {
	file := &model.File{}

	result := r.db.Select("*").Where("uuid = ?", id).First(file)
	if result.Error != nil {
		return &model.File{}, result.Error
	}
	return file, nil
}

func (r *Repository) Delete(id string) error {
	result := r.db.Where("uuid = ?", id).Delete(model.File{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *Repository) ListByCourse(courseID string) ([]*model.File, error) {
	var f []model.File

	result := r.db.Where("course_id = ?", courseID).Find(&f)
	if result.Error != nil {
		return []*model.File{}, result.Error
	}

	files := make([]*model.File, 0)

	for _, file := range f {
		files = append(files, &file)
	}
	return files, nil
}

func (r *Repository) ListByGroup(groupID string) ([]*model.File, error) {
	var f []model.File

	result := r.db.Where("group_id = ?", groupID).Find(&f)
	if result.Error != nil {
		return []*model.File{}, result.Error
	}

	files := make([]*model.File, 0)

	for _, file := range f {
		files = append(files, &file)
	}
	return files, nil
}
