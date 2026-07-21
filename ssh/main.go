package main

import (
	"context"
	"fmt"

	"dagger/ssh/internal/dagger"
)

type Ssh struct {
	AlpineVersion string

	Key  *dagger.Secret
	Host string
	User string
}

func New(
	// +optional
	alpineVersion string,
	key *dagger.Secret,
	host, user string,
) *Ssh {
	return &Ssh{
		alpineVersion,
		key, host, user,
	}
}

const (
	mountedKeyPath = "/tmp/mounted-key"
	keyPath        = "/tmp/key"
)

func (m *Ssh) container() *dagger.Container {
	alpineOpts := dagger.AlpineOpts{
		Packages: []string{"openssh-client-default"},
	}

	if m.AlpineVersion != "" {
		alpineOpts.AlpineVersion = m.AlpineVersion
	}

	return dag.Alpine(alpineOpts).
		Container().
		WithMountedSecret(mountedKeyPath, m.Key).
		WithExec([]string{
			"sh", "-c",
			`tr -d "\r" < ` + mountedKeyPath + ` > ` + keyPath,
		}).
		WithExec([]string{
			"sh", "-c",
			`echo -e "\n" >> ` + keyPath,
		}).
		WithExec([]string{
			"sh", "-c",
			`chmod 600 ` + keyPath,
		})
}

func (m *Ssh) Run(
	ctx context.Context,
	command string,
) error {
	_, err := m.container().
		WithExec([]string{
			"sh", "-c",
			fmt.Sprintf(
				`ssh -o StrictHostKeyChecking=accept-new -i %s %s@%s "%s"`,
				keyPath, m.User, m.Host, command,
			),
		}).
		Sync(ctx)

	return err
}
