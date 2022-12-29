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
	"os"

	"github.com/spf13/cobra"
	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/pkg/v3/cobrautl"
)

// NewKPutCommand returns the cobra command for "kput".
func NewKPutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kput [options] <key> <value> (<value> can also be given from stdin)",
		Short: "Puts the given key into the store",
		Long: `
Puts the given key into the store.

When <value> begins with '-', <value> is interpreted as a flag.
Insert '--' for workaround:

$ kput <key> -- <value>
$ kput -- <key> <value>

For example,
$ cat file | kput <key>
will store the content of the file to <key>.
`,
		Run: kputCommandFunc,
	}
	return cmd
}

// kputCommandFunc executes the "kput" command.
func kputCommandFunc(cmd *cobra.Command, args []string) {
	key, value, _ := getKPutOp(args)
	ctx, cancel := commandCtx(cmd)
	kine := mustKineClientCfgFromCmd(cmd)
	val, err := kine.Get(ctx, key)
	if err != nil {
		err = kine.Create(ctx, key, []byte(value))
	} else {
		err = kine.Update(ctx, key, val.Modified, []byte(value))
	}
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}
	val, err = kine.Get(ctx, key)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}
	cancel()
	resp := &clientv3.PutResponse{
		Header: &etcdserverpb.ResponseHeader{
			Revision: val.Modified,
		},
		PrevKv: &mvccpb.KeyValue{
			Key:            val.Key,
			Value:          val.Data,
			Version:        0,
			ModRevision:    val.Modified,
			CreateRevision: val.Modified,
		},
	}
	display.Put(*resp)
}

func getKPutOp(args []string) (string, string, []clientv3.OpOption) {
	if len(args) == 0 {
		cobrautl.ExitWithError(cobrautl.ExitBadArgs, fmt.Errorf("kput command needs 1 argument and input from stdin or 2 arguments"))
	}

	key := args[0]

	var value string
	var err error
	value, err = argOrStdin(args, os.Stdin, 1)
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitBadArgs, fmt.Errorf("kput command needs 1 argument and input from stdin or 2 arguments"))
	}

	opts := []clientv3.OpOption{}

	return key, value, opts
}
