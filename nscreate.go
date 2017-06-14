// Renders the yaml file for a new namespace and yaml files for the associated resourcequota and
// rolebinding.

// It then outputs a bash script that can be run to add the [rendered] items above to a kubernetes
// cluster. This is a simple proof of concept that is meant to use go's standard library, but would
// likely be replace by somem parsing library or kubernetes packagine tool such as Helm.

package main

import (
	"bufio"
	"bytes"
	"flag"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type Meta struct {
	Namespace   string
	Group       string
	Role        string
	LimitCPU    string
	LimitMemory string
}

func main() {
	ns := flag.String("ns", "", "namespace name to use")
	group := flag.String("group", "", "ldap group to associate")
	role := flag.String("role", "", "rbac clusterrole to bind to")
	cpu := flag.String("cpu", "", "cpu limit for namespace")
	mem := flag.String("mem", "", "memory limit for namespace")
	flag.Parse()

	meta := Meta{
		Namespace:   *ns,
		Group:       *group,
		Role:        *role,
		LimitCPU:    *cpu,
		LimitMemory: *mem,
	}
	meta.validate()

	manifests := listManifests()
	for _, m := range manifests {
		o := meta.parseTemplate(m)
		writeFiles(meta.Namespace, o, m)
	}

	log.Printf("Render complete. \n\nTo deploy, run:\n$ kubectl apply -f %s/namespace.yaml && kubectl apply -f %s/", meta.Namespace, meta.Namespace)
}

func (m Meta) validate() {
	switch {
	case m.Namespace == "":
		log.Fatalln("Must specify namespace flag. Use --help to see params.")
	case m.Group == "":
		log.Fatalln("Must specify group flag. Use --help to see params.")
	case m.Role == "":
		log.Fatalln("Must specify role flag. Use --help to see params.")
	case m.LimitCPU == "":
		log.Fatalln("Must specify cpu flag. Use --help to see params.")
	case m.LimitMemory == "":
		log.Fatalln("Must specify memory flag. Use --help to see params.")
	}
}

func (m Meta) parseTemplate(fileName string) string {
	dirname := "." + string(filepath.Separator) + "manifests" + string(filepath.Separator)
	t, err := template.ParseFiles(dirname + fileName)
	if err != nil {
		log.Fatalf("Parsing yaml template failed. File: %s Error: %s.", fileName, err.Error())
	}
	buff := bytes.NewBufferString("")
	if err = t.Execute(buff, m); err != nil {
		log.Fatalf("Failed to execute parse on template. File: %s Error: %s.", fileName, err.Error())
	}

	return buff.String()
}

func writeFiles(namespace string, manifest string, fileName string) {
	if err := os.MkdirAll(namespace, 0777); err != nil {
		log.Fatalf("Failed to create output dir for rendered manifest. Namespace: %s Error: %s.", namespace, err.Error())
	}
	f, err := os.Create(namespace + string(filepath.Separator) + fileName)
	if err != nil {
		log.Fatalf("Failed to create new manifest file. File: %s, Error: %s.", fileName, err.Error())
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(manifest)
	w.Flush()
}

func listManifests() []string {
	dirname := "." + string(filepath.Separator) + "manifests"
	d, err := os.Open(dirname)
	if err != nil {
		log.Fatalf("Failed to find manifests directory. Error: %s.", err.Error())
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		log.Fatalf("Failed to read files from manifests directory. Error: %s.", err.Error())
	}

	ms := []string{}
	for _, file := range files {
		if !file.Mode().IsRegular() {
			continue
		}
		if filepath.Ext(file.Name()) == ".yaml" {
			ms = append(ms, file.Name())
		}
	}

	return ms
}
