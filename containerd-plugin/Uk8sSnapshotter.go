package main

import (
	"context"

	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/snapshots"
	"github.com/containerd/containerd/snapshots/native"
)

type Uk8sSnapshotter struct {
	native snapshots.Snapshotter
}

func NewUk8sSnapshotter(root, domain string, udisk any) (snapshots.Snapshotter, error) {
	sn, err := native.NewSnapshotter(root)
	if err != nil {
		return nil, err
	}
	return &Uk8sSnapshotter{
		native: sn,
	}, nil
}

func (s *Uk8sSnapshotter) Stat(ctx context.Context, key string) (snapshots.Info, error) {
	return s.native.Stat(ctx, key)
}

func (s *Uk8sSnapshotter) Update(ctx context.Context, Info snapshots.Info, fieldpaths ...string) (snapshots.Info, error) {
	return s.native.Update(ctx, Info, fieldpaths...)
}

func (s *Uk8sSnapshotter) Usage(ctx context.Context, key string) (snapshots.Usage, error) {
	return s.native.Usage(ctx, key)
}

func (s *Uk8sSnapshotter) Mounts(ctx context.Context, key string) ([]mount.Mount, error) {
	return s.native.Mounts(ctx, key)
}

func (s *Uk8sSnapshotter) Prepare(ctx context.Context, key, parent string, opts ...snapshots.Opt) ([]mount.Mount, error) {
	return s.native.Prepare(ctx, key, parent)
}

func (s *Uk8sSnapshotter) View(ctx context.Context, key, parent string, opts ...snapshots.Opt) ([]mount.Mount, error) {
	return s.native.View(ctx, key, parent)
}

func (s *Uk8sSnapshotter) Commit(ctx context.Context, name, key string, opts ...snapshots.Opt) error {
	return s.native.Commit(ctx, name, key, opts...)
}

func (s *Uk8sSnapshotter) Remove(ctx context.Context, key string) error {
	return s.native.Remove(ctx, key)
}

func (s *Uk8sSnapshotter) Walk(ctx context.Context, fn snapshots.WalkFunc, filters ...string) error {
	return s.native.Walk(ctx, fn, filters...)
}

func (s *Uk8sSnapshotter) Close() error {
	return s.native.Close()
}
