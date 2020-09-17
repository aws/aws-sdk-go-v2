/*
Package processcreds is a credential Provider to retrieve `credential_process`
credentials.

WARNING: The following describes a method of sourcing credentials from an external
process. This can potentially be dangerous, so proceed with caution. Other
credential providers should be preferred if at all possible. If using this
option, you should make sure that the config file is as locked down as possible
using security best practices for your operating system.

Concurrency and caching

The Provider is not safe to be used concurrently, and does not provide any
caching of credentials retrieved. You should wrap the Provider with a
`aws.CredentialsCache` to provide concurrency safety, and caching of
credentials.

Loading credentials with the SDKs AWS Config

You can use credentials from a AWS shared config `credential_process` in a
variety of ways.

One way is to setup your shared config file, located in the default
location, with the `credential_process` key and the command you want to be
called. You also need to set the AWS_SDK_LOAD_CONFIG environment variable
(e.g., `export AWS_SDK_LOAD_CONFIG=1`) to use the shared config file.

    [default]
    credential_process = /command/to/call

Loading configuration using external will use the credential process to
retrieve credentials. NOTE: If there are credentials in the profile you are
using, the credential process will not be used.

    // Initialize a session to load credentials.
	cfg, _ := config.LoadDefaultConfig()

    // Create S3 service client to use the credentials.
    svc := s3.NewFromConfig(cfg)

Loading credentials with the Provider directly

Another way to use the `credential_process` method is by using `NewProvider`
function and providing a command to be executed to retrieve credentials:

    // Create credentials using the Provider.
	provider := processcreds.NewProvider("/path/to/command")

    // Create the service client value configured for credentials.
    svc := s3.New(s3.Options{
		Credentials: provider,
	})

If you need more control, you can set any configurable options in the
credentials using one or more option functions.

    provider := processcreds.NewProvider("/path/to/command",
        func(o *processcreds.Options) {
			// Override the provider's default timeout
            o.Timeout = 2 * time.Minute
        })

You can also use your own `exec.Cmd` value by satisfying a value that satisfies
the `NewCommandBuilder` interface and use the `NewProviderCommand` constructor.

	// Create an exec.Cmd
	cmdBuilder := processcreds.NewCommandBuilderFunc(
		func(ctx context.Context) (*exec.Cmd, error) {
			cmd := exec.CommandContext(ctx,
				"customCLICommand",
				"-a", "argument",
			)
			cmd.Env = []string{
				"ENV_VAR_FOO=value",
				"ENV_VAR_BAR=other_value",
			}

			return cmd, nil
		},
	)

	// Create credentials using your exec.Cmd and custom timeout
	provider := processcreds.NewProviderCommand(cmdBuilder,
		func(opt *processcreds.Provider) {
			// optionally override the provider's default timeout
			opt.Timeout = 1 * time.Second
		})
*/
package processcreds

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/sdkio"
)

const (
	// ProviderName is the name this credentials provider will label any
	// returned credentials Value with.
	ProviderName = `ProcessProvider`

	// DefaultTimeout default limit on time a process can run.
	DefaultTimeout = time.Duration(1) * time.Minute
)

// ProviderError is an error indicating failure initializing or executing the
// process credentials provider
type ProviderError struct {
	Err error
}

// Error returns the error message.
func (e *ProviderError) Error() string {
	return fmt.Sprintf("process provider error: %v", e.Err)
}

// Unwrap returns the underlying error the provider error wraps.
func (e *ProviderError) Unwrap() error {
	return e.Err
}

// Provider satisfies the credentials.Provider interface, and is a
// client to retrieve credentials from a process.
type Provider struct {
	// Provides a constructor for exec.Cmd that are invoked by the provider for
	// retrieving credentials. Use this to provide custom creation of exec.Cmd
	// with things like environment variables, or other configuration.
	//
	// The provider defaults to the DefaultNewCommand function.
	commandBuilder NewCommandBuilder

	options Options
}

// Options is the configuration options for configuring the Provider.
type Options struct {
	// ExpiryWindow will allow the credentials to trigger refreshing prior to
	// the credentials actually expiring. This is beneficial so race conditions
	// with expiring credentials do not cause request to fail unexpectedly
	// due to ExpiredTokenException exceptions.
	//
	// For example, an ExpiryWindow of 10s would cause calls to the
	// Credentials.IsExpired() method to return true 10 seconds before the
	// credentials would of actually expired.
	//
	// If ExpiryWindow is 0 or less, it will be ignored.
	ExpiryWindow time.Duration

	// Timeout limits the time a process can run.
	Timeout time.Duration
}

// NewCommandBuilder provides the interface for specifying how command will be
// created that the Provider will use to retrieve credentials with.
type NewCommandBuilder interface {
	NewCommand(context.Context) (*exec.Cmd, error)
}

// NewCommandBuilderFunc provides a wrapper type around a function pointer to
// satisfy the NewCommandBuilder interface.
type NewCommandBuilderFunc func(context.Context) (*exec.Cmd, error)

// NewCommand calls the underlying function pointer the builder was initialized with.
func (fn NewCommandBuilderFunc) NewCommand(ctx context.Context) (*exec.Cmd, error) {
	return fn(ctx)
}

// DefaultNewCommandBuilder provides the default NewCommandBuilder
// implementation used by the provider. It takes a command and arguments to
// invoke. The command will also be initialized with the current process
// environment variables, stderr, and stdin pipes.
type DefaultNewCommandBuilder struct {
	Args []string
}

// NewCommand returns an initialized exec.Cmd with the builder's initialized
// Args. The command is also initialized current process environment variables,
// stderr, and stdin pipes.
func (b DefaultNewCommandBuilder) NewCommand(ctx context.Context) (*exec.Cmd, error) {
	var cmdArgs []string
	if runtime.GOOS == "windows" {
		cmdArgs = []string{"cmd.exe", "/C"}
	} else {
		cmdArgs = []string{"sh", "-c"}
	}

	if len(b.Args) == 0 {
		return nil, &ProviderError{
			Err: fmt.Errorf("failed to prepare command: command must not be empty"),
		}
	}

	cmdArgs = append(cmdArgs, b.Args...)
	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = os.Environ()

	cmd.Stderr = os.Stderr // display stderr on console for MFA
	cmd.Stdin = os.Stdin   // enable stdin for MFA

	return cmd, nil
}

// NewProvider returns a pointer to a new Credentials object wrapping the
// Provider.
//
// The provider defaults to the DefaultNewCommandBuilder for creating command
// the Provider will use to retrieve credentials with.
func NewProvider(command string, options ...func(*Options)) *Provider {
	var args []string

	// Ensure that the command arguments are not set if the provided command is
	// empty. This will error out when the command is executed since no
	// arguments are specified.
	if len(command) > 0 {
		args = []string{command}
	}

	commanBuilder := DefaultNewCommandBuilder{
		Args: args,
	}
	return NewProviderCommand(commanBuilder, options...)
}

// NewProviderCommand returns a pointer to a new Credentials object with the
// specified command, and default timeout duration. Use this to provide custom
// creation of exec.Cmd for options like environment variables, or other
// configuration.
func NewProviderCommand(builder NewCommandBuilder, options ...func(*Options)) *Provider {
	p := &Provider{
		commandBuilder: builder,
		options: Options{
			Timeout: DefaultTimeout,
		},
	}

	for _, option := range options {
		option(&p.options)
	}

	return p
}

type credentialProcessResponse struct {
	Version         int
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string
	SessionToken    string
	Expiration      *time.Time
}

// Retrieve executes the credential process command and returns the
// credentials, or error if the command fails.
func (p *Provider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	out, err := p.executeCredentialProcess(ctx)
	if err != nil {
		return aws.Credentials{Source: ProviderName}, err
	}

	// Serialize and validate response
	resp := &credentialProcessResponse{}
	if err = json.Unmarshal(out, resp); err != nil {
		return aws.Credentials{Source: ProviderName}, &ProviderError{
			Err: fmt.Errorf("parse failed of process output: %s, error: %w", out, err),
		}
	}

	if resp.Version != 1 {
		return aws.Credentials{Source: ProviderName}, &ProviderError{
			Err: fmt.Errorf("wrong version in process output (not 1)"),
		}
	}

	if len(resp.AccessKeyID) == 0 {
		return aws.Credentials{Source: ProviderName}, &ProviderError{
			Err: fmt.Errorf("missing AccessKeyId in process output"),
		}
	}

	if len(resp.SecretAccessKey) == 0 {
		return aws.Credentials{Source: ProviderName}, &ProviderError{
			Err: fmt.Errorf("missing SecretAccessKey in process output"),
		}
	}

	creds := aws.Credentials{
		Source:          ProviderName,
		AccessKeyID:     resp.AccessKeyID,
		SecretAccessKey: resp.SecretAccessKey,
		SessionToken:    resp.SessionToken,
	}

	// Handle expiration
	if resp.Expiration != nil {
		creds.CanExpire = true
		creds.Expires = (*resp.Expiration).Add(-p.options.ExpiryWindow)
	}

	return creds, nil
}

// executeCredentialProcess starts the credential process on the OS and
// returns the results or an error.
func (p *Provider) executeCredentialProcess(ctx context.Context) ([]byte, error) {
	if p.options.Timeout >= 0 {
		var cancelFunc func()
		ctx, cancelFunc = context.WithTimeout(ctx, p.options.Timeout)
		defer cancelFunc()
	}

	cmd, err := p.commandBuilder.NewCommand(ctx)
	if err != nil {
		return nil, err
	}

	// get creds json on process's stdout
	output := bytes.NewBuffer(make([]byte, 0, int(8*sdkio.KibiByte)))
	if cmd.Stdout != nil {
		cmd.Stdout = io.MultiWriter(cmd.Stdout, output)
	} else {
		cmd.Stdout = output
	}

	execCh := make(chan error, 1)
	go executeCommand(cmd, execCh)

	select {
	case execError := <-execCh:
		if execError == nil {
			break
		}
		select {
		case <-ctx.Done():
			return output.Bytes(), &ProviderError{
				Err: fmt.Errorf("credential process timed out: %w", execError),
			}
		default:
			return output.Bytes(), &ProviderError{
				Err: fmt.Errorf("error in credential_process: %w", execError),
			}
		}
	}

	out := output.Bytes()
	if runtime.GOOS == "windows" {
		// windows adds slashes to quotes
		out = bytes.ReplaceAll(out, []byte(`\"`), []byte(`"`))
	}

	return out, nil
}

func executeCommand(cmd *exec.Cmd, exec chan error) {
	// Start the command
	err := cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	exec <- err
}
