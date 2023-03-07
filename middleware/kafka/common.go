package kafka

const (
	VideoTopic      = "video"
	VideoGroupId    = "video-group"
	FavoriteTopic   = "favorite"
	FavoriteGroupId = "favorite-group"
	ChatTopic       = "chat"
	ChatGroupId     = "chat-group"
	RelationTopic   = "relation"
	RelationGroupId = "relation-group"
)

var (
	VideoClient         VideoProducer
	VideoServerGroup    VideoConsumerGroup
	FavoriteClient      FavoriteProducer
	FavoriteServerGroup FavoriteConsumerGroup
	ChatClient          ChatProducer
	ChatServerGroup     ChatConsumerGroup
	RelationClient      RelationProducer
	RelationServerGroup RelationConsumerGroup
)
