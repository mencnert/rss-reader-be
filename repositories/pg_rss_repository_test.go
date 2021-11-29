package repositories

import (
	"testing"
)

func TestInterfaceImpl(t *testing.T) {
	var _ RssRepository = &pgRssRepository{}
}
