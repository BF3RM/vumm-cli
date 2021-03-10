package workspace

import (
	"testing"
)

type modDependencyTest struct {
	name, version  string
	expectedTag    string
	hasConstraints bool
}

type modDependencyStringTest struct {
	nameAndVersion string
	expectedTag    string
	hasConstraints bool
}

func TestResolveModDependency(t *testing.T) {
	tests := []modDependencyTest{
		{"realitymod", "latest", "latest", false},
		{"realitymod", "0.1.0", "", true},
		{"realitymod", ">=0.1.0", "", true},
	}

	for _, test := range tests {
		dep := ResolveModDependency(test.name, test.version)
		if dep.Tag != test.expectedTag {
			t.Errorf("Tag %s not equal to expected %s", dep.Tag, test.expectedTag)
		}
		hasConstraints := dep.VersionConstraints != nil
		if hasConstraints != test.hasConstraints {
			t.Errorf("VersionConstraints %v not equal to expected %v", hasConstraints, test.hasConstraints)
		}
	}
}

func TestResolveModDependencyFromString(t *testing.T) {
	tests := []modDependencyStringTest{
		{"realitymod@latest", "latest", false},
		{"realitymod@0.1.0", "", true},
		{"realitymod@>=0.1.0", "", true},
	}

	for _, test := range tests {
		dep := ResolveModDependencyFromString(test.nameAndVersion)
		if dep.Tag != test.expectedTag {
			t.Errorf("Tag %s not equal to expected %s", dep.Tag, test.expectedTag)
		}
		hasConstraints := dep.VersionConstraints != nil
		if hasConstraints != test.hasConstraints {
			t.Errorf("VersionConstraints %v not equal to expected %v", hasConstraints, test.hasConstraints)
		}
	}
}
