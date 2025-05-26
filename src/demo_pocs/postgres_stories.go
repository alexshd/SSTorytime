//******************************************************************
//
// Demo of node by node addition
//
//******************************************************************

package main

import (
	"fmt"
        SST "SSTorytime"
)

//******************************************************************

func main() {

	load_arrows := false
	ctx := SST.Open(load_arrows)

	fmt.Println("Reset..")

	var n1,n2,n3,n4,n5,n6 SST.Node
	var lnk12,lnk23,lnk34,lnk25,lnk56 SST.Link
	
	nd := SSTInsertNode("Mary had a little lamb","home and away")

	fmr.Println(nd)

	n2.NPtr = SST.NodePtr{ CPtr : 2, Class: SST.LT128}
	n2.S = "Whose fleece was dull and grey"
	n2.Chap = "home and away"

	n3.NPtr = SST.NodePtr{ CPtr : 3, Class: SST.LT128}
	n3.S = "And everytime she washed it clean"
	n3.Chap = "home and away"

	n4.NPtr = SST.NodePtr{ CPtr : 4, Class: SST.LT128}
	n4.S = "It just went to roll in the hay"
	n4.Chap = "home and away"

	n5.NPtr = SST.NodePtr{ CPtr : 5, Class: SST.LT128}
	n5.S = "And every bar that Mary went"
	n5.Chap = "home and away"

	n6.NPtr = SST.NodePtr{ CPtr : 6, Class: SST.LT128}
	n6.S = "Was hot and loud and gay"
	n6.Chap = "home and away"

	lnk12.Arr = 77
	lnk12.Wgt = 0.34
	lnk12.Ctx = []string{"fairy castles","angel air"}
	lnk12.Dst = n2.NPtr

	lnk23.Arr = 77
	lnk23.Wgt = 0.34
	lnk23.Ctx = []string{"fairy castles","angel air"}
	lnk23.Dst = n3.NPtr

	lnk34.Arr = 77
	lnk34.Wgt = 0.34
	lnk34.Ctx = []string{"fairy castles","angel air"}
	lnk34.Dst = n4.NPtr

	lnk25.Arr = 77
	lnk25.Wgt = 0.34
	lnk25.Ctx = []string{"steamy hot tubs"}
	lnk25.Dst = n5.NPtr

	lnk56.Arr = 77
	lnk56.Wgt = 0.34
	lnk56.Ctx = []string{"steamy hot tubs","lady gaga"}
	lnk56.Dst = n6.NPtr

	sttype := SST.LEADSTO

	n1 = SST.IdempDBAddNode(ctx, n1)
	n2 = SST.IdempDBAddNode(ctx, n2)
	SST.AppendDBLinkToNode(ctx,n1.NPtr,lnk12,sttype)

	n2 = SST.IdempDBAddNode(ctx, n2)
	n3 = SST.IdempDBAddNode(ctx, n3)
	SST.AppendDBLinkToNode(ctx,n2.NPtr,lnk23,sttype)

	n3 = SST.IdempDBAddNode(ctx, n3)
	n4 = SST.IdempDBAddNode(ctx, n4)
	SST.AppendDBLinkToNode(ctx,n3.NPtr,lnk34,sttype)

	n2 = SST.IdempDBAddNode(ctx, n2)
	n5 = SST.IdempDBAddNode(ctx, n5)
	SST.AppendDBLinkToNode(ctx,n2.NPtr,lnk25,sttype)

	n5 = SST.IdempDBAddNode(ctx, n5)
	n6 = SST.IdempDBAddNode(ctx, n6)
	SST.AppendDBLinkToNode(ctx,n5.NPtr,lnk56,sttype)

	fmt.Println("----------------------------------")
	fmt.Println("Node section hypersurfaces:")

	const maxdepth = 8

	for depth := 0; depth < maxdepth; depth++ {
		val := SST.GetFwdConeAsNodes(ctx,n1.NPtr,sttype,depth)
		fmt.Println("As NodePtr(s) fwd from",n1,"depth",depth)
		for l := range val {
			fmt.Println("   - Step",val[l])
		}

	}

	fmt.Println("----------------------------------")
	fmt.Println("Link section hypersurfaces:")

	for depth := 0; depth < maxdepth; depth++ {
		val := SST.GetFwdConeAsLinks(ctx,n1.NPtr,sttype,depth)
		fmt.Println("Search as Links fwd from",n1,"depth",depth)
		for l := range val {
			fmt.Println("   - Step",val[l])
		}
	}

	fmt.Println("----------------------------------")
	fmt.Println("Link proper time normal paths:")

	for depth := 0; depth < maxdepth; depth++ {

		fmt.Println("Searching paths of length",depth,"/",maxdepth,"from",n1.NPtr)

		paths,_ := SST.GetFwdPathsAsLinks(ctx,n1.NPtr,sttype,depth)

		for p := range paths {

			if len(paths[p]) > 1 {
			
				fmt.Println("    Path",p," len",len(paths[p]))

				for l := 0; l < len(paths[p]); l++ {
					fmt.Println("    ",l,"xx  --> ",paths[p][l].Dst,"weight",paths[p][l].Wgt)
				}
			}
		}
	}

	SST.Close(ctx)
}



//*****************************************

func SSTInsertNode(s string,chapter string) {

	var node SST.Node

	node.S = s
	node.Chap = chapter  
        node.L,node.NPtr.Class = StorageClass(n.S)
	node.NPtr.CPtr = DB_NODE_OFFSETS[node.NPtr.Class]
	DB_NODE_OFFSETS[node.NPtr.Class]++

}






