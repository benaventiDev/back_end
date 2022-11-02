package main

import (
	"fmt"
	"unsafe"
)

func Test1() {
	var part_start int32
	//var part_size int32
	if LoggedInPartition.Part_type == Primary {
		part_start = LoggedInPartition.Part.Part_start
		//part_size = LoggedInPartition.Part.Part_size
	} else if LoggedInPartition.Part_type == Logic {
		part_start = LoggedInPartition.Logica.Part_start
		//part_size = LoggedInPartition.Logica.Part_size
	}
	sb, err := ReadSB(LoggedInPartition.Path, part_start)
	if err != "" {
		fmt.Println("Error reading SB")
	}
	bitmapBlock := ReadBitmapBlock(LoggedInPartition.Path, sb.S_bm_block_start, sb.S_blocks_count)
	bitmapInode := ReadBitmapInode(LoggedInPartition.Path, sb.S_bm_inode_start, sb.S_inodes_count)
	/*
		err1 := os.WriteFile("block.txt", bitmapBlock, 0644)
		if err1 != nil {
			fmt.Println("Error reading SB")
		}
		err2 := os.WriteFile("bitmap.txt", bitmapInode, 0644)
		if err2 != nil {
			fmt.Println("Error reading SB")
		}*/
	fmt.Println("bitmapBlock")
	for i := 0; i < 10; i++ {
		fmt.Printf("%v ", bitmapBlock[i])
	}
	fmt.Println()
	fmt.Println("bitmapInode")
	for i := 0; i < 10; i++ {
		fmt.Printf("%v ", bitmapInode[i])
	}
	fmt.Println()

	root_inode, err := ReadInode(LoggedInPartition.Path, sb.S_inode_start)
	if err != "" {
		fmt.Println("Error reading root_inode")
	}
	fmt.Println("root_inode")
	PrettyStruct(root_inode)
	for i := 0; i < 15; i++ {
		if root_inode.I_block[i] != -1 {
			PrintFolderBlocks(sb, root_inode.I_block[i])
		}
	}

}

func PrintFolderBlocks(sb Superbloque, index int32) {
	fb := FolderBlock{}
	fb, err := ReadFolderBlock(LoggedInPartition.Path, sb.S_block_start+index*int32(unsafe.Sizeof(fb)))
	if err != "" {
		fmt.Println("Error reading folder block")
	}
	fmt.Println("\nFolderBlock")
	for i := 0; i < 4; i++ {
		if fb.B_content[i].B_inodo != -1 && Byte12ToString(fb.B_content[i].B_name) != ".." && Byte12ToString(fb.B_content[i].B_name) != "." {
			fmt.Printf("Name: %v, Inode: %v\n", Byte12ToString(fb.B_content[i].B_name), fb.B_content[i].B_inodo)
			inode := Inode{}
			inode, err := ReadInode(LoggedInPartition.Path, sb.S_inode_start+fb.B_content[i].B_inodo*int32(unsafe.Sizeof(inode)))
			if err != "" {
				fmt.Println("PrintFolderBlocks:Error reading inode")
			}
			PrettyStruct(inode)
			if inode.Inode_type == Inode_folder {
				PrintFolderBlocks(sb, inode.I_block[0])
			}

		}
	}

}
