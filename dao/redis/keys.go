package redis

// redis key 注意使用命名空间的方式，方便查询和拆分
const (
	// KeyFavoriteUserIdPrefix key->tiktok:favorite:user:userid, value: videoid的set集合，用户点赞视频集合
	KeyFavoriteUserIdPrefix = "tiktok:favorite:user:"
	// KeyFavoriteVideoIdPrefix key->tiktok:favorite:video:videoid, value: count,视频点赞次数
	KeyFavoriteVideoIdPrefix = "tiktok:favorite:video:"
	// KeyPublishUserIdPrefix key->tiktok:publish:user:userid, value: videoid的set集合，用户发布视频集合
	KeyPublishUserIdPrefix = "tiktok:publish:user:"
	// KeyFollowUserIdPrefix key->tiktok:follow:user:userid, value: to_userid的set集合，用户关注列表集合
	KeyFollowUserIdPrefix = "tiktok:follow:user:"
	// KeyFollowerUserIdPrefix key->tiktok:follower:user:to_userid, value: userid的set集合，用户粉丝列表集合
	KeyFollowerUserIdPrefix = "tiktok:follower:user:"
)
