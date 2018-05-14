package nav

import (
	"fmt"
	"testing"
	"zeus/linmath"
)

func Test_loadMesh(t *testing.T) {

	mesh := NewMesh("D:\\Work\\NavMesh\\Assets\\map_export\\nav_test\\navmesh.bin")
	finder := NewMeshPathFinder(mesh)

	path, _ := finder.FindPath(linmath.NewVector3(38, 1000, 72), linmath.NewVector3(61, 1000, 30))

	fmt.Println(path)
}

func BenchmarkGetPoint(b *testing.B) {

}

// func BenchmarkFindPath(b *testing.B) {
// 	mesh := NewMesh("D:\\Workspace\\wktimefire\\res\\space\\1\\navmesh.bin")
// 	finder := NewMeshPathFinder(mesh)
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		x1 := r.Float32()
// 		z1 := r.Float32()
// 		x2 := r.Float32()
// 		z2 := r.Float32()

// 		//finder.FindPath(x1, z1, x2, z2)
// 		if _, err := finder.FindPath(x1, z1, x2, z2); err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// }

func Test_basic(t *testing.T) {
	b := make(chan bool)

	go func() {
		<-b
		fmt.Println("hi")
		b <- true
	}()

	fmt.Println("begin")
	b <- true
	fmt.Println("end")

	<-b
	fmt.Println("xxx")
}
