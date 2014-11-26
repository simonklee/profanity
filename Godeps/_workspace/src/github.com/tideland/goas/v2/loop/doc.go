// Tideland Go Application Support - Loop
//
// Copyright (C) 2013-2014 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// A typical Go idiom for concurrent applications is running
// a loop in the background doing a select on one or more channels.
// Stopping those loops or getting aware of internal errors
// requires extra efforts. The loop package helps to manage
// this kind of goroutines.
package loop

// EOF
