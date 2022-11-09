package main

//func TestGetPeerVersion(t *testing.T) {
//	router := setupRouter()
//
//	w := httptest.NewRecorder()
//
//	for i := 0; i < 1000; i++ {
//		req, _ := http.NewRequest("GET", "/fabric/peer/", nil)
//		router.ServeHTTP(w, req)
//
//		assert.Equal(t, 200, w.Code)
//	}
//}
//
//func TestESSearchAll(t *testing.T) {
//	router := setupRouter()
//
//	w := httptest.NewRecorder()
//
//	for i := 0; i < 1000; i++ {
//		req, _ := http.NewRequest("GET", "/fabric/dashboard/smart-contracts?filter={}", nil)
//		router.ServeHTTP(w, req)
//		assert.Equal(t, 200, w.Code)
//	}
//}

//func TestInstallWithDeployCC(t *testing.T) {
//	router := setupRouter()
//
//	w := httptest.NewRecorder()
//
//	var jsonData = []byte(`{
//		"cc_name": "basic",
//		"cc_path": "asset-transfer-basic/chaincode-go/",
//		"cc_language": "go"
//	}`)
//
//	req, _ := http.NewRequest("POST", "/fabric/dashboard/deployCC", bytes.NewBuffer(jsonData))
//	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
//	router.ServeHTTP(w, req)
//	assert.Equal(t, 200, w.Code)
//}
