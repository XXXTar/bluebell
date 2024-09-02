package logic

import (
	"Project/bluebell/dao/mysql"
	"Project/bluebell/models"
)

func GetCommunityList() ([]*models.Community, error) {
	//查询数据库 查找到所有的community并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
