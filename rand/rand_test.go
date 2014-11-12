package rand_test

import (
	"bitbucket.org/syb-devs/gotools/rand"
	"testing"
)

var alphaNumTests = []struct {
	seed     int64
	size     int
	expected []string
}{
	{442342545, 0, []string{"", ""}},
	{352672, 8, []string{"a7jDf73H", "WEHe5CG1", "dA8iZGRf"}},
	{1, 128, []string{"RFbD56TI2smTyVsGd5Xav0yu99ZAMPTA7z7s575klKiz9pyKl17ltLSvQmntzYlkmifsd2X28mLGpj0sdzvhNhjpXmkI0TwXU3Pqj71n5gwFoigtDswxBrlgaWd1iKZU"}},
}

func TestAlphaNum(t *testing.T) {
	for _, test := range alphaNumTests {
		rand.Seed(test.seed)
		for _, expected := range test.expected {
			actual := rand.AlphaNum(test.size)
			if actual != expected {
				t.Errorf("AlphaNum(%v) with seed %v, expecting %v, got %v", test.size, test.seed, expected, actual)
			}
		}
	}
}
