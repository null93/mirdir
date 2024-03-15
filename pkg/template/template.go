package template

import (
	"os"
	"path/filepath"
	"sort"
	"syscall"
)

type Template struct {
	TemplateDirectory string
	InputPaths        []Path
}

func NewTemplate(tmplDir string) (*Template, error) {
	inputs := []Path{}
	walkErr := filepath.Walk(tmplDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, relErr := filepath.Rel(tmplDir, path)
		if relErr != nil {
			return relErr
		}
		infosys := info.Sys()
		if info.IsDir() {
			inputs = append(inputs, Path{
				Type:        Directory,
				Path:        relPath,
				OwnerUid:    infosys.(*syscall.Stat_t).Uid,
				OwnerGid:    infosys.(*syscall.Stat_t).Gid,
				Permissions: info.Mode().Perm(),
			})
			return nil
		}
		if info.Mode()&os.ModeSymlink != 0 {
			linkTarget, readlinkErr := os.Readlink(path)
			if readlinkErr != nil {
				return readlinkErr
			}
			inputs = append(inputs, Path{
				Type:        Link,
				Path:        relPath,
				OwnerUid:    infosys.(*syscall.Stat_t).Uid,
				OwnerGid:    infosys.(*syscall.Stat_t).Gid,
				Permissions: info.Mode().Perm(),
				LinkTarget:  linkTarget,
			})
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		inputs = append(inputs, Path{
			Type:        File,
			Path:        relPath,
			Content:     content,
			OwnerUid:    infosys.(*syscall.Stat_t).Uid,
			OwnerGid:    infosys.(*syscall.Stat_t).Gid,
			Permissions: info.Mode().Perm(),
		})
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}
	sort.Slice(inputs, func(i, j int) bool { return inputs[i].Type < inputs[j].Type })
	tmpl := Template{
		TemplateDirectory: tmplDir,
		InputPaths:        inputs,
	}
	return &tmpl, nil
}

func (tmpl *Template) Render(destDir string, preservePerms bool, values map[string]string) ([]Path, error) {
	outputs := []Path{}
	for key, input := range tmpl.InputPaths {
		mode := input.Permissions
		if !preservePerms && (input.Type == File || input.Type == Link) {
			mode = os.FileMode(0666) // value before umask is applied
		}
		if !preservePerms && input.Type == Directory {
			mode = os.FileMode(0777) // value before umask is applied
		}
		renderedPath := input.GetRenderedPath(values)
		renderedLinkTarget := input.GetRenderedLinkTarget(values)
		renderedContent, renderErr := input.GetRenderedContent(values)
		if renderErr != nil {
			return nil, renderErr
		}
		output := Path{
			Type:        input.Type,
			Path:        filepath.Join(destDir, renderedPath),
			Content:     renderedContent,
			LinkTarget:  renderedLinkTarget,
			OwnerUid:    input.OwnerUid,
			OwnerGid:    input.OwnerGid,
			Permissions: mode,
			Input:       &tmpl.InputPaths[key],
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}
