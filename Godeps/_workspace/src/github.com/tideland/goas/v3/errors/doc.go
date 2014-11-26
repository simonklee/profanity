// Tideland Go Application Support - Errors
//
// Copyright (C) 2013-2014 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Typical errors in Go are often created using errors.New()
// or fmt.Errorf(). Those errors only contain a string as information.
// When trying to differentiate between errors or to carry helpful
// payload own types are needed.
//
// The errors package allows to easily created formatted errors
// with New() like with the fmt.Errorf function, but also with an
// error code. Additionlly a Messages instance has to be passed
// to map the error code to their according messages.
//
// If an error alreay exists use Annotate(). This way the original
// error will be stored and can be retrieved with Annotated(). Also
// its error message will be appended to the created error separated
// by a colon.
//
// All errors additionally contain their package, filename and line
// number. These information can be retrieved using Location(). In
// case of a chain of annotated errors those can be retrieved as a
// slice of errors with Stack().
package errors

//--------------------
// IMPORTS
//--------------------

import (
	"github.com/tideland/goas/v1/version"
)

//--------------------
// VERSION
//--------------------

// PackageVersion returns the version of the version package.
func PackageVersion() version.Version {
	return version.New(3, 2, 0)
}

// EOF
