package repository

import (
	"cargo-m/internal/common"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var dataSource *DataSource

type IRepository[T any] interface {
	List() ([]T, *gorm.DB)
	Create(entity *T) error
	Update(entity *T) error
	DeleteById(id uint) error
	Delete(query interface{}, args ...interface{}) error
	Page(sortParams *string, page common.Page[T]) (common.Page[T], error)
	One(sortParams *string, query interface{}, args ...interface{}) (T, error)
	SelectList(sortParams *string, query interface{}, args ...interface{}) ([]T, error)
	SelectPage(sortParams *string, page common.Page[T], query interface{}, args ...interface{}) (common.Page[T], error)

	Where(query interface{}, args ...interface{}) *gorm.DB

	GetEntity() *T
	DefaultSort() string
	InitializeRepository()
	GetConnection() *gorm.DB
	GetSort(sortParams *string) string
}

// BaseRepository 定于一 ----------------------------------------------------
type BaseRepository[T any] struct {
	dataSource *DataSource
}

func (br *BaseRepository[T]) DefaultSort() string {
	return "sort asc"
}

func (br *BaseRepository[T]) GetSort(sortParams *string) string {
	if sortParams != nil && *sortParams != "" {
		return *sortParams
	}
	return br.DefaultSort()
}

func (br *BaseRepository[T]) Create(entity *T) error {
	return br.GetConnection().Create(entity).Error
}

func (br *BaseRepository[T]) Update(entity *T) error {
	return br.GetConnection().Save(entity).Error
}

func (br *BaseRepository[T]) Delete(query interface{}, args ...interface{}) error {
	return br.Where(query, args...).Update("Valid", 0).Error
}

func (br *BaseRepository[T]) DeleteById(id uint) error {
	return br.Where("`id` = ?", id).Update("Valid", 0).Error
}

// List 获取说有有效数据, 使用默认排序条件
func (br *BaseRepository[T]) List() ([]T, *gorm.DB) {
	var res []T
	tx := br.Where("`Valid` = ?", 1).Order(br.DefaultSort()).Find(&res)
	return res, tx
}

// One 按照条件查询首条数据
func (br *BaseRepository[T]) One(sortParams *string, query interface{}, args ...interface{}) (T, error) {
	var res T
	sort := br.GetSort(sortParams)
	tx := br.Where(query, args...).Order(sort).First(&res)
	if tx.Error != nil {
		return res, tx.Error
	}
	return res, nil
}

// SelectList 按照条件查询列表数据
func (br *BaseRepository[T]) SelectList(sortParams *string, query interface{}, args ...interface{}) ([]T, error) {
	var res []T
	sort := br.GetSort(sortParams)
	tx := br.Where(query, args...).Order(sort).Find(&res)
	if tx.Error != nil {
		return res, tx.Error
	}
	return res, nil
}

// Page 分页查询所有有效数据
func (br *BaseRepository[T]) Page(sortParams *string, page common.Page[T]) (common.Page[T], error) {
	var res []T
	sort := br.GetSort(sortParams)
	count := int64(0)
	query := br.Where("`Valid` = ?", 1)
	query.Count(&count)
	page.Total = &count
	first := page.GetFirst()
	// 分页数据大于当前查询结果最大数据
	if page.Total != nil && *page.Total < int64(first) {
		page.PageNo = 1
		first = page.GetFirst()
	}
	tx := query.Order(sort).Offset(first).Limit(page.PageSize).Find(&res)
	if tx.Error != nil {
		return page, tx.Error
	}
	page.Result = &res
	return page, nil
}

// SelectPage 按照条件分页查询列表数据
func (br *BaseRepository[T]) SelectPage(sortParams *string, page common.Page[T], query interface{}, args ...interface{}) (common.Page[T], error) {
	var res []T
	sort := br.GetSort(sortParams)
	count := int64(0)
	br.Where(query, args...).Count(&count)
	page.Total = &count
	first := page.GetFirst()
	// 分页数据大于当前查询结果最大数据
	if page.Total != nil && *page.Total < int64(first) {
		page.PageNo = 1
		first = page.GetFirst()
	}
	tx := br.Where(query, args...).Order(sort).Offset(first).Limit(page.PageSize).Find(&res)
	if tx.Error != nil {
		return page, tx.Error
	}
	page.Result = &res
	return page, nil
}

func (br *BaseRepository[T]) Where(query interface{}, args ...interface{}) *gorm.DB {
	connection := br.GetConnection()
	return connection.Model(br.GetEntity()).Where(query, args...)
}

func (br *BaseRepository[T]) GetEntity() *T {
	return new(T)
}

func (br *BaseRepository[T]) InitializeRepository() {
	if dataSource != nil {
		err := br.GetConnection().AutoMigrate(br.GetEntity())
		if err != nil {
			panic(err)
		}
	}
}

func (br *BaseRepository[T]) GetConnection() *gorm.DB {
	dataSource := br.dataSource
	if dataSource == nil {
		panic("dataSource is nil")
	}
	return dataSource.GetConnection()
}

// DataSource --------------------------------------

type DataSource struct {
	db *gorm.DB
}

func (ds DataSource) GetConnection() *gorm.DB {
	if ds.db == nil {
		panic("database not initialized")
	}
	return ds.db
}

func NewDataSource() *DataSource {
	if dataSource != nil {
		return dataSource
	}
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Panicln("数据库连接失败" + err.Error())
	}
	dataSource = &DataSource{db: db}
	return dataSource
}
