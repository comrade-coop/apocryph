// this package contains ipfs helper functions
package ipfs

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"google.golang.org/protobuf/proto"
)

// Adds a file (or a directory) from the local filesystem to IPFS
func AddFsFile(node *rpc.HttpApi, path string) (path.Resolved, error) {
	stats, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	serialFile, err := files.NewSerialFile(path, false, stats)
	if err != nil {
		return nil, err
	}
	return node.Unixfs().Add(context.Background(), serialFile)
}

// Saves an IPFS file or directory represented by the given files.Node to the specified local file system path.
func GetFsFile(node *rpc.HttpApi, file files.Node, path path.Path, savePath string) error {
	file, err := node.Unixfs().Get(context.Background(), path)
	if err != nil {
		return err
	}
	return files.WriteTo(file, savePath)
}

// Adds a file from a protobuf slice to IPFS
func AddProtobufFile(node *rpc.HttpApi, msg proto.Message) (cid.Cid, error) {
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return cid.Cid{}, err
	}
	path, err := node.Unixfs().Add(context.Background(), files.NewBytesFile(bytes))
	if err != nil {
		return cid.Cid{}, err
	}
	return path.Cid(), nil
}

// Adds a file from a protobuf slice to IPFS
func GetProtobufFile(node *rpc.HttpApi, cid cid.Cid, msg proto.Message) error {
	fileNode, err := node.Unixfs().Get(context.Background(), path.IpfsPath(cid))
	if err != nil {
		return err
	}
	file, ok := fileNode.(files.File)
	if !ok {
		return fmt.Errorf("Supplied CID (%s) not a file", cid.String())
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	return proto.Unmarshal(bytes, msg)
}
