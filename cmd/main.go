/* SPDX-License-Identifier: Apache-2.0 */
/* Copyright Contributors to the vidx2pidx project. */

package main

// main is the main entrypoint of vidx2pidx
func main() {
	cmd := NewCli()
	ExitOnError(cmd.Execute())
}
