package middleware

import (
	"encoding/json"
	"net/http"
	"nms-controller/logic"
)

func GetGroups(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//取得定時被覆值的全局變數
	data := logic.GroupDatas
	// send the response
	json.NewEncoder(w).Encode(data)
}
