package start

import (
	"fmt"

	"github.com/josephbudd/okp/shared/paths"
)

func initPaths() (err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("startup.initPaths: %w", err)
		}
	}()

	err = paths.Init()
	return
}
