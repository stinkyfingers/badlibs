package libscontroller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	libs "github.com/stinkyfingers/badlibs/models"

)

type mockStorage struct {
	lib *libs.Lib
	libs []libs.Lib
	err error
}

func (m *mockStorage) Get(id string) (*libs.Lib, error) {
	return m.lib, m.err
}

func (m *mockStorage)All() ([]libs.Lib, error){
	return m.libs, m.err
}
func (m *mockStorage)Delete(id string) error {
	return m.err
}
func (m *mockStorage)Update(l *libs.Lib) (*libs.Lib, error) {
	return m.lib, m.err

}
func (m *mockStorage)Create(l *libs.Lib) (*libs.Lib, error) {
	return m.lib, m.err
}

func TestAllLibs(t *testing.T) {
	expected := []libs.Lib{{Title: "test"}}
	s := Server{
		Storage: &mockStorage{
			libs: expected,
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/lib/all", nil)
	rec := httptest.NewRecorder()
	s.AllLibs(rec, req)
	res := rec.Result()
	defer res.Body.Close()
	require.Equal(t, http.StatusOK,  res.StatusCode)
	var ls []libs.Lib
	err := json.NewDecoder(res.Body).Decode(&ls)
	require.Nil(t ,err)
	require.Equal(t, expected, ls)
}