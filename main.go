package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"

	"github.com/go-git/go-git/plumbing/transport"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	gitsRootDir string
)

// toDirPath example: arg -> github.com/nu50218/gits
func toDirPath(arg string) (string, error) {
	if len(arg) == 0 {
		return "", errors.New("arg is empty")
	}

	// unquote if quoted
	if unquoted, err := strconv.Unquote(arg); err == nil {
		arg = unquoted
	}

	// check if url
	if u, err := url.Parse(arg); err == nil {
		arg = path.Join(u.Hostname(), u.Path)
	}

	// trim .git
	arg = strings.TrimSuffix(arg, ".git")

	return path.Clean(arg), nil
}

func askAuthMethod() (transport.AuthMethod, error) {
	const (
		authMethodBasic       = "basic_auth"
		authMethodAccessToken = "access_token"
	)

	prompt := promptui.Select{
		Label: "authentication method",
		Items: interface{}([]string{
			authMethodBasic,
			authMethodAccessToken,
		}),
		HideHelp: true,
	}

	_, res, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to ask authentication method: %w", err)
	}

	switch res {
	case authMethodBasic:
		fmt.Print("username: ")
		var username string
		if _, err := fmt.Scan(&username); err != nil {
			return nil, fmt.Errorf("failed to read username: %w", err)
		}
		fmt.Print("password: ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr, "")
		if err != nil {
			return nil, fmt.Errorf("failed to read password from stdin: %v", err)
		}
		return &http.BasicAuth{
			Username: username,
			Password: string(password),
		}, nil

	case authMethodAccessToken:
		fmt.Print("token: ")
		token, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Fprintln(os.Stderr, "")
		if err != nil {
			return nil, fmt.Errorf("failed to read token from stdin: %v", err)
		}
		return &http.TokenAuth{
			Token: string(token),
		}, nil

	default:
		return nil, errors.New("unimplemented authentication method")
	}
}

func cloneAction(ctx *cli.Context) error {
	arg := ctx.Args().First()

	dirPath, err := toDirPath(arg)
	if err != nil {
		return NewError(ErrorTypeInvalidArgument, "failed to parse argument: %v", err)
	}

	cloneURL := "https://" + dirPath + ".git"
	cloneDirPath := path.Join(gitsRootDir, dirPath)

	fmt.Fprintf(os.Stderr, `Cloning '%s' into '%s'`, cloneURL, cloneDirPath)
	fmt.Fprintln(os.Stderr, "")

	_, err = git.PlainClone(cloneDirPath, false, &git.CloneOptions{
		URL:      cloneURL,
		Progress: os.Stderr,
	})

	// if error, try to clone with authentication
	// NOTE: err != transport.ErrAuthenticationRequired ???????????
	if err != nil && err.Error() == transport.ErrAuthenticationRequired.Error() {
		var authMethod transport.AuthMethod
		authMethod, err = askAuthMethod()
		if err != nil {
			return NewError(ErrorTypeGeneral, "failed to clone with authentication: %v", err)
		}
		_, err = git.PlainClone(cloneDirPath, false, &git.CloneOptions{
			URL:      "https://" + dirPath + ".git",
			Progress: os.Stderr,
			Auth:     authMethod,
		})
	}

	if err != nil {
		return NewError(ErrorTypeGeneral, "failed to clone: %v", err)
	}

	return nil
}

func main() {
	// get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get home dir: %v", err)
	}
	gitsRootDir = path.Join(homeDir, "gits")

	// check '$ git' is available
	if err := exec.Command("git", "help").Run(); errors.Is(err, exec.ErrNotFound) {
		fmt.Fprintln(os.Stderr, "you need to install git to use gits")
		os.Exit(1)
	}

	app := &cli.App{
		Name:      "gits",
		ArgsUsage: "gits [path]",
		Action:    cloneAction,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		var myErr *Error
		if !errors.As(err, &myErr) {
			os.Exit(1)
		}

		switch myErr.Type {
		case ErrorTypeGeneral:
			os.Exit(1)
		case ErrorTypeInvalidArgument:
			os.Exit(128)
		}
	}
}
