package service

import (
	"context"
	"fmt"

	"github.com/artorias501/blog-service/internal/domain/entity"
	"github.com/artorias501/blog-service/internal/domain/repository"
	"github.com/artorias501/blog-service/internal/domain/valueobject"
)

// Mock Tag Repository
type mockTagRepository struct {
	tags        map[string]*entity.Tag
	createFn    func(ctx context.Context, tag *entity.Tag) error
	getByIDFn   func(ctx context.Context, id valueobject.TagID) (*entity.Tag, error)
	getByNameFn func(ctx context.Context, name valueobject.TagName) (*entity.Tag, error)
	updateFn    func(ctx context.Context, tag *entity.Tag) error
	deleteFn    func(ctx context.Context, id valueobject.TagID) error
	listFn      func(ctx context.Context, params repository.ListParams) (*repository.TagListResult, error)
}

func newMockTagRepository() *mockTagRepository {
	return &mockTagRepository{
		tags: make(map[string]*entity.Tag),
	}
}

func (m *mockTagRepository) Create(ctx context.Context, tag *entity.Tag) error {
	if m.createFn != nil {
		return m.createFn(ctx, tag)
	}
	m.tags[tag.ID().String()] = tag
	return nil
}

func (m *mockTagRepository) GetByID(ctx context.Context, id valueobject.TagID) (*entity.Tag, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return m.tags[id.String()], nil
}

func (m *mockTagRepository) GetByName(ctx context.Context, name valueobject.TagName) (*entity.Tag, error) {
	if m.getByNameFn != nil {
		return m.getByNameFn(ctx, name)
	}
	for _, t := range m.tags {
		if t.Name().String() == name.String() {
			return t, nil
		}
	}
	return nil, nil
}

func (m *mockTagRepository) Update(ctx context.Context, tag *entity.Tag) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, tag)
	}
	m.tags[tag.ID().String()] = tag
	return nil
}

func (m *mockTagRepository) Delete(ctx context.Context, id valueobject.TagID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	delete(m.tags, id.String())
	return nil
}

func (m *mockTagRepository) List(ctx context.Context, params repository.ListParams) (*repository.TagListResult, error) {
	if m.listFn != nil {
		return m.listFn(ctx, params)
	}
	var tags []*entity.Tag
	for _, t := range m.tags {
		tags = append(tags, t)
	}
	return &repository.TagListResult{
		Items: tags,
		Total: int64(len(tags)),
	}, nil
}

func (m *mockTagRepository) ListByPostID(ctx context.Context, postID valueobject.PostID) ([]*entity.Tag, error) {
	return []*entity.Tag{}, nil
}

func (m *mockTagRepository) GetOrCreate(ctx context.Context, name valueobject.TagName) (*entity.Tag, error) {
	tag, _ := m.GetByName(ctx, name)
	if tag != nil {
		return tag, nil
	}
	tagName, _ := valueobject.NewTagName(name.String())
	newTag := entity.NewTag(tagName)
	m.Create(ctx, newTag)
	return newTag, nil
}

func (m *mockTagRepository) Search(ctx context.Context, keyword string, params repository.ListParams) (*repository.TagListResult, error) {
	return m.List(ctx, params)
}

// Mock Tag Cache Repository
type mockTagCacheRepository struct {
	cache        map[string]*entity.Tag
	listCache    map[string]*repository.TagListResult
	getFn        func(ctx context.Context, id valueobject.TagID) (*entity.Tag, error)
	setFn        func(ctx context.Context, tag *entity.Tag) error
	deleteFn     func(ctx context.Context, id valueobject.TagID) error
	getByNameFn  func(ctx context.Context, name valueobject.TagName) (*entity.Tag, error)
	setByNameFn  func(ctx context.Context, name valueobject.TagName, tag *entity.Tag) error
	getListFn    func(ctx context.Context, params repository.ListParams) (*repository.TagListResult, error)
	setListFn    func(ctx context.Context, params repository.ListParams, result *repository.TagListResult) error
	deleteListFn func(ctx context.Context) error
	invalidateFn func(ctx context.Context, id valueobject.TagID) error
}

func newMockTagCacheRepository() *mockTagCacheRepository {
	return &mockTagCacheRepository{
		cache:     make(map[string]*entity.Tag),
		listCache: make(map[string]*repository.TagListResult),
	}
}

func (m *mockTagCacheRepository) Get(ctx context.Context, id valueobject.TagID) (*entity.Tag, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return m.cache[id.String()], nil
}

func (m *mockTagCacheRepository) Set(ctx context.Context, tag *entity.Tag) error {
	if m.setFn != nil {
		return m.setFn(ctx, tag)
	}
	m.cache[tag.ID().String()] = tag
	return nil
}

func (m *mockTagCacheRepository) Delete(ctx context.Context, id valueobject.TagID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	delete(m.cache, id.String())
	return nil
}

func (m *mockTagCacheRepository) GetByName(ctx context.Context, name valueobject.TagName) (*entity.Tag, error) {
	if m.getByNameFn != nil {
		return m.getByNameFn(ctx, name)
	}
	return nil, nil
}

func (m *mockTagCacheRepository) SetByName(ctx context.Context, name valueobject.TagName, tag *entity.Tag) error {
	if m.setByNameFn != nil {
		return m.setByNameFn(ctx, name, tag)
	}
	return nil
}

func (m *mockTagCacheRepository) GetList(ctx context.Context, params repository.ListParams) (*repository.TagListResult, error) {
	if m.getListFn != nil {
		return m.getListFn(ctx, params)
	}
	key := fmt.Sprintf("tag:list:page:%d:size:%d", params.Page, params.PageSize)
	return m.listCache[key], nil
}

func (m *mockTagCacheRepository) SetList(ctx context.Context, params repository.ListParams, result *repository.TagListResult) error {
	if m.setListFn != nil {
		return m.setListFn(ctx, params, result)
	}
	key := fmt.Sprintf("tag:list:page:%d:size:%d", params.Page, params.PageSize)
	m.listCache[key] = result
	return nil
}

func (m *mockTagCacheRepository) DeleteList(ctx context.Context) error {
	if m.deleteListFn != nil {
		return m.deleteListFn(ctx)
	}
	m.listCache = make(map[string]*repository.TagListResult)
	return nil
}

func (m *mockTagCacheRepository) GetByPostID(ctx context.Context, postID valueobject.PostID) ([]*entity.Tag, error) {
	return nil, nil
}

func (m *mockTagCacheRepository) SetByPostID(ctx context.Context, postID valueobject.PostID, tags []*entity.Tag) error {
	return nil
}

func (m *mockTagCacheRepository) InvalidateTag(ctx context.Context, id valueobject.TagID) error {
	if m.invalidateFn != nil {
		return m.invalidateFn(ctx, id)
	}
	delete(m.cache, id.String())
	return nil
}

// Mock Comment Repository
type mockCommentRepository struct {
	comments         map[string]*entity.Comment
	createFn         func(ctx context.Context, comment *entity.Comment) error
	getByIDFn        func(ctx context.Context, id valueobject.CommentID) (*entity.Comment, error)
	updateFn         func(ctx context.Context, comment *entity.Comment) error
	deleteFn         func(ctx context.Context, id valueobject.CommentID) error
	listFn           func(ctx context.Context, params repository.ListParams) (*repository.CommentListResult, error)
	listByPostIDFn   func(ctx context.Context, postID valueobject.PostID, params repository.ListParams) (*repository.CommentListResult, error)
	countByPostIDFn  func(ctx context.Context, postID valueobject.PostID) (int64, error)
	approveFn        func(ctx context.Context, id valueobject.CommentID) error
	rejectFn         func(ctx context.Context, id valueobject.CommentID) error
	markAsSpamFn     func(ctx context.Context, id valueobject.CommentID) error
	deleteByPostIDFn func(ctx context.Context, postID valueobject.PostID) error
}

func newMockCommentRepository() *mockCommentRepository {
	return &mockCommentRepository{
		comments: make(map[string]*entity.Comment),
	}
}

func (m *mockCommentRepository) Create(ctx context.Context, comment *entity.Comment) error {
	if m.createFn != nil {
		return m.createFn(ctx, comment)
	}
	m.comments[comment.ID().String()] = comment
	return nil
}

func (m *mockCommentRepository) GetByID(ctx context.Context, id valueobject.CommentID) (*entity.Comment, error) {
	if m.getByIDFn != nil {
		return m.getByIDFn(ctx, id)
	}
	return m.comments[id.String()], nil
}

func (m *mockCommentRepository) Update(ctx context.Context, comment *entity.Comment) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, comment)
	}
	m.comments[comment.ID().String()] = comment
	return nil
}

func (m *mockCommentRepository) Delete(ctx context.Context, id valueobject.CommentID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	delete(m.comments, id.String())
	return nil
}

func (m *mockCommentRepository) List(ctx context.Context, params repository.ListParams) (*repository.CommentListResult, error) {
	if m.listFn != nil {
		return m.listFn(ctx, params)
	}
	var comments []*entity.Comment
	for _, c := range m.comments {
		comments = append(comments, c)
	}
	return &repository.CommentListResult{
		Items: comments,
		Total: int64(len(comments)),
	}, nil
}

func (m *mockCommentRepository) ListByPostID(ctx context.Context, postID valueobject.PostID, params repository.ListParams) (*repository.CommentListResult, error) {
	if m.listByPostIDFn != nil {
		return m.listByPostIDFn(ctx, postID, params)
	}
	return m.List(ctx, params)
}

func (m *mockCommentRepository) ListByStatus(ctx context.Context, status string, params repository.ListParams) (*repository.CommentListResult, error) {
	return m.List(ctx, params)
}

func (m *mockCommentRepository) CountByPostID(ctx context.Context, postID valueobject.PostID) (int64, error) {
	if m.countByPostIDFn != nil {
		return m.countByPostIDFn(ctx, postID)
	}
	return int64(len(m.comments)), nil
}

func (m *mockCommentRepository) Approve(ctx context.Context, id valueobject.CommentID) error {
	if m.approveFn != nil {
		return m.approveFn(ctx, id)
	}
	return nil
}

func (m *mockCommentRepository) Reject(ctx context.Context, id valueobject.CommentID) error {
	if m.rejectFn != nil {
		return m.rejectFn(ctx, id)
	}
	return nil
}

func (m *mockCommentRepository) MarkAsSpam(ctx context.Context, id valueobject.CommentID) error {
	if m.markAsSpamFn != nil {
		return m.markAsSpamFn(ctx, id)
	}
	return nil
}

func (m *mockCommentRepository) DeleteByPostID(ctx context.Context, postID valueobject.PostID) error {
	if m.deleteByPostIDFn != nil {
		return m.deleteByPostIDFn(ctx, postID)
	}
	return nil
}

// Mock Comment Cache Repository
type mockCommentCacheRepository struct {
	cache                map[string]*entity.Comment
	listCache            map[string]*repository.CommentListResult
	countCache           map[string]int64
	getFn                func(ctx context.Context, id valueobject.CommentID) (*entity.Comment, error)
	setFn                func(ctx context.Context, comment *entity.Comment) error
	deleteFn             func(ctx context.Context, id valueobject.CommentID) error
	getListByPostIDFn    func(ctx context.Context, postID valueobject.PostID, params repository.ListParams) (*repository.CommentListResult, error)
	setListByPostIDFn    func(ctx context.Context, postID valueobject.PostID, params repository.ListParams, result *repository.CommentListResult) error
	deleteListByPostIDFn func(ctx context.Context, postID valueobject.PostID) error
	getCountByPostIDFn   func(ctx context.Context, postID valueobject.PostID) (int64, error)
	setCountByPostIDFn   func(ctx context.Context, postID valueobject.PostID, count int64) error
	invalidateFn         func(ctx context.Context, id valueobject.CommentID) error
	invalidateByPostIDFn func(ctx context.Context, postID valueobject.PostID) error
}

func newMockCommentCacheRepository() *mockCommentCacheRepository {
	return &mockCommentCacheRepository{
		cache:      make(map[string]*entity.Comment),
		listCache:  make(map[string]*repository.CommentListResult),
		countCache: make(map[string]int64),
	}
}

func (m *mockCommentCacheRepository) Get(ctx context.Context, id valueobject.CommentID) (*entity.Comment, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return m.cache[id.String()], nil
}

func (m *mockCommentCacheRepository) Set(ctx context.Context, comment *entity.Comment) error {
	if m.setFn != nil {
		return m.setFn(ctx, comment)
	}
	m.cache[comment.ID().String()] = comment
	return nil
}

func (m *mockCommentCacheRepository) Delete(ctx context.Context, id valueobject.CommentID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, id)
	}
	delete(m.cache, id.String())
	return nil
}

func (m *mockCommentCacheRepository) GetListByPostID(ctx context.Context, postID valueobject.PostID, params repository.ListParams) (*repository.CommentListResult, error) {
	if m.getListByPostIDFn != nil {
		return m.getListByPostIDFn(ctx, postID, params)
	}
	return nil, nil
}

func (m *mockCommentCacheRepository) SetListByPostID(ctx context.Context, postID valueobject.PostID, params repository.ListParams, result *repository.CommentListResult) error {
	if m.setListByPostIDFn != nil {
		return m.setListByPostIDFn(ctx, postID, params, result)
	}
	return nil
}

func (m *mockCommentCacheRepository) DeleteListByPostID(ctx context.Context, postID valueobject.PostID) error {
	return nil
}

func (m *mockCommentCacheRepository) GetCountByPostID(ctx context.Context, postID valueobject.PostID) (int64, error) {
	if m.getCountByPostIDFn != nil {
		return m.getCountByPostIDFn(ctx, postID)
	}
	return -1, nil // Return -1 to indicate cache miss
}

func (m *mockCommentCacheRepository) SetCountByPostID(ctx context.Context, postID valueobject.PostID, count int64) error {
	if m.setCountByPostIDFn != nil {
		return m.setCountByPostIDFn(ctx, postID, count)
	}
	m.countCache[postID.String()] = count
	return nil
}

func (m *mockCommentCacheRepository) InvalidateComment(ctx context.Context, id valueobject.CommentID) error {
	if m.invalidateFn != nil {
		return m.invalidateFn(ctx, id)
	}
	delete(m.cache, id.String())
	return nil
}

func (m *mockCommentCacheRepository) InvalidateByPostID(ctx context.Context, postID valueobject.PostID) error {
	if m.invalidateByPostIDFn != nil {
		return m.invalidateByPostIDFn(ctx, postID)
	}
	return nil
}
