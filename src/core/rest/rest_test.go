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
	Nodes     []core.Node
	NextLseq  int64
	CalledMap map[string]int
}

func newMockStorage() *mockStorage {
	mockStorage := new(mockStorage)
	mockStorage.Nodes = []core.Node{
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

	mockStorage.NextLseq = 4

	mockStorage.CalledMap = make(map[string]int)

	return mockStorage
}

func (ms *mockStorage) GetSubTreeNodes(root, last int64) ([]core.Node, error) {
	ms.CalledMap["GetSubTreeNodes"] += 1

	var res []core.Node
	for _, child := range ms.Nodes {
		if child.Lseq > last {
			res = append(res, child)
		}
	}
	return res, nil
}

func (ms *mockStorage) CreateNode(ref int64, content []byte) (core.Node, error) {
	ms.CalledMap["CreateNode"] += 1
	newNode := core.Node{
		Lseq:   ms.NextLseq,
		Parent: ref,
		Value:  content,
	}

	ms.Nodes = append(ms.Nodes, newNode)
	ms.NextLseq += 1

	return newNode, nil
}

func TestGetNodes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ms := newMockStorage()
	restAPI := NewRest(ms)

	// Define the test cases
	testCases := []struct {
		name              string
		root              string
		last              string
		expectedStatus    int
		expectedBody      []byte
		expectedCalledMap map[string]int
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
			expectedCalledMap: map[string]int{"GetSubTreeNodes": 1},
		},
		{
			name:              "Invalid query parameters lseq",
			root:              "0",
			last:              "Invalid",
			expectedStatus:    http.StatusBadRequest,
			expectedBody:      []byte(`"Invalid query parameters"`),
			expectedCalledMap: map[string]int{},
		},
		{
			name:              "Invalid query parameters root",
			root:              "Invalid",
			last:              "1",
			expectedStatus:    http.StatusBadRequest,
			expectedBody:      []byte(`"Invalid query parameters"`),
			expectedCalledMap: map[string]int{},
		},
	}

	// Run the test cases
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ms.CalledMap = make(map[string]int)

			// Create a new request and recorder
			rec := httptest.NewRecorder()

			// Bind the request to a gin context and call the handler
			c, _ := gin.CreateTestContext(rec)
			c.Params = append(c.Params, gin.Param{Key: "root", Value: tt.root})
			c.Params = append(c.Params, gin.Param{Key: "last", Value: tt.last})
			restAPI.GetNodes(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			assert.Equal(t, tt.expectedBody, rec.Body.Bytes())

			assert.Equal(t, tt.expectedCalledMap, ms.CalledMap)
		})
	}
}

func TestAddNode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ms := newMockStorage()
	rest := NewRest(ms)
	router := gin.New()
	router.POST("/nodes", rest.AddNode)

	newNode := core.Node{
		Parent: 1,
		Value:  []byte(`"test content"`),
	}

	responseNode := core.Node{
		Parent: 1,
		Lseq:   4,
		Value:  []byte(`"test content"`),
	}

	tests := []struct {
		name              string
		reqBody           interface{}
		expectedStatus    int
		expectedBody      []byte
		expectedNodes     []core.Node
		expectedCalledMap map[string]int
	}{
		{
			name:              "invalid ref format",
			reqBody:           []byte(`{ "ref": -1, "content": "some"}`),
			expectedStatus:    http.StatusBadRequest,
			expectedBody:      []byte(`"Invalid request body or parameters"`),
			expectedNodes:     nil,
			expectedCalledMap: map[string]int{},
		},
		{
			name:           "add node success",
			reqBody:        newNode,
			expectedStatus: http.StatusCreated,
			expectedBody: func() []byte {
				body, _ := json.Marshal(responseNode)
				return body
			}(),
			expectedNodes:     append(ms.Nodes, responseNode),
			expectedCalledMap: map[string]int{"CreateNode": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.CalledMap = make(map[string]int)

			rec := httptest.NewRecorder()

			reqBodyBytes, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/nodes", bytes.NewReader(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Bind the request to a gin context and call the handler
			c, _ := gin.CreateTestContext(rec)
			c.Request = req
			rest.AddNode(c)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.Equal(t, tt.expectedBody, rec.Body.Bytes())

			if tt.expectedNodes != nil {
				assert.Equal(t, tt.expectedNodes, ms.Nodes)
			}

			assert.Equal(t, tt.expectedCalledMap, ms.CalledMap)
		})
	}
}
