// this package contains ipfs helper functions
package ipfs_utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/kubo/client/rpc"
)

// Adds a file (or a directory) to IPFS using the provided IPFS node connection and the local unix file path.
// It returns the IPFS Content Identifier (CID) of the added file upon success or an error if the addition fails.
//
// # Parameters:
//   - node: A pointer to an established IPFS node connection used for interacting with IPFS.
//   - path: The local unix file path to be added to IPFS.
//   - verbose: (Optional) A boolean flag that controls whether additional information, including error messages, should be printed to the console.
//
// # Returns:
//   - string: The IPFS Content Identifier (CID) of the added file if successful, on failure this value will be an empty string.
//   - error: An error if the file addition to IPFS fails. If successful, this value will be nil.
//
// # Example Usage:
//
//	// Add the file to IPFS with verbose output
//	cid, err := AddFile(ipfsNode, "/path/to/file", true)
func AddFile(node *rpc.HttpApi, path string, verbose ...bool) (string, error) {
	ctx := context.Background()
	stats, err := os.Stat(path)
	if err != nil {
		msg := fmt.Sprintf("Error getting file stats for %v", path)
		handleError(err, msg, verbose)
		return "", err
	}
	serialFile, err := files.NewSerialFile(path, false, stats)
	if err != nil {
		handleError(err, "Error creating SerialFile", verbose)
		return "", err
	}
	response, err := node.Unixfs().Add(ctx, serialFile)
	if err != nil {
		handleError(err, "Error adding file to IPFS", verbose)
		return "", err
	}
	return response.String(), nil
}

func handleError(err error, msg string, verbose []bool) {
	if len(verbose) > 0 && verbose[0] == true {
		fmt.Fprintf(os.Stderr, "%v: %v\n", msg, err)
	}
}

// Retrieves a file or directory from IPFS using the provided IPFS node connection and Content Identifier (CID).
//
// # Parameters:
//   - node: A pointer to an established IPFS node connection used for interacting with IPFS.
//   - cid: The Content Identifier (CID) of the IPFS content to retrieve.
//   - savePath: (Optional) The local file system path where the retrieved content should be saved. If not provided, no local save operation is performed.
//
// # Returns:
//   - files.Node: The retrieved IPFS content as a files.Node if successful. The type of files.Node may vary based on the content type (file, directory, ...etc).
//   - error: An error if the retrieval or save operation fails. If successful, this value will be nil.
//
// # Example Usage:
//
//	// Retrieve IPFS content with CID 'QmXyz123' and print errors if they occur
//	content, err := RetrieveFile(ipfsNode, "QmXyz123", "/path/on/local/file/system", true)
func RetreiveFile(node *rpc.HttpApi, cid string, savePath ...string) (files.Node, error) {
	ctx := context.Background()

	response, err := node.Unixfs().Get(ctx, path.New(cid))
	if err != nil {
		fmt.Printf("Error retreiving file from IPFS:%v\n", err)
		return nil, err
	}
	if len(savePath) > 0 {
		SaveFile(response, savePath[0], true)
	}
	return response, nil
}

// Saves an IPFS file or directory represented by the given files.Node to the specified local file system path.
//
// # Parameters:
//   - file: The IPFS file or directory, represented as a files.Node, that you want to save to the local file system.
//   - path: The local file system path where the IPFS content should be saved. If the path does not exist, it will be created.
//   - verbose: (Optional) A boolean flag that controls whether additional information, including error messages, should be printed to the console.
//
// # Example Usage:
//
//	// Save the IPFS content to '/path/on/local/file/system' with verbose output
//	SaveFile(file, "/path/to/filename", true)
//
// # Returns:
//
//   - error: error upon falure, nil upon success
//
// Note:
//   - If the specified path does not exist, it will be created to save the IPFS content. and if the file exists it does not ovveride it and returns an error
func SaveFile(file files.Node, path string, verbose ...bool) error {
	err := files.WriteTo(file, path)
	if err != nil {
		handleError(err, "Error saving file", verbose)
		return err
	}
	return nil
}

// ConnectToLocalNode establishes a connection to the local IPFS node.
//
// This function creates and returns a pointer to an IPFS node connection (rpc.HttpApi) that can be used for interacting with the local IPFS node.
//
// Returns:
//   - *rpc.HttpApi: A pointer to the established IPFS node connection if successful.
//   - error: An error if the connection to the local node cannot be established. If successful, this value will be nil.
func ConnectToLocalNode() (*rpc.HttpApi, error) {
	node, err := rpc.NewLocalApi()
	if err != nil {
		fmt.Printf("Error creating local node connection: %v\n", err)
		return nil, err
	}
	return node, nil
}

func CreateIpfsService(node *rpc.HttpApi, pName string, endpoint string) error {
	ctx := context.Background()
	request := node.Request("p2p/listen", pName, endpoint)
	_, err := request.Send(ctx)
	if err != nil {
		log.Fatalf("Could not send request: %v", err)
		return err
	}
	return nil
}

func ForwardConnection(node *rpc.HttpApi, pName string, endpoint, target string) error {
	target = fmt.Sprintf("/p2p/%v", target)
	ctx := context.Background()
	request := node.Request("p2p/forward", pName, endpoint, target)
	_, err := request.Send(ctx)
	if err != nil {
		fmt.Printf("Could not send request: %v", err)
		return err
	}
	return nil
}
