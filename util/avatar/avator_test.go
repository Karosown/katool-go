package avatar

import (
	"fmt"
	"testing"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/karosown/katool-go/util/pathutil"
)

func TestAvatar(t *testing.T) {

	avatar, err := NewAvatarGenerator().SetFont("google-sans").SetColors("#ffffff", "#5c6bc0").GenerateAvatarFile("KaTool")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(avatar.File.Name())
	path := pathutil.NewWrapper(fileutil.CurrentPath()).AfterLayer("/output")
	path.ForceCreate()
	err = fileutil.WriteBytesToFile(path.Path+avatar.FileName, avatar.File.Content())
	if err != nil {
		t.Fatal(err)
	}
}
