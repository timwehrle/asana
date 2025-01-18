package login

import (
	"bytes"
	"github.com/google/shlex"
	"github.com/stretchr/testify/require"
	"github.com/timwehrle/asana/pkg/factory"
	"github.com/timwehrle/asana/pkg/iostreams"
	"testing"
)

func TestNewCmdLogin(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		stdin    string
		stdinTTY bool
		wants    LoginOptions
		wantsErr bool
	}{
		{
			name:  "with-token and without workspace",
			cli:   "--with-token",
			stdin: "test-token\n",
			wants: LoginOptions{
				Token: "test-token",
			},
			wantsErr: true,
		},
		{
			name:  "with-token and with workspace",
			cli:   "--with-token --workspace \"Test Workspace\"",
			stdin: "test-token\n",
			wants: LoginOptions{
				Token:     "test-token",
				Workspace: "Test Workspace",
			},
		},
		{
			name: "with workspace and without token",
			cli:  "--workspace \"Test Workspace\"",
			wants: LoginOptions{
				Interactive: true,
				Workspace:   "Test Workspace",
			},
		},
		{
			name: "interactive login run",
			cli:  "",
			wants: LoginOptions{
				Interactive: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, stdin, _, _ := iostreams.Test()
			f := &factory.Factory{
				IOStreams: ios,
			}

			ios.IsStdoutTTY = true
			ios.IsStdinTTY = tt.stdinTTY
			if tt.stdin != "" {
				stdin.WriteString(tt.stdin)
			}

			argv, err := shlex.Split(tt.cli)
			require.NoError(t, err)

			var gotOpts *LoginOptions
			cmd := NewCmdLogin(*f, func(opts *LoginOptions) error {
				gotOpts = opts
				return nil
			})

			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(&bytes.Buffer{})
			cmd.SetErr(&bytes.Buffer{})

			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			require.Equal(t, tt.wants.Token, gotOpts.Token)
			require.Equal(t, tt.wants.Workspace, gotOpts.Workspace)
			require.Equal(t, tt.wants.Interactive, gotOpts.Interactive)
		})
	}
}
