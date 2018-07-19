// Copyright (c) 2012 VMware, Inc.

package gosigar

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
	"strings"
	"unsafe"

	"github.com/prashanthpai/sunrpc"
)

const (
	NFS_PROGRAM          = uint32(100003)
	NFS_PROGRAM_VERSION2 = uint32(2)
	NFS_PROGRAM_VERSION3 = uint32(3)
	NFS_PROGRAM_VERSION4 = uint32(4)
)

func bytePtrToString(ptr *int8) string {
	bytes := (*[10000]byte)(unsafe.Pointer(ptr))

	n := 0
	for bytes[n] != 0 {
		n++
	}

	return string(bytes[0:n])
}

func chop(buf []byte) []byte {
	return buf[0 : len(buf)-1]
}

// try to ping nfs filesystem to check whether active
func FsPing(self FileSystem) error {
	if !strings.HasPrefix(self.SysTypeName, "nfs") {
		return nil
	}

	comma := strings.Index(self.DevName, ":")
	if comma <= 0 {
		return nil
	}

	if RpcPingNfs(self.DevName[:comma], NFS_PROGRAM_VERSION2) != nil {
		if RpcPingNfs(self.DevName[:comma], NFS_PROGRAM_VERSION3) != nil {
			return RpcPingNfs(self.DevName[:comma], NFS_PROGRAM_VERSION4)
		}
	}

	return nil
}

func RpcPingNfs(host string, version uint32) error {
	port, err := sunrpc.PmapGetPort(host+":"+strconv.Itoa(111), NFS_PROGRAM, version, sunrpc.IPProtoTCP)
	if err != nil || port == 0 {
		return errors.New(fmt.Sprintf("error when getting nfs program port(%d) %v", port, err))
	}

	procedureID := sunrpc.ProcedureID{
		ProgramNumber:   NFS_PROGRAM,
		ProgramVersion:  version,
		ProcedureNumber: 0,
	}

	sunrpc.RegisterProcedure(sunrpc.Procedure{procedureID, "Nfs.ProcNull"}, true)
	conn, err := net.Dial("tcp", host+":"+strconv.Itoa(int(port)))
	if err != nil {
		return err
	}
	client := rpc.NewClientWithCodec(sunrpc.NewClientCodec(conn, nil))
	defer client.Close()
	return client.Call("Nfs.ProcNull", nil, nil)
}
