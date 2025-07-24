//******************************************************************
//
// Try out neighbour search for all ST stypes together
//
// Prepare:
// cd examples
// ../src/N4L-db -u chinese.n4l
//
//******************************************************************

package main

import (
	"fmt"

        SST "SSTorytime"
)

var path [8][]string

//******************************************************************

func main() {

	load_arrows := true
	ctx := SST.Open(load_arrows)

	Solve(ctx)

	SST.Close(ctx)
}

//******************************************************************

func Solve(ctx SST.PoSST) {

	// Contra colliding wavefronts as path integral solver

	const maxdepth = 16

	start_bc := "i6"

	p1 := SST.GetDBNodePtrMatchingName(ctx,start_bc,"")
	p2 := SST.GetDBNodePtrMatchingNCC(ctx,start_bc,"",nil,nil)

	if Diff (p1,p2) {
		fmt.Println("Failed",p1,p2)
	}
}

// **********************************************************

func Diff(left,right []SST.NodePtr) bool {

	// Return coordinate pairs of partial paths to splice

	if len(left) != len(right) {
		return true
	}

	for l := 0; l < len(left); l++ {
		if left[l] != right[l] {
			fmt.Println("Mismatch:",left[l],right[l])
			return true
		}
	}

	return false
}







