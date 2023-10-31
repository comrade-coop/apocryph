package ipdr

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/types"
	iface "github.com/ipfs/boxo/coreiface"
	"github.com/ipfs/boxo/coreiface/path"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/go-cid"
)

type IpdrTransport interface {
	types.ImageTransport

	NewDestinationReference(tag string) IpdrImageReference
	NewReference(p path.Path, tag string) IpdrImageReference
}

type ipdrTransport struct {
	ipfs iface.CoreAPI
}

func NewIpdrTransport(ipfs iface.CoreAPI) IpdrTransport {
	return &ipdrTransport{ipfs: ipfs}
}

func (*ipdrTransport) Name() string {
	return "ipdr"
}

func (t *ipdrTransport) ParseReference(reference string) (types.ImageReference, error) {
	base, tag, _ := strings.Cut(reference, ":")
	if tag == "" {
		tag = "latest"
	}
	path := path.New(base)
	if path.IsValid() != nil && base != "" {
		return nil, path.IsValid()
	}
	return &ipdrImageReference{transport: t, path: path, tag: tag}, nil
}

func (t *ipdrTransport) NewDestinationReference(tag string) IpdrImageReference {
	if tag == "" {
		tag = "latest"
	}
	return &ipdrImageReference{transport: t, path: path.New(""), tag: tag}
}

func (t *ipdrTransport) NewReference(p path.Path, tag string) IpdrImageReference {
	if tag == "" {
		tag = "latest"
	}
	return &ipdrImageReference{transport: t, path: p, tag: tag}
}

func (*ipdrTransport) ValidatePolicyConfigurationScope(scope string) error {
	return nil
}

type IpdrImageReference interface {
	types.ImageReference

	Path() path.Path
	Tag() string
}

type ipdrImageReference struct {
	transport *ipdrTransport
	path      path.Path
	tag       string
}

func (r *ipdrImageReference) Transport() types.ImageTransport {
	return r.transport
}

func (r *ipdrImageReference) Path() path.Path {
	return r.path
}

func (r *ipdrImageReference) Tag() string {
	return r.tag
}

func (r *ipdrImageReference) StringWithinTransport() string {
	return fmt.Sprintf("%s:%s", r.path.String(), r.tag)
}

func (*ipdrImageReference) DockerReference() reference.Named {
	return nil
}

func (r *ipdrImageReference) PolicyConfigurationIdentity() string {
	return "" // TODO?
}

func (r *ipdrImageReference) PolicyConfigurationNamespaces() []string {
	return []string{}
}

func (r *ipdrImageReference) NewImage(ctx context.Context, sys *types.SystemContext) (types.ImageCloser, error) {
	src, err := r.NewImageSource(ctx, sys)
	if err != nil {
		return nil, err
	}
	return image.FromSource(ctx, sys, src)
}

func (r *ipdrImageReference) NewImageSource(ctx context.Context, sys *types.SystemContext) (types.ImageSource, error) {
	return &ipdrImageSource{
		reference: r,
		ipfs:      r.transport.ipfs,
	}, nil // TODO
}

func (r *ipdrImageReference) NewImageDestination(ctx context.Context, sys *types.SystemContext) (types.ImageDestination, error) {
	return &ipdrImageDestination{
		reference:    r,
		ipfs:         r.transport.ipfs,
		blobs:        make(map[string]files.Node),
		manifests:    make(map[string]files.Node),
		cidsToRemove: []cid.Cid{},
	}, nil
}

func (*ipdrImageReference) DeleteImage(ctx context.Context, sys *types.SystemContext) error {
	return errors.ErrUnsupported
}
