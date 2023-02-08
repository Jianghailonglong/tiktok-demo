package service

import (
	"bytes"
	"context"
	"github.com/huandu/go-clone"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"tiktok-demo/common"
	"tiktok-demo/conf"
	"tiktok-demo/dao/mysql"
	"tiktok-demo/logger"
	"tiktok-demo/middleware/snowflake"
)

var (
	client *minio.Client
	err    error
)

func InitMinio() error {
	client, err = minio.New(conf.Config.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.Config.AccessKeyID, conf.Config.SecretAccessKey, ""),
		Secure: conf.Config.UseSsL})
	if err != nil {
		logger.Log.Error("", zap.Any("error", err.Error()))
		return err
	}
	return nil
}

// GetFeed 获取视频流
func GetFeed(latestTime string, userID int) (feedResponse *common.FeedResponse, err error) {
	// 1、获取视频流以及各个视频作者的信息
	videoAndAuthors, err := mysql.GetFeed(latestTime)
	if err != nil {
		return
	}
	// 2、TODO 作者的粉丝数，关注数，是否被关注
	// 3、TODO 视频的点赞数，评论数，是否被点赞
	feedResponse = assembleFeed(videoAndAuthors)
	return
}

// assembleFeed 组装视频流信息
func assembleFeed(videoAndAuthors []mysql.VideoAndAuthor) (feedResponse *common.FeedResponse) {
	feedResponse = new(common.FeedResponse)
	feedResponse.VideoList = make([]common.Video, len(videoAndAuthors))
	for i := 0; i < len(videoAndAuthors); i++ {
		feedResponse.VideoList[i].Id = int64(videoAndAuthors[i].Video.Id)
		feedResponse.VideoList[i].PlayUrl = videoAndAuthors[i].PlayUrl
		feedResponse.VideoList[i].CoverUrl = videoAndAuthors[i].CoverUrl
		feedResponse.VideoList[i].CommentCount = int64(videoAndAuthors[i].CommentCount)
		feedResponse.VideoList[i].FavoriteCount = int64(videoAndAuthors[i].FavoriteCount)
		feedResponse.VideoList[i].Title = videoAndAuthors[i].Title
		feedResponse.VideoList[i].Author = common.User{
			Id:   int64(videoAndAuthors[i].User.Id),
			Name: videoAndAuthors[i].Username,
			// TODO FollowCount FollowerCount IsFollow
		}
	}
	// 更新为最早发布视频的时间
	if len(videoAndAuthors) > 0 {
		feedResponse.NextTime = videoAndAuthors[len(videoAndAuthors)-1].PublishTime
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

// PublishList 通过userID获取发布的全部视频
func PublishList(userID int64) (videoPublishListResponse *common.VideoPublishListResponse, err error) {
	// 1、获取userID发布的视频
	videos, err := mysql.GetVideoListByUserID(userID)
	if err != nil {
		return
	}
	// 2、获取视频的作者，唯一作者
	videoAuthor, err := mysql.GetUserByUserID(userID)
	if err != nil {
		return
	}
	// 3、TODO 作者的粉丝数，关注数，是否被关注
	// 4、TODO 视频的点赞数，评论数，是否被点赞
	videoPublishListResponse = assemblePublishList(videos, videoAuthor)
	return
}

func assemblePublishList(videos []mysql.Video, videoAuthor mysql.User) (videoPublishListResponse *common.VideoPublishListResponse) {
	videoPublishListResponse = new(common.VideoPublishListResponse)
	videoPublishListResponse.VideoList = make([]common.Video, len(videos))
	for i := 0; i < len(videos); i++ {
		videoPublishListResponse.VideoList[i].Id = int64(videos[i].Id)
		videoPublishListResponse.VideoList[i].PlayUrl = videos[i].PlayUrl
		videoPublishListResponse.VideoList[i].CoverUrl = videos[i].CoverUrl
		videoPublishListResponse.VideoList[i].CommentCount = int64(videos[i].CommentCount)
		videoPublishListResponse.VideoList[i].FavoriteCount = int64(videos[i].FavoriteCount)
		videoPublishListResponse.VideoList[i].Title = videos[i].Title

		videoPublishListResponse.VideoList[i].Author = common.User{
			Id:   int64(videoAuthor.Id),
			Name: videoAuthor.Username,
			// TODO FollowCount FollowerCount IsFollow
		}
	}
	return
}
