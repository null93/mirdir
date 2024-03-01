package template

import (
	"bytes"
	"fmt"
	. "github.com/null93/mirdir/pkg/utils"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"text/template"
)

type PathType string

const (
	File      PathType = "file"
	Directory PathType = "directory"
	Link      PathType = "link"
)

type Template struct {
	TemplateDirectory string
	InputPaths        []Path
}

type Path struct {
	Type        PathType
	Path        string
	Content     []byte
	OwnerUid    uint32
	OwnerGid    uint32
	Permissions os.FileMode
	LinkTarget  string
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
	for _, input := range tmpl.InputPaths {
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
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}

func (p *Path) Print() {
	if p.Type == Directory {
		fmt.Printf("%s %4d:%-4d %s/\n", p.Permissions.String(), p.OwnerUid, p.OwnerGid, p.Path)
	}
	if p.Type == File {
		fmt.Printf("%s %4d:%-4d %s\n%s\n", p.Permissions.String(), p.OwnerUid, p.OwnerGid, p.Path, string(p.Content))
	}
	if p.Type == Link {
		fmt.Printf("%s %4d:%-4d %s -> %s\n", p.Permissions.String(), p.OwnerUid, p.OwnerGid, p.Path, p.LinkTarget)
	}
}

func (p *Path) Write(preserveOwnership bool) error {
	if p.Type == Directory {
		if err := os.MkdirAll(p.Path, p.Permissions); err != nil {
			return err
		}
	}
	if p.Type == File {
		if err := os.WriteFile(p.Path, p.Content, p.Permissions); err != nil {
			return err
		}
	}
	if p.Type == Link {
		if Exists(p.Path) {
			if err := os.Remove(p.Path); err != nil {
				return err
			}
		}
		if err := os.Symlink(p.LinkTarget, p.Path); err != nil {
			return err
		}
	}
	if preserveOwnership {
		if err := os.Lchown(p.Path, int(p.OwnerUid), int(p.OwnerGid)); err != nil {
			return err
		}
	}
	return nil
}

func (p *Path) IsTemplate() bool {
	return strings.HasSuffix(p.Path, ".tpl")
}

func (p *Path) GetRenderedContent(values interface{}) ([]byte, error) {
	if p.IsTemplate() && p.Type == File {
		tmpl, parseErr := template.New("").Parse(string(p.Content))
		if parseErr != nil {
			return nil, parseErr
		}
		var renderedContentBuffer bytes.Buffer
		execErr := tmpl.Execute(&renderedContentBuffer, values)
		if execErr != nil {
			return nil, execErr
		}
		return renderedContentBuffer.Bytes(), nil
	}
	return p.Content, nil
}

func (p *Path) GetRenderedLinkTarget(values map[string]string) string {
	if p.Type != Link {
		return p.LinkTarget
	}
	result := p.LinkTarget
	for key, value := range values {
		result = strings.ReplaceAll(result, "["+key+"]", value)
	}
	return result
}

func (p *Path) GetRenderedPath(values map[string]string) string {
	result := p.Path
	for key, value := range values {
		result = strings.ReplaceAll(result, "["+key+"]", value)
	}
	return strings.TrimSuffix(result, ".tpl")
}

func (p *Path) IsDir() bool {
	return p.Type == Directory
}
