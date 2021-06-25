/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the vidx2pidx project. */

package main

import (
	"encoding/xml"
	"errors"
	"fmt"
)

type PdscXML struct {
	XMLName xml.Name `xml:"package"`
	Vendor  string   `xml:"vendor"`
	URL     string   `xml:"url"`
	Name    string   `xml:"name"`

	ReleasesTag struct {
		XMLName  xml.Name     `xml:"releases"`
		Releases []ReleaseTag `xml:"release"`
	} `xml:"releases"`
}

type ReleaseTag struct {
	XMLName xml.Name `xml:"release"`
	Version string   `xml:"version,attr"`
	Date    string   `xml:"Date,attr"`
}

func (p *PdscXML) MatchTag(pdscTag PdscTag) error {
	diff := ""

	if p.Vendor != pdscTag.Vendor {
		diff += fmt.Sprintf(" Vendor('%s' != '%s')", p.Vendor, pdscTag.Vendor)
	}
	if p.URL != pdscTag.URL {
		diff += fmt.Sprintf(" URL('%s' != '%s')", p.URL, pdscTag.URL)
	}
	if p.Name != pdscTag.Name {
		diff += fmt.Sprintf(" Name('%s' != '%s')", p.Name, pdscTag.Name)
	}
	if p.LatestVersion() != pdscTag.Version {
		diff += fmt.Sprintf(" Version('%s' != '%s')", p.LatestVersion(), pdscTag.Version)
	}

	if len(diff) > 0 {
		message := fmt.Sprintf("Pdsc tag '%s%s' does not match the actual file:%s", pdscTag.URL, pdscTag.Name, diff)
		return errors.New(message)
	}

	return nil
}

func (p *PdscXML) LatestVersion() string {
	releases := p.ReleasesTag.Releases
	if len(releases) > 0 {
		return releases[0].Version
	}
	return ""
}

func (p *PdscXML) Tag() PdscTag {
	return PdscTag{
		Vendor:  p.Vendor,
		URL:     p.URL,
		Name:    p.Name,
		Version: p.LatestVersion(),
	}
}
