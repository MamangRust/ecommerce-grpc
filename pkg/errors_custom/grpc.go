package errors_custom

import (
	"ecommerce/internal/pb"
	"encoding/json"
)

func GrpcErrorToJson(err *pb.ErrorResponse) string {
	jsonData, _ := json.Marshal(err)
	return string(jsonData)
}
