package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUrlParts(t *testing.T) {
	test_cases := []struct {
		url   string
		parts []string
	}{
		{"", []string{}},
		{"/auth/v1/user", []string{"auth", "v1", "user"}},
		{"//auth//v1//user//", []string{"auth", "v1", "user"}},
		{"/inventory/category/:cat_id/products/:prod_id", []string{"inventory", "category", ":cat_id", "products", ":prod_id"}},
	}

	for _, testcase := range test_cases {
		assert.Equal(t, testcase.parts, getUrlParts(testcase.url))
	}
}

func TestFindTrailingNode(t *testing.T) {
	router := NewRouter()
	router.insertUrl("/auth/v1/user")
	node, params := router.findTrailerNode("/auth/v1/user")
	assert.NotNil(t, node)
	assert.Empty(t, params)
	node, params = router.findTrailerNode("/auth/v2/user")
	assert.Nil(t, node)
	assert.Empty(t, params)

	router.insertUrl("/inventory/category/:cat_id/products/:prod_id")
	node, params = router.findTrailerNode("/inventory/category/20/products/10011")
	assert.NotNil(t, node)
	assert.NotEmpty(t, params)
	assert.Equal(t, 2, len(params))
	assert.Equal(t, "20", params["cat_id"])
	assert.Equal(t, "20", params["cat_id"])

	router.insertUrl("/inventory/category/1/products/1")
	node, params = router.findTrailerNode("/inventory/category/1/products/1")
	assert.NotNil(t, node)
	assert.Empty(t, params)	
	
	assert.Panics(t, func() {
		router.insertUrl("/inventory/category/:cat_id_new/products/:prod_id")
	})
}
