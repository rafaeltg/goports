package domaintest

import (
	"fmt"

	"github.com/rafaeltg/goports/internal/core/domain"
)

type portMatcher struct {
	expected *domain.Port
}

func (m *portMatcher) Matches(x interface{}) bool {
	actual, ok := x.(*domain.Port)
	if !ok {
		return false
	}

	return (m.expected.ID == actual.ID) &&
		(m.expected.Name == actual.Name) // TODO: compare other fields
}

func (m *portMatcher) String() string {
	return fmt.Sprintf("%v", m.expected)
}

func PortMatcher(expected *domain.Port) *portMatcher {
	return &portMatcher{expected}
}

type portsMatcher struct {
	expected domain.Ports
}

func (m *portsMatcher) Matches(x interface{}) bool {
	actual, ok := x.(domain.Ports)
	if !ok || len(actual) != len(m.expected) {
		return false
	}

	for i := range m.expected {
		if !PortMatcher(&m.expected[i]).Matches(&actual[i]) {
			return false
		}
	}

	return true
}

func (m *portsMatcher) String() string {
	return fmt.Sprintf("%v", m.expected)
}

func PortsMatcher(expected domain.Ports) *portsMatcher {
	return &portsMatcher{expected}
}
