package main

import "fmt"

type InterfaceNode interface {
	Print(string)
	Clone() InterfaceNode
}

// 文件类
type File struct {
	Name string
}

func (f *File) Print(indentation string) {
	fmt.Println(indentation + f.Name)
}

func (f *File) Clone() InterfaceNode {
	return &File{Name: f.Name + "_Clone"}
}

// 文件夹类
type Folder struct {
	Children []InterfaceNode
	Name     string
}

func (f *Folder) Print(indentation string) {
	fmt.Println(indentation + f.Name)
	for _, i := range f.Children {
		i.Print(indentation + indentation)
	}
}

func (f *Folder) Clone() InterfaceNode {
	cloneFolder := &Folder{
		Name: f.Name + "_Clone",
	}
	var tmpChildren []InterfaceNode
	for _, child := range f.Children {
		cp := child.Clone()
		tmpChildren = append(tmpChildren, cp)
	}
	cloneFolder.Children = tmpChildren
	return cloneFolder
}

func main() {
	f1 := &File{Name: "file1"}
	f2 := &File{Name: "file2"}
	f3 := &File{Name: "file3"}

	folder1 := &Folder{
		Children: []InterfaceNode{f1},
		Name:     "文件夹Folder1",
	}

	folder2 := &Folder{
		Children: []InterfaceNode{folder1, f2, f3},
		Name:     "文件夹Folder2",
	}

	fmt.Println("\n打印文件夹Folder2的层级：")
	folder2.Print("  ")

	cloneFolder := folder2.Clone()
	fmt.Println("\n打印复制文件夹Folder2的层级：")
	cloneFolder.Print("  ")
}
