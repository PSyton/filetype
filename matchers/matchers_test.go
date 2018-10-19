package matchers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/filetype.v1/types"
)

func TestMatch(t *testing.T) {

	type s struct {
		T types.Type
		F TypeMatcher
	}

	// Prepare all matchers
	testMathers := map[string]s{}
	for mType, aMatcher := range Matchers {
		testMathers[mType.Extension] = s{
			T: mType,
			F: aMatcher,
		}
	}

	err := filepath.Walk("../fixtures", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext[0:1] == "." {
			ext = ext[1:]
		}

		matcher, ok := testMathers[ext]

		if !ok {
			assert.Fail(t, fmt.Sprintf("No matcher for fixture %s", path))
			return nil
		}

		buff, err := ioutil.ReadFile(path)

		if err != nil {
			assert.Fail(t, fmt.Sprintf("Read error: %s", path))
			return nil
		}

		assert.Equal(t, matcher.T, matcher.F(buff), fmt.Sprintf("Match failed for %s", path))

		return nil
	})

	if err != nil {
		fmt.Print(err)
	}

	assert.Nil(t, err, "Failed walt fixtures")
}
