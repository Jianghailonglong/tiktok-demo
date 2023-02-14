package service

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/huandu/go-clone"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"tiktok-demo/common"
	"tiktok-demo/conf"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/dao/redis"
	"tiktok-demo/logger"
	"tiktok-demo/middleware/kafka"
	"tiktok-demo/middleware/snowflake"
	"tiktok-demo/util"
	"time"
)

var (
	client *minio.Client
	err    error
)

func InitMinio() error {
	client, err = minio.New(conf.Config.MinioConfig.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.Config.AccessKeyID, conf.Config.SecretAccessKey, ""),
		Secure: conf.Config.UseSsL})
	if err != nil {
		logger.Log.Error("", zap.Any("error", err.Error()))
		return err
	}
	logger.Log.Info("init minio success")
	return nil
}

// GetFeed 获取视频流
func GetFeed(c *gin.Context, latestTime string, userId int) (feedResponse *common.FeedResponse, err error) {
	// 1、获取视频流
	videoList, err := mysql.GetFeed(latestTime)
	if err != nil {
		return
	}
	size := len(videoList)
	var wg sync.WaitGroup
	wg.Add(size * 4)
	// 2、获取用户信息
	resUsers := make([]common.User, size)

	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			if user, err := GetCommonUserInfoById(int64(userId), int64(videoList[i].AuthorId)); nil == err {
				resUsers[i] = user
			} else {
				logger.Log.Error("获取用户信息失败")
			}
		}(i)
	}
	// 3、获取视频点赞数
	videoLikeCntsList := make([]int, size)
	// 1）从redis获取
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			if cnt, err := redis.GetVideoFavoriteCount(c, videoList[i].Id); nil == err && cnt != 0 {
				videoLikeCntsList[i] = cnt
			} else {
				logger.Log.Error("从redis获取视频点赞次数失败")
				// 2）从mysql获取
				if cnt, err := mysql.GetFavoriteCount(videoList[i].Id); nil == err {
					videoLikeCntsList[i] = int(cnt)
					// 3）传给mysql更新点赞数
					err = kafka.FavoriteClient.SendMessage(strconv.Itoa(videoList[i].Id), strconv.Itoa(int(cnt)))
					if err != nil {
						logger.Log.Error("FavoriteClient.SendMessage failed", zap.Any("error", err))
					}
				} else {
					logger.Log.Error("从mysql获取点赞次数失败")
				}
			}
		}(i)
	}
	// 4、获取视频评论数
	videoCommentCntsList := make([]int, size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			if cnt, err := mysql.GetCommentCount(videoList[i].Id); nil == err {
				videoCommentCntsList[i] = int(cnt)
			} else {
				logger.Log.Error("获取评论次数失败")
			}
		}(i)
	}
	// 5、用户是否点赞视频
	isFavoriteVideoList := make([]bool, size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			// 1）从redis查看关系，没有关系还需要去mysql进一步确认
			if isFavorite, err := redis.CheckIsFavorite(c, userId, videoList[i].Id); err == nil && isFavorite {
				isFavoriteVideoList[i] = isFavorite
			} else {
				logger.Log.Error("从redis获取用户是否点赞失败")
				// 2）从mysql获取
				if isFavorite := mysql.CheckFavorite(userId, videoList[i].Id); nil == err {
					isFavoriteVideoList[i] = isFavorite
					// 3）传给mysql更新点赞关系
					if isFavorite {
						err = kafka.FavoriteClient.SendMessage(strconv.Itoa(videoList[i].Id), strconv.Itoa(userId)+":add")
					} else {
						err = kafka.FavoriteClient.SendMessage(strconv.Itoa(videoList[i].Id), strconv.Itoa(userId)+":del")
					}
					if err != nil {
						logger.Log.Error("FavoriteClient.SendMessage failed", zap.Any("error", err))
					}
				} else {
					logger.Log.Error("从mysql获取用户是否点赞失败")
				}
			}
		}(i)
	}
	wg.Wait()
	feedResponse = assembleFeed(videoList, resUsers, videoLikeCntsList, videoCommentCntsList, isFavoriteVideoList)
	return
}

// assembleFeed 组装视频流信息
func assembleFeed(videos []mysql.Video, resUsers []common.User, videoLikeCntsList []int,
	videoCommentCntsList []int, isFavoriteVideoList []bool) (feedResponse *common.FeedResponse) {
	feedResponse = new(common.FeedResponse)
	feedResponse.VideoList = make([]common.Video, len(videos))
	for i := 0; i < len(videos); i++ {
		feedResponse.VideoList[i].Id = int64(videos[i].Id)
		feedResponse.VideoList[i].PlayUrl = videos[i].PlayUrl
		feedResponse.VideoList[i].CoverUrl = videos[i].CoverUrl
		feedResponse.VideoList[i].Title = videos[i].Title
		feedResponse.VideoList[i].Author = resUsers[i]
		feedResponse.VideoList[i].FavoriteCount = int64(videoLikeCntsList[i])
		feedResponse.VideoList[i].CommentCount = int64(videoCommentCntsList[i])
		feedResponse.VideoList[i].IsFavorite = isFavoriteVideoList[i]
	}
	// 按照不同策略进行排序，1、默认发布时间倒序 2、点赞数排序 3、评论数排序
	// 下面仅仅为了演示使用
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(3)
	sortVideoCxt := util.SortVideoContext{}
	if r == 1 {
		sortVideoCxt.SetSortVideoStrategy(new(util.SortVideoByFavoriteCount))
		sortVideoCxt.SortVideo(feedResponse.VideoList)
	} else {
		sortVideoCxt.SetSortVideoStrategy(new(util.SortVideoByCommentCount))
		sortVideoCxt.SortVideo(feedResponse.VideoList)
	}
	// 为了让最新发布视频也可以展示，更新为当前时间
	if len(videos) > 0 {
		feedResponse.NextTime = time.Now().Unix()
	}
	return
}

// PublishVideo 发布视频
func PublishVideo(userID int, title string, inputBuffer *bytes.Buffer) error {
	var buffer1 *bytes.Buffer
	var buffer2 *bytes.Buffer
	// 深拷贝
	buffer1 = clone.Clone(inputBuffer).(*bytes.Buffer)
	buffer2 = clone.Clone(inputBuffer).(*bytes.Buffer)

	wait := sync.WaitGroup{}
	wait.Add(2)
	var videoError error
	var imageError error
	var imageName string
	var videoName string
	ctx := context.Background()
	go func(name *string) {
		*name, videoError = processVideo(ctx, buffer1)
		if videoError != nil {
			logger.Log.Error(videoError.Error())
		}
		wait.Done()
	}(&videoName)
	go func(name *string) {
		*name, imageError = processImage(ctx, buffer2)
		if imageError != nil {
			logger.Log.Error(imageError.Error())
		}
		wait.Done()
	}(&imageName)
	wait.Wait()
	if videoError != nil {
		return videoError
	}
	if imageError != nil {
		return imageError
	}
	// 访问路径需要存储到mysql中
	if videoName != "" && imageName != "" {
		err = mysql.InsertVideo(userID, title, videoName, imageName)
		if err != nil {
			return err
		}
	}
	return err
}

func processVideo(c context.Context, b *bytes.Buffer) (videoName string, err error) {
	bucketName := conf.Config.Video.BucketName
	videoId, err := snowflake.GenID()
	if err != nil {
		logger.Log.Error(err.Error())
		return "", err
	}
	videoName = strconv.Itoa(int(videoId)) + ".mp4" // 视频拼凑成.mp4格式
	contentType := conf.Config.Video.ContentType

	object, err := client.PutObject(c, bucketName, videoName, b,
		-1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		logger.Log.Error("上传失败", zap.Any("error", err.Error()))
		return "", err
	}
	logger.Log.Info("Successfully uploaded video", zap.Any("videoName", videoName), zap.Any("videoSize", object.Size))
	return videoName, nil
}

func processImage(c context.Context, b *bytes.Buffer) (imageName string, err error) {
	imageBytes := bytes.NewBuffer(nil)
	// 1、创建临时视频文件，作为生成图片的依据，最新版minio存储的都是文件夹类型，无法作为源文件依据
	file, err := os.Create("./test.mp4")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	_, err = io.Copy(file, b)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	// 2、生成封面图
	imageBytes = genImageBytes()
	if imageBytes == nil {
		logger.Log.Error("生成封面图失败")
		return "", err
	}

	// 3、上传封面图
	bucketName := conf.Config.Image.BucketName
	imageId, err := snowflake.GenID()
	if err != nil {
		logger.Log.Error(err.Error())
		return "", err
	}
	imageName = strconv.Itoa(int(imageId)) + ".jpg" // 视频拼凑成.mp4格式
	contentType := conf.Config.Image.ContentType
	object, err := client.PutObject(c, bucketName, imageName, imageBytes,
		-1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		logger.Log.Error("上传失败", zap.Any("error", err.Error()))
		return "", err
	}
	err = os.Remove("./test.mp4")
	if err != nil {
		logger.Log.Error("删除临时视频文件失败", zap.Any("error", err.Error()))
		return "", err
	}
	err = os.Remove("./test.jpg")
	if err != nil {
		logger.Log.Error("删除临时图片文件失败", zap.Any("error", err.Error()))
		return "", err
	}
	logger.Log.Info("Successfully uploaded image", zap.Any("imageName", imageName), zap.Any("imageSize", object.Size))
	return imageName, nil
}

func genImageBytes() *bytes.Buffer {
	var b bytes.Buffer
	// /usr/bin/ffmpeg -i ./test.mp4 -ss 00:00:01 -f image2 -vcodec mjpeg -vframes 1 -y ./test.jpg
	cmdArguments := []string{"-i", "./test.mp4", "-ss", "00:00:01", "-f", "image2", "-vcodec", "mjpeg", "-vframes", "1", "-y", "./test.jpg"}
	cmd := exec.Command("ffmpeg", cmdArguments...)
	err := cmd.Run()
	if err != nil {
		return nil
	}
	logger.Log.Info(cmd.String())
	file, err := os.Open("./test.jpg")
	if err != nil {
		logger.Log.Error("打开文件失败")
		return nil
	}
	_, err = io.Copy(&b, file)
	if err != nil {
		logger.Log.Error("io.Copy failed")
		return nil
	}
	return &b
}

// PublishList 通过userId获取发布的全部视频
func PublishList(c *gin.Context, userId int64) (videoPublishListResponse *common.VideoPublishListResponse, err error) {
	// 1、获取userID发布的视频
	videos, err := mysql.GetVideoListByUserID(userId)
	if err != nil {
		return
	}
	// 2、获取视频的作者，唯一作者
	author, err := mysql.GetInfoById(userId, userId)
	if err != nil {
		return
	}
	// 3、获取视频点赞数
	size := len(videos)
	var wg sync.WaitGroup
	wg.Add(size * 3)

	// 3、获取视频点赞数
	videoLikeCntsList := make([]int, size)
	// 1）从redis获取
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			if cnt, err := redis.GetVideoFavoriteCount(c, videos[i].Id); nil == err && cnt != 0 {
				videoLikeCntsList[i] = cnt
			} else {
				logger.Log.Error("从redis获取视频点赞次数失败")
				// 2）从mysql获取
				if cnt, err := mysql.GetFavoriteCount(videos[i].Id); nil == err {
					videoLikeCntsList[i] = int(cnt)
					// 3）传给mysql更新点赞数
					err = kafka.FavoriteClient.SendMessage(strconv.Itoa(videos[i].Id), strconv.Itoa(int(cnt)))
					if err != nil {
						logger.Log.Error("FavoriteClient.SendMessage failed", zap.Any("error", err))
					}
				} else {
					logger.Log.Error("从mysql获取点赞次数失败")
				}
			}
		}(i)
	}
	// 4、获取视频评论数
	videoCommentCntsList := make([]int, size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			if cnt, err := mysql.GetCommentCount(videos[i].Id); nil == err {
				videoCommentCntsList[i] = int(cnt)
			} else {
				logger.Log.Error("获取评论次数失败")
			}
		}(i)
	}
	// 5、用户是否点赞视频
	isFavoriteVideoList := make([]bool, size)
	for i := 0; i < size; i++ {
		go func(i int) {
			defer wg.Done()
			// 1）从redis查看关系
			if isFavorite, err := redis.CheckIsFavorite(c, int(userId), videos[i].Id); err == nil && isFavorite {
				isFavoriteVideoList[i] = isFavorite
			} else {
				logger.Log.Error("从redis获取用户是否点赞失败")
				// 2）从mysql获取
				if isFavorite := mysql.CheckFavorite(int(userId), videos[i].Id); nil == err {
					isFavoriteVideoList[i] = isFavorite
					// 3）传给mysql更新点赞关系
					if isFavorite {
						err = kafka.FavoriteClient.SendMessage(strconv.Itoa(videos[i].Id), strconv.Itoa(int(userId))+":add")
					} else {
						err = kafka.FavoriteClient.SendMessage(strconv.Itoa(videos[i].Id), strconv.Itoa(int(userId))+":del")
					}
					if err != nil {
						logger.Log.Error("FavoriteClient.SendMessage failed", zap.Any("error", err))
					}
				} else {
					logger.Log.Error("从mysql获取用户是否点赞失败")
				}
			}
		}(i)
	}
	wg.Wait()
	videoPublishListResponse = assemblePublishList(videos, author, videoLikeCntsList, videoCommentCntsList, isFavoriteVideoList)
	return
}

func assemblePublishList(videos []mysql.Video, author common.User, videoLikeCntsList []int,
	videoCommentCntsList []int, isFavoriteVideoList []bool) (videoPublishListResponse *common.VideoPublishListResponse) {
	videoPublishListResponse = new(common.VideoPublishListResponse)
	videoPublishListResponse.VideoList = make([]common.Video, len(videos))
	for i := 0; i < len(videos); i++ {
		videoPublishListResponse.VideoList[i].Id = int64(videos[i].Id)
		videoPublishListResponse.VideoList[i].PlayUrl = videos[i].PlayUrl
		videoPublishListResponse.VideoList[i].CoverUrl = videos[i].CoverUrl
		videoPublishListResponse.VideoList[i].Title = videos[i].Title
		videoPublishListResponse.VideoList[i].Author = author
		videoPublishListResponse.VideoList[i].FavoriteCount = int64(videoLikeCntsList[i])
		videoPublishListResponse.VideoList[i].CommentCount = int64(videoCommentCntsList[i])
		videoPublishListResponse.VideoList[i].IsFavorite = isFavoriteVideoList[i]
	}
	return
}
