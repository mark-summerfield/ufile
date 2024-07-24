// Copyright © 2024 Mark Summerfield. All rights reserved.

package ufile

import (
    "fmt"
    _ "embed"
    )

//go:embed Version.dat
var Version string

func Hello() string {
    return fmt.Sprintf("Hello ufile v%s", Version)
}
