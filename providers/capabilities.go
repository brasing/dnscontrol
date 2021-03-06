package providers

import (
	"log"
)

//Capability is a bitmasked set of "features" that a provider supports. Only use constants from this package.
type Capability uint32

const (
	// CanUseAlias indicates the provider support ALIAS records (or flattened CNAMES). Up to the provider to translate them to the appropriate record type.
	// If you add something to this list, you probably want to add it to pkg/normalize/validate.go checkProviderCapabilities() or somewhere near there.
	CanUseAlias Capability = iota
	// CanUsePTR indicates the provider can handle PTR records
	CanUsePTR
	// CanUseSRV indicates the provider can handle SRV records
	CanUseSRV
	// CanUseCAA indicates the provider can handle CAA records
	CanUseCAA
	// CantUseNOPURGE indicates NO_PURGE is broken for this provider. To make it
	// work would require complex emulation of an incremental update mechanism,
	// so it is easier to simply mark this feature as not working for this
	// provider.
	CantUseNOPURGE

	// DocOfficiallySupported means it is actively used and maintained by stack exchange
	DocOfficiallySupported
	// DocDualHost means provider allows full management of apex NS records, so we can safely dual-host with anothe provider
	DocDualHost
	// DocCreateDomains means provider can add domains with the `dnscontrol create-domains` command
	DocCreateDomains
)

var providerCapabilities = map[string]map[Capability]bool{}

func ProviderHasCabability(pType string, cap Capability) bool {
	if providerCapabilities[pType] == nil {
		return false
	}
	return providerCapabilities[pType][cap]
}

// DocumentationNote is a way for providers to give more detail about what features they support.
type DocumentationNote struct {
	HasFeature bool
	Comment    string
	Link       string
}

// DocumentationNotes is a full list of notes for a single provider
type DocumentationNotes map[Capability]*DocumentationNote

// ProviderMetadata is a common interface for DocumentationNotes and Capability to be used interchangably
type ProviderMetadata interface{}

// Notes is a collection of all documentation notes, keyed by provider type
var Notes = map[string]DocumentationNotes{}

func unwrapProviderCapabilities(pName string, meta []ProviderMetadata) {
	for _, pm := range meta {
		switch x := pm.(type) {
		case Capability:
			if providerCapabilities[pName] == nil {
				providerCapabilities[pName] = map[Capability]bool{}
			}
			providerCapabilities[pName][x] = true
		case DocumentationNotes:
			if Notes[pName] == nil {
				Notes[pName] = DocumentationNotes{}
			}
			for k, v := range x {
				Notes[pName][k] = v
			}
		default:
			log.Fatalf("Unrecognized ProviderMetadata type: %T", pm)
		}

	}
}

// Can is a small helper for concisely creating Documentation Notes
// comments are variadic for easy ommission. First is comment, second is link, the rest are ignored.
func Can(comments ...string) *DocumentationNote {
	n := &DocumentationNote{
		HasFeature: true,
	}
	n.addStrings(comments)
	return n
}

// Cannot is a small helper for concisely creating Documentation Notes
// comments are variadic for easy ommission. First is comment, second is link, the rest are ignored.
func Cannot(comments ...string) *DocumentationNote {
	n := &DocumentationNote{
		HasFeature: false,
	}
	n.addStrings(comments)
	return n
}

func (n *DocumentationNote) addStrings(comments []string) {
	if len(comments) > 0 {
		n.Comment = comments[0]
	}
	if len(comments) > 1 {
		n.Link = comments[1]
	}
}
