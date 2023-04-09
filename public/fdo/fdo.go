package fdo

import "github.com/ondrejbudai/osbuild-composer-public/public/blueprint"

type Options struct {
	ManufacturingServerURL string
	DiunPubKeyInsecure     string
	DiunPubKeyHash         string
	DiunPubKeyRootCerts    string
}

func FromBP(bpFDO blueprint.FDOCustomization) *Options {
	fdo := Options(bpFDO)
	return &fdo
}
