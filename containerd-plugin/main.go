package main

import (
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"

	snapshotsapiv1 "github.com/containerd/containerd/api/services/snapshots/v1"
	"github.com/containerd/containerd/contrib/snapshotservice"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("invalid args: usage: %s <unix addr> <root>\n", os.Args[0])
		os.Exit(1)
	}

	rpc := grpc.NewServer()
	sn, err := NewUk8sSnapshotter(os.Args[2], "uhub.service.ucloud.cn", "/dev/vdb")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	service := snapshotservice.FromSnapshotter(sn)
	snapshotsapiv1.RegisterSnapshotsServer(rpc, service)

	l, err := net.Listen("unix", os.Args[1])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	if err := rpc.Serve(l); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}
