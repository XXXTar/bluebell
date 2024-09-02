package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

//基于用户投票相关算法:http://www.ruanyifeng.com/blog/algorithm/

//本项目使用简化版的投票分数
//投一票就加432分 86400/200 ->需要200张赞成票可以给你帖子续一天 ——>《redis实战》

//投票的几种情况；
//direction=1时1.之前没投过票，现在投赞成票2.之前投反对，现在改投赞成票
//direction=0时1.之前投赞成票，现在取消，2.之前投反对票，现在取消
//direction=-1时1.之前没投过票，现在投反对票2.之前投赞成，现在改投反对票

//投票的限制:每个帖子自发表之日起，一个星期之内允许用户投票
//1.到期之后将redis中保存的赞成票数及反对票存储到mysql表中
//2.到期之后删除那个KeyPostVotedZSetPF

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票多少分
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID, communityID int64) error {
	pipeline := client.TxPipeline()
	//帖子时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	//帖子分数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})
	//更新把帖子id加到社区的set
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(communityID)))
	pipeline.SAdd(cKey, postID)
	_, err := pipeline.Exec()
	return err
}

func VoteForPost(userID, postID string, value float64) error {
	//1.判断投票限制
	//去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}
	//2和3需要放到一个pipeline事务中操作

	//2.更新分数
	//先查当前用户当前帖子的投票纪录
	ov := client.ZScore(getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()
	//如果这一次投票的值和之前的值一样，就提示不允许投票
	if value == ov {
		return ErrVoteRepeated
	}
	var op float64
	if value > ov {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(ov - value) //计算两次投票的差值
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), op*diff*scorePerVote, postID)

	//3.记录用户为该帖子投过票
	if value == 0 {
		pipeline.ZRem(getRedisKey(KeyPostVotedZSetPF+postID), userID).Result()
	} else {
		pipeline.ZAdd(getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  value, //赞成还是反对
			Member: userID,
		})
	}
	_, err := pipeline.Exec()
	return err
}
