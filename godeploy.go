package main

import (
    "fmt"
    "strings"
    "bytes"
    "log"
    "io/ioutil"
    "encoding/json"
    "net/http"
    "os"
    "os/exec"
)

func deployHandler(w http.ResponseWriter, r *http.Request) {
    buf := new(bytes.Buffer)
    defer func() {
        if _, err := buf.WriteTo(w); err != nil {
            log.Printf("WriteTo: %v", err)
        }
    }()

    requestBody, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError);
        return
    }

    var f interface{}
    if err := json.Unmarshal(requestBody, &f); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest);
        return
    }

    m := f.(map[string]interface{})
    event := r.Header.Get("X-GitHub-Event")

    switch event {
    case "ping":
        fmt.Fprintf(buf, "ping acknowledged from hook_id: %.0f\n", m["hook_id"])

    case "push":
        if err := r.ParseForm(); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        path := r.Form.Get("path")

        fmt.Fprintf(buf, "path: %s\n", path)

        if err := os.Chdir(path); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        output, err := exec.Command("git", "symbolic-ref", "-q", "HEAD").Output()
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        branch := strings.TrimSpace(string(output))
        if branch == m["ref"].(string) {
            output, _:= exec.Command("git", "pull").CombinedOutput();
            
            fmt.Fprintf(buf, "output:\n%s\n", output)
            log.Printf("\npath: %s\nbranch: %s\noutput:\n%s====", path, branch, output)
        }

    default:
        http.Error(w, fmt.Sprintf("Unknown event type: %s", event), http.StatusNotImplemented)
        return
    }
}

func main() {
    f, err := os.OpenFile("/var/log/godeploy", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()

    log.SetOutput(f)

    http.HandleFunc("/", deployHandler)
    http.ListenAndServe(":8001", nil)
}