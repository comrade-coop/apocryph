// SPDX-License-Identifier: GPL-3.0

package ipdr

import (
	"context"
	"errors"
	"io"

	"github.com/containers/image/v5/manifest"
	"github.com/containers/image/v5/types"
	iface "github.com/ipfs/kubo/core/coreiface"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/boxo/files"
	"github.com/opencontainers/go-digest"
)

type ipdrImageSource struct {
	reference *ipdrImageReference
	ipfs      iface.CoreAPI
}

func sizeOrMinusOne(n files.Node) int64 {
	size, err := n.Size()
	if err != nil {
		size = -1
	}
	return size
}

func (s *ipdrImageSource) Reference() types.ImageReference {
	return s.reference
}

func (s *ipdrImageSource) Close() error {
	return nil
}

func (s *ipdrImageSource) GetManifest(ctx context.Context, instanceDigest *digest.Digest) ([]byte, string, error) {
	var p path.Path
	var err error
	if instanceDigest == nil {
		p, err = path.Join(s.reference.path, "manifests", s.reference.tag)
	} else {
		p, err = path.Join(s.reference.path, "manifests", instanceDigest.String())
	}
	if err != nil {
		return nil, "", err
	}
	manifestNode, err := s.ipfs.Unixfs().Get(ctx, p)
	if err != nil {
		return nil, "", err
	}
	manifestFile := files.ToFile(manifestNode)
	if manifestFile == nil {
		return nil, "", errors.New("manifest.json expected file")
	}
	bytes, err := io.ReadAll(manifestFile)
	if err != nil {
		return nil, "", err
	}
	return bytes, manifest.GuessMIMEType(bytes), nil
}

func (s *ipdrImageSource) GetBlob(ctx context.Context, blobInfo types.BlobInfo, _ types.BlobInfoCache) (io.ReadCloser, int64, error) {
	p, err := path.Join(s.reference.path, "blobs", blobInfo.Digest.String())
	if err != nil {
		return nil, 0, err
	}
	blobNode, err := s.ipfs.Unixfs().Get(ctx, p)
	if err != nil {
		return nil, 0, err
	}
	blobFile := files.ToFile(blobNode)
	if blobFile == nil {
		return nil, -1, errors.New("blob expected file")
	}
	return blobFile, sizeOrMinusOne(blobFile), nil
}

func (s *ipdrImageSource) HasThreadSafeGetBlob() bool {
	return true
}

func (s *ipdrImageSource) GetSignatures(ctx context.Context, instanceDigest *digest.Digest) ([][]byte, error) {
	return nil, nil // TODO, not yet supported
}

func (s *ipdrImageSource) LayerInfosForCopy(ctx context.Context, instanceDigest *digest.Digest) ([]types.BlobInfo, error) {
	return nil, nil
}
