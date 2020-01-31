package album

import (
	"context"
	"database/sql"
	"errors"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

var errCRUD = errors.New("error crud")

func TestCreateAlbumRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     CreateAlbumRequest
		wantError bool
	}{
		{"success", CreateAlbumRequest{Name: "test"}, false},
		{"required", CreateAlbumRequest{Name: ""}, true},
		{"too long", CreateAlbumRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateAlbumRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     UpdateAlbumRequest
		wantError bool
	}{
		{"success", UpdateAlbumRequest{Name: "test"}, false},
		{"required", UpdateAlbumRequest{Name: ""}, true},
		{"too long", UpdateAlbumRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService(&mockRepository{}, logger)

	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, 0, count)

	// successful creation
	album, err := s.Create(ctx, CreateAlbumRequest{Name: "test"})
	assert.Nil(t, err)
	assert.NotEmpty(t, album.ID)
	id := album.ID
	assert.Equal(t, "test", album.Name)
	assert.NotEmpty(t, album.CreatedAt)
	assert.NotEmpty(t, album.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation error in creation
	_, err = s.Create(ctx, CreateAlbumRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// unexpected error in creation
	_, err = s.Create(ctx, CreateAlbumRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	_, _ = s.Create(ctx, CreateAlbumRequest{Name: "test2"})

	// update
	album, err = s.Update(ctx, id, UpdateAlbumRequest{Name: "test updated"})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", album.Name)
	_, err = s.Update(ctx, "none", UpdateAlbumRequest{Name: "test updated"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, id, UpdateAlbumRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// unexpected error in update
	_, err = s.Update(ctx, id, UpdateAlbumRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	album, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", album.Name)
	assert.Equal(t, id, album.ID)

	// query
	albums, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, len(albums))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	album, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, album.ID)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)
}

type mockRepository struct {
	items []entity.Album
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Album, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.Album{}, sql.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int) ([]entity.Album, error) {
	return m.items, nil
}

func (m *mockRepository) Create(ctx context.Context, album entity.Album) error {
	if album.Name == "error" {
		return errCRUD
	}
	m.items = append(m.items, album)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, album entity.Album) error {
	if album.Name == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.ID == album.ID {
			m.items[i] = album
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
