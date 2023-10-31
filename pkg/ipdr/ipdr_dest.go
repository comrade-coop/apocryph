package ipdr

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/containers/image/v5/types"
	iface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/boxo/coreiface/options"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
	"github.com/opencontainers/go-digest"
)

var addOptions = []options.UnixfsAddOption{options.Unixfs.CidVersion(1)}

type ipdrImageDestination struct {
	reference    *ipdrImageReference
	ipfs         iface.CoreAPI
	cidsToRemove []cid.Cid
	blobs        map[string]files.Node
	manifests    map[string]files.Node
}

func (d *ipdrImageDestination) Reference() types.ImageReference {
	return d.reference
}
func (d *ipdrImageDestination) Close() error {
	for _, c := range d.cidsToRemove {
		err := d.ipfs.Dag().Remove(context.Background(), c)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *ipdrImageDestination) SupportedManifestMIMETypes() []string {
	return nil
}
func (d *ipdrImageDestination) SupportsSignatures(ctx context.Context) error {
	return errors.ErrUnsupported // TODO
}
func (d *ipdrImageDestination) DesiredLayerCompression() types.LayerCompression {
	return types.PreserveOriginal
}
func (d *ipdrImageDestination) AcceptsForeignLayerURLs() bool {
	return false
}
func (d *ipdrImageDestination) MustMatchRuntimeOS() bool {
	return false
}
func (d *ipdrImageDestination) IgnoresEmbeddedDockerReference() bool {
	return true
}

func (d *ipdrImageDestination) PutBlob(ctx context.Context, stream io.Reader, inputInfo types.BlobInfo, cache types.BlobInfoCache, isConfig bool) (types.BlobInfo, error) {
	var digester digest.Digester
	if inputInfo.Digest == "" {
		digester = digest.Canonical.Digester()
		stream = io.TeeReader(stream, digester.Hash())
	}

	path, err := d.ipfs.Unixfs().Add(ctx, files.NewReaderFile(stream), addOptions...)
	if err != nil {
		return inputInfo, err
	}
	inputInfo.URLs = append(inputInfo.URLs, fmt.Sprintf("ipfs://%s", path.Cid().String()))
	inputInfo.URLs = append(inputInfo.URLs, fmt.Sprintf("https://ipfs.io/ipfs/%s", path.Cid().String()))

	d.cidsToRemove = append(d.cidsToRemove, path.Cid())

	node, err := d.ipfs.Unixfs().Get(ctx, path)
	if err != nil {
		return inputInfo, err
	}

	if digester != nil {
		inputInfo.Digest = digester.Digest()
	}

	d.blobs[inputInfo.Digest.String()] = node

	inputInfo.Size = sizeOrMinusOne(node)

	return inputInfo, nil
}

func (d *ipdrImageDestination) HasThreadSafePutBlob() bool {
	return true
}

func (d *ipdrImageDestination) TryReusingBlob(ctx context.Context, info types.BlobInfo, cache types.BlobInfoCache, canSubstitute bool) (bool, types.BlobInfo, error) {
	var node files.Node
	for _, v := range info.URLs {
		u, err := url.Parse(v)
		if err == nil && (u.Scheme == "ipfs" || u.Host == "ipfs.io") {
			p := path.New(u.Path)
			if p.IsValid() == nil {
				var err error
				node, err = d.ipfs.Unixfs().Get(ctx, p)
				if err != nil {
					return false, info, err
				}
				d.blobs[info.Digest.String()] = node
				return true, info, nil
			}
		}
	}
	return false, info, nil
}

func (d *ipdrImageDestination) PutManifest(ctx context.Context, manifest []byte, instanceDigest *digest.Digest) error {
	manifestFile := files.NewBytesFile(manifest)

	if instanceDigest != nil {
		d.manifests[instanceDigest.String()] = manifestFile
	} else {
		d.manifests[digest.Canonical.FromBytes(manifest).String()] = manifestFile
		manifestFile = files.NewBytesFile(manifest)
		d.manifests[d.reference.tag] = manifestFile
	}

	return nil
}

func (d *ipdrImageDestination) PutSignatures(ctx context.Context, signatures [][]byte, instanceDigest *digest.Digest) error {
	return errors.ErrUnsupported
}

func (d *ipdrImageDestination) Commit(ctx context.Context, unparsedToplevel types.UnparsedImage) error {
	rootDir := files.NewMapDirectory(map[string]files.Node{
		"manifests": files.NewMapDirectory(d.manifests),
		"blobs":     files.NewMapDirectory(d.blobs),
	})
	path, err := d.ipfs.Unixfs().Add(ctx, rootDir, addOptions...)
	if err != nil {
		return err
	}
	if d.reference.path.Mutable() {
		_, err := d.ipfs.Name().Publish(ctx, d.reference.path)
		if err != nil {
			return err
		}
	} else {
		d.reference.path = path
		// TODO: Pin?
	}
	d.cidsToRemove = []cid.Cid{}
	return nil
}
