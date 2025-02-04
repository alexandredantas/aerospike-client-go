// Copyright 2014-2022 Aerospike, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aerospike

import (
	"math/rand"

	kvs "github.com/aerospike/aerospike-client-go/v6/proto/kvs"
)

type executeCommand struct {
	readCommand

	// overwrite
	policy       *WritePolicy
	packageName  string
	functionName string
	args         *ValueArray
}

func newExecuteCommand(
	cluster *Cluster,
	policy *WritePolicy,
	key *Key,
	packageName string,
	functionName string,
	args *ValueArray,
) (executeCommand, Error) {
	var err Error
	var partition *Partition
	if cluster != nil {
		partition, err = PartitionForWrite(cluster, &policy.BasePolicy, key)
		if err != nil {
			return executeCommand{}, err
		}
	}

	readCommand, err := newReadCommand(cluster, &policy.BasePolicy, key, nil, partition)
	if err != nil {
		return executeCommand{}, err
	}

	return executeCommand{
		readCommand:  readCommand,
		policy:       policy,
		packageName:  packageName,
		functionName: functionName,
		args:         args,
	}, nil
}

func (cmd *executeCommand) writeBuffer(ifc command) Error {
	return cmd.setUdf(cmd.policy, cmd.key, cmd.packageName, cmd.functionName, cmd.args)
}

func (cmd *executeCommand) getNode(ifc command) (*Node, Error) {
	return cmd.partition.GetNodeWrite(cmd.cluster)
}

func (cmd *executeCommand) prepareRetry(ifc command, isTimeout bool) bool {
	cmd.partition.PrepareRetryWrite(isTimeout)
	return true
}

func (cmd *executeCommand) isRead() bool {
	return false
}

func (cmd *executeCommand) Execute() Error {
	return cmd.execute(cmd)
}

func (cmd *executeCommand) ExecuteGRPC(clnt *ProxyClient) Error {
	cmd.dataBuffer = bufPool.Get().([]byte)
	defer cmd.grpcPutBufferBack()

	err := cmd.prepareBuffer(cmd, cmd.policy.deadline())
	if err != nil {
		return err
	}

	req := kvs.AerospikeRequestPayload{
		Id:          rand.Uint32(),
		Iteration:   1,
		Payload:     cmd.dataBuffer[:cmd.dataOffset],
		WritePolicy: cmd.policy.grpc(),
	}

	conn, err := clnt.grpcConn()
	if err != nil {
		return err
	}

	client := kvs.NewKVSClient(conn)

	ctx := cmd.policy.grpcDeadlineContext()

	res, gerr := client.Execute(ctx, &req)
	if gerr != nil {
		return newGrpcError(gerr, gerr.Error())
	}

	cmd.commandWasSent = true

	defer clnt.returnGrpcConnToPool(conn)

	if res.Status != 0 {
		return newGrpcStatusError(res)
	}

	cmd.conn = newGrpcFakeConnection(res.Payload, nil)
	err = cmd.parseResult(cmd, cmd.conn)
	if err != nil {
		return err
	}

	return nil
}
