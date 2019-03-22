//+build test_live

package core

import (
	"context"
	"encoding/json"
	"os/exec"
	"sort"
	"sync"
	"testing"

	jmespath "github.com/jmespath/go-jmespath"
	"github.com/object88/slog/internal/constants"
	"github.com/object88/slog/kubernetes/client"
	v1 "k8s.io/api/core/v1"
)

func Test_Watcher_Live_GetPods(t *testing.T) {
	// First, use the Kubectl command to get the pods (count and names)

	var wg sync.WaitGroup
	wg.Add(1)

	cmd := exec.CommandContext(context.Background(), "kubectl", "get", "pods", "-o", "json")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal("Failed to get STDOUT pipe")
		return
	}

	var names []string

	go func() {
		defer wg.Done()
		var data interface{}
		input := json.NewDecoder(stdout)
		input.Decode(&data)

		jmes, err := jmespath.Compile(`items[].metadata.name`)
		if err != nil {
			t.Fatal("Failed to compile jmespath query")
			return
		}

		results, err := jmes.Search(data)
		if err != nil {
			t.Fatalf("Error while searching JSON: %s", err.Error())
		}

		intSlice, ok := results.([]interface{})
		if !ok {
			t.Fatalf("Internal error: did not get []interface{}: %#v", results)
		}

		names = make([]string, len(intSlice))
		for k, v := range intSlice {
			name, ok := v.(string)
			if !ok {
				t.Fatalf("Internal error: did not get string at offset %d", k)
			}
			names[k] = name
		}
	}()

	err = cmd.Run()
	if err != nil {
		t.Fatalf("Error while execiting kubectl: %s", err.Error())
	}

	// Wait for the command to finish, and sort the found names
	wg.Wait()
	sort.Strings(names)

	// Now, set up the watcher to find pods
	done := make(chan int)
	results := []string{}
	w := NewWatcher()
	err = w.Connect(client.NewClientFactory(buildFactory()), "default")
	if err != nil {
		t.Fatal("Failed to connect watcher")
		return
	}

	go func() {
		for e := range w.Listen() {
			name := e.Object.(*v1.Pod).GetName()
			results = append(results, name)
		}
		done <- 1
	}()

	err = w.Load(constants.Pods)
	if err != nil {
		t.Fatal("Failed to load")
		return
	}

	// Stop and wait
	w.Stop()
	<-done

	// Check the results.
	sort.Strings(results)
	for k, v := range results {
		if v != names[k] {
			t.Errorf("Mismatched")
		}
	}
}
