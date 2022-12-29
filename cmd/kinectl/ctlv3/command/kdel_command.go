// Copyright 2015 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/pkg/v3/cobrautl"
)

// NewKDelCommand returns the cobra command for "kdel".
func NewKDelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kdel [options] <key>",
		Short: "Remove the specified key",
		Run:   kdelCommandFunc,
	}
	return cmd
}

// kdelCommandFunc executes the "kdel" command.
func kdelCommandFunc(cmd *cobra.Command, args []string) {
	key, _ := getKDelOp(args)
	ctx, cancel := commandCtx(cmd)
	kine := mustKineClientCfgFromCmd(cmd)
	val, err := kine.Get(ctx, key)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}
	if err := kine.Delete(ctx, key, val.Modified); err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}
	cancel()
	resp := &clientv3.DeleteResponse{
		Header: &etcdserverpb.ResponseHeader{
			Revision: val.Modified,
		},
		Deleted: 1,
		PrevKvs: []*mvccpb.KeyValue{
			{
				Key:            val.Key,
				Value:          val.Data,
				Version:        0,
				ModRevision:    val.Modified,
				CreateRevision: val.Modified,
			},
		},
	}
	display.Del(*resp)
}

func getKDelOp(args []string) (string, []clientv3.OpOption) {
	if len(args) != 1 {
		cobrautl.ExitWithError(cobrautl.ExitBadArgs, fmt.Errorf("kdel command needs one argument as key"))
	}

	opts := []clientv3.OpOption{}
	key := args[0]

	return key, opts
}
