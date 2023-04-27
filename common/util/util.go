package util

import (
	"errors"
	"github.com/gofrs/uuid"
	"log"
	"os"
)

const (
	NodePath        = "./data"
	NodeFile        = "./data/key"
	LatestBlockPath = "./data"
	LatestBlockFile = "./data/blockchain.json"
)

func DeleteFile(path string) error {
	return os.Remove(path)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetLocalNodeId() (string, error) {

	if ok, _ := PathExists(NodeFile); ok {
		//存在
		bs, err := os.ReadFile(NodeFile)
		if err != nil {
			return "", err
		}
		return string(bs), nil
	} else {
		//不存在
		err := os.MkdirAll(NodePath, os.ModePerm)
		if err != nil {
			log.Println(err.Error())
			return "", err
		}
		err = os.Chmod(NodePath, 0777)
		if err != nil {
			return "", err
		}
		f, err := os.Create(NodeFile)
		if err != nil {
			return "", err
		}
		defer f.Close()

		u, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		nodeId := u.String()
		_, err = f.WriteString(nodeId)
		if err != nil {
			return "", err
		}
		return nodeId, nil
	}

}

func WriteLatestBlock(content string) error {
	if ok, _ := PathExists(LatestBlockFile); !ok {
		//不存在
		err := os.MkdirAll(LatestBlockPath, os.ModePerm)
		if err != nil {
			//log.Println(err.Error())
			return err
		}
		err = os.Chmod(LatestBlockPath, 0777)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(LatestBlockFile, []byte(content), 0777)
}

func ReadLatestBlock() ([]byte, error) {
	if ok, _ := PathExists(LatestBlockFile); ok {
		bs, err := os.ReadFile(LatestBlockFile)
		if err != nil {
			return nil, err
		}
		return bs, nil
	} else {
		return nil, errors.New("not found blockchain.json")
	}
}
