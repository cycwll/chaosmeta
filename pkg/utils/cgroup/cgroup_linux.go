/*
 * Copyright 2022-2023 Chaos Meta Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cgroup

import (
	"context"
	"fmt"
	"github.com/ChaosMetaverse/chaosmetad/pkg/crclient"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils"
	"github.com/ChaosMetaverse/chaosmetad/pkg/utils/cmdexec"
	"strconv"
	"strings"
)

func NewCgroup(ctx context.Context, cgroupPath string, configCmdStr string) error {
	if err := cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("mkdir %s%s%s", cgroupPath, utils.CmdSplit, configCmdStr)); err != nil {
		return err
	}

	return nil
}

func GetContainerCgroupPath(ctx context.Context, cr, containerID, subSys string) (string, error) {
	client, err := crclient.GetClient(ctx, cr)
	if err != nil {
		return "", fmt.Errorf("get %s client error: %s", cr, err.Error())
	}

	pid, err := client.GetPidById(context.Background(), containerID)
	if err != nil {
		return "", fmt.Errorf("get pid of container[%s] error: %s", containerID, err.Error())
	}

	cPath, err := GetpidCurCgroup(ctx, pid, subSys)
	if err != nil {
		return "", fmt.Errorf("get cgroup[%s] path of process[%d] error: %s", subSys, pid, err.Error())
	}

	return cPath, nil
}

func GetBlkioCPath(uid string) string {
	return fmt.Sprintf("%s/%s/%s_%s", utils.RootCgroupPath, BLKIO, BlkioCgroupName, uid)
}

func CheckPidListBlkioCgroup(ctx context.Context, pidList []int) error {
	for _, unitP := range pidList {
		oldPath, err := GetpidCurCgroup(ctx, unitP, BLKIO)
		if err != nil {
			return fmt.Errorf("get old cgroup path of process[%d] error: %s", unitP, err.Error())
		}

		if strings.Index(oldPath, BlkioCgroupName) >= 0 {
			return fmt.Errorf("%d is in experiment[%s]", unitP, oldPath)
		}
	}

	return nil
}

func GetPidListCurCgroup(ctx context.Context, pidList []int, subSys string) (map[int]string, error) {
	var re = make(map[int]string)
	for _, unitP := range pidList {
		oldPath, err := GetpidCurCgroup(ctx, unitP, subSys)
		if err != nil {
			return nil, fmt.Errorf("get old cgroup path of process[%d] error: %s", unitP, err.Error())
		}
		re[unitP] = oldPath
	}

	return re, nil
}

func GetpidCurCgroup(ctx context.Context, pid int, subSys string) (string, error) {
	reByte, err := cmdexec.RunBashCmdWithOutput(ctx, fmt.Sprintf("cat /proc/%d/cgroup | grep -w %s", pid, subSys))
	if err != nil {
		return "", fmt.Errorf("run cmd error: %s", err.Error())
	}

	out := strings.TrimSpace(string(reByte))
	sArr := strings.Split(out, ":")
	if strings.Index(sArr[1], subSys) < 0 {
		return "", fmt.Errorf("out string is not valid: %s", out)
	}

	return sArr[2], nil
}

func MovePidListToCgroup(ctx context.Context, pidList []int, cgroupPath string) error {
	for _, unit := range pidList {
		if err := MoveTaskToCgroup(ctx, unit, cgroupPath); err != nil {
			return fmt.Errorf("move pid[%d] to cgroup[%s] error: %s", unit, cgroupPath, err.Error())
		}
	}

	return nil
}

func MoveTaskToCgroup(ctx context.Context, pid int, cgroupPath string) error {
	if err := cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("echo %d > %s/tasks", pid, cgroupPath)); err != nil {
		return err
	}

	return nil
}

//func MoveProcToCgroup(pid int, cgroupPath string) error {
//	if err := cmdexec.RunBashCmdWithoutOutput(fmt.Sprintf("echo %d > %s/cgroup.procs", pid, cgroupPath)); err != nil {
//		return err
//	}
//
//	return nil
//}

func GetPidStrListByCgroup(ctx context.Context, cgroupPath string) ([]int, error) {
	reByte, err := cmdexec.RunBashCmdWithOutput(ctx, fmt.Sprintf("cat %s/tasks", cgroupPath))
	if err != nil {
		return nil, fmt.Errorf("run cmd error: %s", err.Error())
	}

	var pidList []int
	strList := strings.Split(string(reByte), "\n")
	for _, unit := range strList {
		if unit == "" {
			continue
		}

		pid, err := strconv.Atoi(unit)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid pid: %s", unit, err.Error())
		}

		pidList = append(pidList, pid)
	}

	return pidList, nil
}

func RemoveCgroup(ctx context.Context, cgroupPath string) error {
	if err := cmdexec.RunBashCmdWithoutOutput(ctx, fmt.Sprintf("rmdir %s", cgroupPath)); err != nil {
		return fmt.Errorf("cmd exec error: %s", err.Error())
	}

	return nil
}