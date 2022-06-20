package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	policyFname = flag.String("policy-file", "./policy.hujson", "filename for policy file")
	timeout     = flag.Duration("timeout", 5*time.Minute, "timeout for the entire CI run")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	tailnet := os.Getenv("TS_TAILNET")
	apiKey := os.Getenv("TS_API_KEY")

	switch flag.Arg(0) {
	case "apply":
		controlEtag, err := getACLETag(ctx, tailnet, apiKey)
		if err != nil {
			log.Fatal(err)
		}

		localEtag, err := sumFile(*policyFname)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("control: %s", controlEtag)
		log.Printf("local:   %s", localEtag)

		if controlEtag == localEtag {
			log.Println("no update needed, doing nothing")
			os.Exit(0)
		}

		if err := applyNewACL(ctx, tailnet, apiKey, *policyFname, controlEtag); err != nil {
			log.Fatal(err)
		}

	case "test":
		controlEtag, err := getACLETag(ctx, tailnet, apiKey)
		if err != nil {
			log.Fatal(err)
		}

		localEtag, err := sumFile(*policyFname)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("control: %s", controlEtag)
		log.Printf("local:   %s", localEtag)

		if controlEtag == localEtag {
			log.Println("no updates found, doing nothing")
			os.Exit(0)
		}

		if err := testNewACLs(ctx, tailnet, apiKey, *policyFname); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("usage: %s [options] <test|apply>", os.Args[0])
	}
}

func sumFile(fname string) (string, error) {
	fin, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer fin.Close()

	h := sha256.New()
	_, err = io.Copy(h, fin)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\"%x\"", h.Sum(nil)), nil
}

func applyNewACL(ctx context.Context, tailnet, apiKey, policyFname, oldEtag string) error {
	fin, err := os.Open(policyFname)
	if err != nil {
		return err
	}
	defer fin.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.tailscale.com/api/v2/tailnet/%s/acl", tailnet), fin)
	if err != nil {
		return err
	}

	req.SetBasicAuth(apiKey, "")
	req.Header.Set("Content-Type", "application/hujson")
	req.Header.Set("If-Match", oldEtag)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("wanted HTTP status code %d but got %d", http.StatusOK, resp.StatusCode)
	}

	return nil
}

func testNewACLs(ctx context.Context, tailnet, apiKey, policyFname string) error {
	fin, err := os.Open(policyFname)
	if err != nil {
		return err
	}
	defer fin.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.tailscale.com/api/v2/tailnet/%s/acl/validate", tailnet), fin)
	if err != nil {
		return err
	}

	req.SetBasicAuth(apiKey, "")
	req.Header.Set("Content-Type", "application/hujson")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var ate ACLTestError
		err := json.NewDecoder(resp.Body).Decode(&ate)
		if err != nil {
			return err
		}

		return ate
	}

	return nil
}

type ACLTestError struct {
	Message string               `json:"message"`
	Data    []ACLTestErrorDetail `json:"data"`
}

func (ate ACLTestError) Error() string {
	var sb strings.Builder

	fmt.Fprintln(&sb, ate.Message)
	fmt.Fprintln(&sb)

	for _, data := range ate.Data {
		fmt.Fprintf(&sb, "For user %s:\n", data.User)
		for _, err := range data.Errors {
			fmt.Fprintf(&sb, "- %s\n", err)
		}
	}

	return sb.String()
}

type ACLTestErrorDetail struct {
	User   string   `json:"user"`
	Errors []string `json:"errors"`
}

func getACLETag(ctx context.Context, tailnet, apiKey string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://api.tailscale.com/api/v2/tailnet/%s/acl", tailnet), nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(apiKey, "")
	req.Header.Set("Accept", "application/hujson")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("wanted HTTP status code %d but got %d", http.StatusOK, resp.StatusCode)
	}

	return resp.Header.Get("ETag"), nil
}
