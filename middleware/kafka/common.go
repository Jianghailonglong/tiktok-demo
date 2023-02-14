package kafka

const (
	FavoriteTopic   = "favorite"
	FavoriteGroupId = "favorite-group"
)

var (
	FavoriteClient      FavoriteProducer
	FavoriteServerGroup FavoriteConsumerGroup
)
