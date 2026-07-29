package storage

type FileMetadata struct {
	FileID    int64  `dynamodb:"FileID"`
	Filename  string `dynamodbav:"Filename"`
	Extension string `dynamodbav:"Extension"`
	Size      int64  `dynamodbav:"Size"`
	Type      string `dynamodbav:"Type"`
}
