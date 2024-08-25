package auth

import "go.mongodb.org/mongo-driver/mongo"

type CredentialStorage interface {
	getPassword(source, method, path string) (string, error)
}
type MysqlRepository struct {
	client mongo.Client
}

func NewMysqlRepository() *MysqlRepository {
	return &MysqlRepository{client: mongo.Client{}}
}

func (mr *MysqlRepository) GetPassword(source, method, path string) (string, error) {

	return "", nil
}
