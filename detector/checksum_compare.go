package detector

import (
	"talisman/gitrepo"
	"talisman/talismanrc"
	"talisman/utility"
)

type ChecksumCompare struct {
	additions    []gitrepo.Addition
	ignoreConfig *talismanrc.TalismanRC
}

//NewChecksumCompare returns new instance of the ChecksumCompare
func NewChecksumCompare(gitAdditions []gitrepo.Addition, talismanRCConfig *talismanrc.TalismanRC) *ChecksumCompare {
	cc := ChecksumCompare{additions: gitAdditions, ignoreConfig: talismanRCConfig}
	return &cc
}

func (cc *ChecksumCompare) IsScanNotRequired(addition gitrepo.Addition) bool {
	currentCollectiveChecksum := utility.CollectiveSHA256Hash([]string{string(addition.Path)})
	declaredCheckSum := ""
	for _, ignore := range cc.ignoreConfig.FileIgnoreConfig {
		if addition.Matches(ignore.FileName) {
			currentCollectiveChecksum = utility.CollectiveSHA256Hash([]string{ignore.FileName})
			declaredCheckSum = ignore.Checksum
		}

	}
	return currentCollectiveChecksum == declaredCheckSum

}

//FilterIgnoresBasedOnChecksums filters the file ignores from the talismanrc.TalismanRC which doesn't have any checksum value or having mismatched checksum value from the .talsimanrc
func (cc *ChecksumCompare) FilterIgnoresBasedOnChecksums() talismanrc.TalismanRC {
	finalIgnores := []talismanrc.FileIgnoreConfig{}
	for _, ignore := range cc.ignoreConfig.FileIgnoreConfig {
		currentCollectiveChecksum := cc.calculateCollectiveChecksumForPattern(ignore.FileName, cc.additions)
		// Compare with previous checksum from talismanrc.FileIgnoreConfig
		if ignore.Checksum == currentCollectiveChecksum {
			finalIgnores = append(finalIgnores, ignore)
		}
	}
	rc := talismanrc.TalismanRC{}
	rc.FileIgnoreConfig = finalIgnores
	return rc
}

func (cc *ChecksumCompare) calculateCollectiveChecksumForPattern(fileNamePattern string, additions []gitrepo.Addition) string {
	var patternpaths []string
	currentCollectiveChecksum := ""
	for _, addition := range additions {
		if addition.Matches(fileNamePattern) {
			patternpaths = append(patternpaths, string(addition.Path))
		}
	}
	// Calculate current collective checksum
	patternpaths = utility.UniqueItems(patternpaths)
	if len(patternpaths) != 0 {
		currentCollectiveChecksum = utility.CollectiveSHA256Hash(patternpaths)
	}
	return currentCollectiveChecksum
}
