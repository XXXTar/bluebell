package controller

import "Project/bluebell/models"

//专门用来放接口文档用到的model
//因为我们的接口文档返回的数据格式是一致的，但是具体的data类型不一样
type _ResponsePostList struct {
	Code    ResCode                 `json:"code"`
	Message string                  `json:"message"`
	Data    []*models.ApiPostDetail `json:"data"`
}
