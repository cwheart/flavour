package flavour

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlavour(t *testing.T) {
	f := New()
	assert.NotNil(t, f)
}
