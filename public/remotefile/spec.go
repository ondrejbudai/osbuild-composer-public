package remotefile

import "github.com/ondrejbudai/osbuild-composer-public/public/worker/clienterrors"

type Spec struct {
	URL             string
	Content         []byte
	ResolutionError *clienterrors.Error
}
