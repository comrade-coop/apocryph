// SPDX-License-Identifier: GPL-3.0

package ipdr

import (
	"context"
	"fmt"
	"strings"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/types"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	iface "github.com/ipfs/kubo/core/coreiface"
)

// An implementation of [types.ImageTransport] for IPDR-stored images
type IpdrTransport interface {
	types.ImageTransport

	// Create an [IpdrImageReference] that will be used as a destination of a push or copy command.
	NewDestinationReference(tag string) IpdrImageReference
	// Create an [IpdrImageReference] from an IPFS path and tag. Alternatively, these can be created through ParseReference() using path.String() + ":" + tag, or through [github.com/containers/image/v5/transports/alltransports.ParseImageName], using "ipdr:" + path.String() + ":" + tag.
	NewReference(p path.Path, tag string) IpdrImageReference
}

// Implements an [IpdrTransport]
type ipdrTransport struct {
	ipfs iface.CoreAPI
}

// Create a new IPDR transport given a connection to IPFS.
func NewIpdrTransport(ipfs iface.CoreAPI) IpdrTransport {
	// TODO: it would be nice if we could create the IPFS connection on the fly using [types.SystemContext]; sadly, there aren't many way to extend the SystemContext struct with such data.
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
	path, err := path.NewPath(base)
	if err != nil && base != "" {
		return nil, err
	}
	return &ipdrImageReference{transport: t, path: path, tag: tag}, nil
}

func (t *ipdrTransport) NewDestinationReference(tag string) IpdrImageReference {
	if tag == "" {
		tag = "latest"
	}
	return &ipdrImageReference{transport: t, path: nil, tag: tag}
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

// A reference to an image stored in IPDR. Extends [types.ImageReference] with methods for getting the exact IPFS Path where the image is store as well as the tag it is being accessed through.
// If the path is mutable, writes to the image (through [types.ImageReference.NewImageDestination]) will result in an IPNS update of the mutable data; otherwise writes to the image will result in the [IpdrImageReference] changing its Path, and thus its serialized StringWithinTransport() / [github.com/containers/image/v5/transports.ImageName] form
type IpdrImageReference interface {
	types.ImageReference

	// Path returns the IPFS path that holds all of the image's contents.
	Path() path.Path
	// Tag returns the string tag that the image is accessed as. Typically, "latest" is used.
	Tag() string
}

// Implements an [IpdrImageReference]
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

func (r *ipdrImageReference) DeleteImage(ctx context.Context, sys *types.SystemContext) error {
	return r.transport.ipfs.Pin().Rm(ctx, r.path)
}
