package bnb

import (
	"testing"

	"github.com/0xcregis/easynode/store"
)

func TestParseTx(t *testing.T) {

	str := `
{
  "id": 1694487255459894000,
  "hash": "0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21",
  "txTime": "",
  "txStatus": "",
  "blockNumber": "18117360",
  "from": "0x2c2ab61d2506308c0017f26c36e81e5b22942d57",
  "to": "0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768",
  "value": "0x0",
  "fee": "",
  "gasPrice": "0x2d3fec87b",
  "maxFeePerGas": "0x40f25e72a",
  "gas": "0x25bd4",
  "gasUsed": "",
  "baseFeePerGas": "",
  "maxPriorityFeePerGas": "0x5f5e100",
  "input": "0x4a21a2df0000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000021fb3f",
  "blockHash": "0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36",
  "transactionIndex": "0x1e",
  "type": "0x2",
  "receipt": "{\"id\":1694487256952455000,\"blockHash\":\"0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36\",\"logsBloom\":\"0x00000000000040010010000000000000000000000000000000000000000200000000000000000000000000010000000003000000000000000000000000000000080000000000000000000008000000000000000080000001000000000000000000000000020000000800000000000800000300001000000000000010000000000000000000010000000000000000000000002000000000000000010000000000004000000402020000000000000000000880000000000000000000000000200000000002000000000010000000000000000000002000003000080002004020000000000040000000000000000000000000100000000080000000000000000000\",\"contractAddress\":\"\",\"transactionIndex\":\"0x1e\",\"type\":\"0x2\",\"transactionHash\":\"0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21\",\"gasUsed\":\"0x24b86\",\"blockNumber\":\"18117360\",\"cumulativeGasUsed\":\"0x245b06\",\"from\":\"0x2c2ab61d2506308c0017f26c36e81e5b22942d57\",\"to\":\"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768\",\"effectiveGasPrice\":\"0x2d3fec87b\",\"logs\":[{\"blockHash\":\"0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36\",\"address\":\"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768\",\"logIndex\":\"0x38\",\"data\":\"{\\\"data\\\":\\\"0x\\\",\\\"eip\\\":721,\\\"token\\\":\\\"{\\\\\\\"total\\\\\\\":\\\\\\\"0x1200\\\\\\\"}\\\"}\",\"removed\":false,\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"0x0000000000000000000000002c2ab61d2506308c0017f26c36e81e5b22942d57\",\"0x000000000000000000000000000000000000000000000000000000000000051a\"],\"blockNumber\":\"0x11472f0\",\"transactionIndex\":\"0x1e\",\"transactionHash\":\"0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21\"},{\"blockHash\":\"0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36\",\"address\":\"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768\",\"logIndex\":\"0x39\",\"data\":\"{\\\"data\\\":\\\"0x\\\",\\\"eip\\\":721,\\\"token\\\":\\\"{\\\\\\\"total\\\\\\\":\\\\\\\"0x1200\\\\\\\"}\\\"}\",\"removed\":false,\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"0x0000000000000000000000002c2ab61d2506308c0017f26c36e81e5b22942d57\",\"0x000000000000000000000000000000000000000000000000000000000000051b\"],\"blockNumber\":\"0x11472f0\",\"transactionIndex\":\"0x1e\",\"transactionHash\":\"0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21\"},{\"blockHash\":\"0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36\",\"address\":\"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768\",\"logIndex\":\"0x3a\",\"data\":\"{\\\"data\\\":\\\"0x\\\",\\\"eip\\\":721,\\\"token\\\":\\\"{\\\\\\\"total\\\\\\\":\\\\\\\"0x1200\\\\\\\"}\\\"}\",\"removed\":false,\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"0x0000000000000000000000002c2ab61d2506308c0017f26c36e81e5b22942d57\",\"0x000000000000000000000000000000000000000000000000000000000000051c\"],\"blockNumber\":\"0x11472f0\",\"transactionIndex\":\"0x1e\",\"transactionHash\":\"0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21\"},{\"blockHash\":\"0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36\",\"address\":\"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768\",\"logIndex\":\"0x3b\",\"data\":\"{\\\"data\\\":\\\"0x\\\",\\\"eip\\\":721,\\\"token\\\":\\\"{\\\\\\\"total\\\\\\\":\\\\\\\"0x1200\\\\\\\"}\\\"}\",\"removed\":false,\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"0x0000000000000000000000002c2ab61d2506308c0017f26c36e81e5b22942d57\",\"0x000000000000000000000000000000000000000000000000000000000000051d\"],\"blockNumber\":\"0x11472f0\",\"transactionIndex\":\"0x1e\",\"transactionHash\":\"0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21\"},{\"blockHash\":\"0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36\",\"address\":\"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768\",\"logIndex\":\"0x3c\",\"data\":\"{\\\"data\\\":\\\"0x\\\",\\\"eip\\\":721,\\\"token\\\":\\\"{\\\\\\\"total\\\\\\\":\\\\\\\"0x1200\\\\\\\"}\\\"}\",\"removed\":false,\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"0x0000000000000000000000002c2ab61d2506308c0017f26c36e81e5b22942d57\",\"0x000000000000000000000000000000000000000000000000000000000000051e\"],\"blockNumber\":\"0x11472f0\",\"transactionIndex\":\"0x1e\",\"transactionHash\":\"0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21\"},{\"blockHash\":\"0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36\",\"address\":\"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768\",\"logIndex\":\"0x3d\",\"data\":\"{\\\"data\\\":\\\"0x\\\",\\\"eip\\\":721,\\\"token\\\":\\\"{\\\\\\\"total\\\\\\\":\\\\\\\"0x1200\\\\\\\"}\\\"}\",\"removed\":false,\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"0x0000000000000000000000002c2ab61d2506308c0017f26c36e81e5b22942d57\",\"0x000000000000000000000000000000000000000000000000000000000000051f\"],\"blockNumber\":\"0x11472f0\",\"transactionIndex\":\"0x1e\",\"transactionHash\":\"0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21\"},{\"blockHash\":\"0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36\",\"address\":\"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768\",\"logIndex\":\"0x3e\",\"data\":\"{\\\"data\\\":\\\"0x\\\",\\\"eip\\\":721,\\\"token\\\":\\\"{\\\\\\\"total\\\\\\\":\\\\\\\"0x1200\\\\\\\"}\\\"}\",\"removed\":false,\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"0x0000000000000000000000002c2ab61d2506308c0017f26c36e81e5b22942d57\",\"0x0000000000000000000000000000000000000000000000000000000000000520\"],\"blockNumber\":\"0x11472f0\",\"transactionIndex\":\"0x1e\",\"transactionHash\":\"0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21\"},{\"blockHash\":\"0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36\",\"address\":\"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768\",\"logIndex\":\"0x3f\",\"data\":\"{\\\"data\\\":\\\"0x\\\",\\\"eip\\\":721,\\\"token\\\":\\\"{\\\\\\\"total\\\\\\\":\\\\\\\"0x1200\\\\\\\"}\\\"}\",\"removed\":false,\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"0x0000000000000000000000002c2ab61d2506308c0017f26c36e81e5b22942d57\",\"0x0000000000000000000000000000000000000000000000000000000000000521\"],\"blockNumber\":\"0x11472f0\",\"transactionIndex\":\"0x1e\",\"transactionHash\":\"0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21\"},{\"blockHash\":\"0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36\",\"address\":\"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768\",\"logIndex\":\"0x40\",\"data\":\"{\\\"data\\\":\\\"0x\\\",\\\"eip\\\":721,\\\"token\\\":\\\"{\\\\\\\"total\\\\\\\":\\\\\\\"0x1200\\\\\\\"}\\\"}\",\"removed\":false,\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"0x0000000000000000000000002c2ab61d2506308c0017f26c36e81e5b22942d57\",\"0x0000000000000000000000000000000000000000000000000000000000000522\"],\"blockNumber\":\"0x11472f0\",\"transactionIndex\":\"0x1e\",\"transactionHash\":\"0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21\"},{\"blockHash\":\"0xbe36cdcfce377f7415bd91be3be10555fc705cd9c48ac077b3de9a1c298c4a36\",\"address\":\"0xd9ec62e6927082ad28b73fb5d4b5e9d571e00768\",\"logIndex\":\"0x41\",\"data\":\"{\\\"data\\\":\\\"0x\\\",\\\"eip\\\":721,\\\"token\\\":\\\"{\\\\\\\"total\\\\\\\":\\\\\\\"0x1200\\\\\\\"}\\\"}\",\"removed\":false,\"topics\":[\"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef\",\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"0x0000000000000000000000002c2ab61d2506308c0017f26c36e81e5b22942d57\",\"0x0000000000000000000000000000000000000000000000000000000000000523\"],\"blockNumber\":\"0x11472f0\",\"transactionIndex\":\"0x1e\",\"transactionHash\":\"0x2b7b684d469c365e0f8d9e2bf94bee672878aff4604b7715a48a7f37432f1a21\"}],\"createTime\":\"2023-09-12\",\"status\":\"0x1\"}"
}

`
	tx, err := ParseTx([]byte(str), store.EthTopic, store.EthTransferSingleTopic, 200)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(tx)
	}

}
