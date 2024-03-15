package template

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	. "github.com/null93/mirdir/pkg/utils"
)

type PathType string

const (
	File      PathType = "file"
	Directory PathType = "directory"
	Link      PathType = "link"
)

type Path struct {
	Type        PathType
	Path        string
	Content     []byte
	OwnerUid    uint32
	OwnerGid    uint32
	Permissions os.FileMode
	LinkTarget  string
	Input       *Path
}

func (p *Path) Print(isWrite bool, verbose bool) {
	if p.Type == Directory {
		fmt.Printf("TPL %s %4d:%-4d %s/\n", p.Input.Permissions.String(), p.Input.OwnerUid, p.Input.OwnerGid, p.Input.Path)
		fmt.Printf("DST %s %4d:%-4d %s/\n", p.Permissions.String(), p.OwnerUid, p.OwnerGid, p.Path)
	}
	if p.Type == File {
		fmt.Printf("TPL %s %4d:%-4d %s\n", p.Input.Permissions.String(), p.Input.OwnerUid, p.Input.OwnerGid, p.Input.Path)
		if verbose {
			fmt.Printf(p.Input.GetDebugContent())
		}
		if isWrite {
			fmt.Printf("DST %s %4d:%-4d %s\n", p.Permissions.String(), p.OwnerUid, p.OwnerGid, p.Path)
			if verbose {
				fmt.Printf(p.GetDebugContent())
			}
		} else {
			fmt.Println("DST <deleted>")
		}
	}
	if p.Type == Link {
		fmt.Printf("TPL %s %4d:%-4d %s -> %s\n", p.Input.Permissions.String(), p.Input.OwnerUid, p.Input.OwnerGid, p.Input.Path, p.Input.LinkTarget)
		fmt.Printf("DST %s %4d:%-4d %s -> %s\n", p.Permissions.String(), p.OwnerUid, p.OwnerGid, p.Path, p.LinkTarget)
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

func (p *Path) Delete() error {
	if p.Type != File {
		return errors.New("will only delete files")
	}
	if err := os.Remove(p.Path); err != nil {
		return err
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

func (p *Path) IsEmptyFile() bool {
	trimmed := strings.TrimSpace(string(p.Content))
	return p.Type == File && len(trimmed) == 0
}

func (p *Path) GetDebugContent () string {
	lines := strings.Split(string(p.Content), "\n")
	result := ""
	for _, line := range lines {
		result += fmt.Sprintf("    %s\n", line)
	}
	return result
}
