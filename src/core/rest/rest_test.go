package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"helicopter/core"
)

type mockStorage struct {
	nodes []core.Node
}

func newMockStorage() *mockStorage {
	mockStorage := new(mockStorage)
	mockStorage.nodes = []core.Node{
		{
			Lseq:   0,
			Parent: -1,
			Value:  []byte("root"),
		},
		{
			Lseq:   1,
			Parent: 0,
			Value:  []byte("child1"),
		},
		{
			Lseq:   2,
			Parent: 0,
			Value:  []byte("child2"),
		},
		{
			Lseq:   3,
			Parent: 2,
			Value:  []byte("child3"),
		},
	}

	return mockStorage
}

func (ms *mockStorage) GetSubTreeNodes(root, last int64) ([]core.Node, error) {
	var res []core.Node
	for _, child := range ms.nodes {
		if child.Lseq > last {
			res = append(res, child)
		}
	}
	return res, nil
}

func (ms *mockStorage) GetNode(nodeID int64) (*core.Node, error) {
	return nil, nil
}

func (ms *mockStorage) CreateNode(ref int64, content []byte) (core.Node, error) {
	node := core.Node{
		Lseq:   1,
		Parent: ref,
		Value:  content,
	}
	return node, nil
}

func TestGetNodes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockStorage := newMockStorage()
	restAPI := NewRest(mockStorage)

	// Define the test cases
	testCases := []struct {
		name           string
		root           string
		last           string
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "Valid range",
			root:           "0",
			last:           "1",
			expectedStatus: http.StatusOK,
			expectedBody: func() []byte {
				nodes, _ := json.Marshal([]*core.Node{
					{
						Lseq:   2,
						Parent: 0,
						Value:  []byte("child2"),
					},
					{
						Lseq:   3,
						Parent: 2,
						Value:  []byte("child3"),
					},
				})
				return nodes
			}(),
		},
		{
			name:           "Invalid query parameters lseq",
			root:           "0",
			last:           "Invalid",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte(`"Invalid query parameters"`),
		},
		{
			name:           "Invalid query parameters root",
			root:           "Invalid",
			last:           "1",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte(`"Invalid query parameters"`),
		},
	}

	// Run the test cases
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request and recorder
			rec := httptest.NewRecorder()

			// Bind the request to a gin context and call the handler
			c, _ := gin.CreateTestContext(rec)
			c.Params = append(c.Params, gin.Param{Key: "root", Value: tt.root})
			c.Params = append(c.Params, gin.Param{Key: "last", Value: tt.last})
			restAPI.GetNodes(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			assert.Equal(t, tt.expectedBody, rec.Body.Bytes())
		})
	}
}

func TestAddNode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rest := NewRest(&mockStorage{})
	router := gin.New()
	router.POST("/nodes", rest.AddNode)

	tests := []struct {
		name     string
		reqBody  interface{}
		expected int
	}{
		{
			name:     "invalid request body",
			reqBody:  []byte(`{ "ref": 0 }`),
			expected: http.StatusBadRequest,
		},
		{
			name: "create node success",
			reqBody: core.Node{
				Parent: 1,
				Value:  []byte("test content"),
			},
			expected: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			reqBodyBytes, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/nodes", bytes.NewReader(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expected, w.Code)

			if tt.expected == http.StatusCreated {
				var responseNode core.Node
				err := json.NewDecoder(w.Body).Decode(&responseNode)
				assert.NoError(t, err)
				assert.Equal(t, int64(1), responseNode.Lseq)
			}
		})
	}
}
