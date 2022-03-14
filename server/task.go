package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"sync"
	"time"

	dm "github.com/zerozwt/BLiveDanmaku"
)

type TaskRequest struct {
	Sender   uint   `json:"sender_uid"`
	Sess     string `json:"sess"`
	JCT      string `json:"jct"`
	MinSec   uint   `json:"min_time"`
	MaxSec   uint   `json:"max_time"`
	MsgType  string `json:"msg_type"`
	Content  string `json:"content"`
	RecvList []uint `json:"recv_list"`
}

type TaskProgressInfo struct {
	StartTs    int64        `json:"start_ts"`
	NumTotal   int          `json:"num_total"`
	NumDone    int          `json:"num_done"`
	MsgType    string       `json:"msg_type"`
	Content    string       `json:"content"`
	ResultList []TaskResult `json:"result_list"`
}

type TaskResult struct {
	UID  int64  `json:"uid"`
	Done bool   `json:"done"`
	Err  string `json:"err"`
	Ts   int64  `json:"ts"`
}

var TaskLock sync.Mutex
var CurrentTask *TaskProgressInfo

func SubmitTask(req *TaskRequest) error {
	TaskLock.Lock()
	defer TaskLock.Unlock()

	if CurrentTask != nil && !CurrentTask.AllDone() {
		return Err233
	}

	CurrentTask = &TaskProgressInfo{
		StartTs:    time.Now().Unix(),
		NumTotal:   len(req.RecvList),
		NumDone:    0,
		MsgType:    req.MsgType,
		Content:    req.Content,
		ResultList: make([]TaskResult, 0, len(req.RecvList)),
	}
	for _, item := range req.RecvList {
		CurrentTask.ResultList = append(CurrentTask.ResultList, TaskResult{
			UID:  int64(item),
			Done: false,
			Err:  "",
			Ts:   0,
		})
	}
	go CurrentTask.doRequest(req)

	return nil
}

func (info *TaskProgressInfo) AllDone() bool {
	return info.NumDone == info.NumTotal
}

func (info *TaskProgressInfo) Clone() *TaskProgressInfo {
	TaskLock.Lock()
	defer TaskLock.Unlock()
	if info == nil {
		return nil
	}

	ret := &TaskProgressInfo{
		StartTs:    info.StartTs,
		NumTotal:   info.NumTotal,
		NumDone:    info.NumDone,
		MsgType:    info.MsgType,
		Content:    info.Content,
		ResultList: make([]TaskResult, 0, len(info.ResultList)),
	}
	for _, item := range info.ResultList {
		ret.ResultList = append(ret.ResultList, TaskResult{
			UID:  item.UID,
			Done: item.Done,
			Err:  item.Err,
			Ts:   item.Ts,
		})
	}

	return ret
}

func (info *TaskProgressInfo) doRequest(req *TaskRequest) {
	logger().Printf("start sending direct message to %d recievers ...", len(req.RecvList))

	mark_fail := func(err error) {
		TaskLock.Lock()
		defer TaskLock.Unlock()
		err_str := err.Error()
		for idx, reciever := range req.RecvList {
			result := TaskResult{
				UID:  int64(reciever),
				Done: true,
				Ts:   time.Now().Unix(),
				Err:  err_str,
			}
			info.ResultList[idx] = result
		}
	}

	var pic *dm.UploadedPic
	ext := ""

	if req.MsgType == "2" {
		logger().Printf("uploading picture %s ...", req.Content)
		if len(req.Content) > 0 && req.Content[0] == '/' {
			req.Content = req.Content[1:]
		}
		pic_data, err := ioutil.ReadFile(filepath.Join(www_root, req.Content))
		if err != nil {
			mark_fail(err)
			return
		}
		pic, err = dm.UploadPic(pic_data, filepath.Base(req.Content), req.Sess, req.JCT)
		if err != nil {
			mark_fail(err)
			return
		}
		ext = filepath.Ext(req.Content)
		if len(ext) > 0 && ext[0] == '.' {
			ext = ext[1:]
		}
		logger().Printf("uploaded picture: %+v", *pic)
	}

	for idx, reciever := range req.RecvList {
		if idx > 0 {
			min_interval := time.Second * time.Duration(req.MinSec)
			max_interval := time.Second * time.Duration(req.MaxSec)
			if max_interval > min_interval {
				min_interval += time.Duration(rand.Int63n(int64(max_interval - min_interval)))
			}
			time.Sleep(min_interval)
		}
		logger().Printf("(%d/%d) sending direct message to %d ...", idx+1, len(req.RecvList), reciever)
		var (
			dm_rsp *dm.SendDirectMsgRsp
			err    error
		)
		if req.MsgType == "1" {
			dm_rsp, err = dm.SendDirectMsg(int64(req.Sender), int64(reciever), req.Content, gDMDevID, req.Sess, req.JCT)
		} else {
			dm_rsp, err = dm.SendDirectMsgPicture(int64(req.Sender), int64(reciever), pic, ext, gDMDevID, req.Sess, req.JCT)
		}
		err_str := ""

		if err != nil {
			logger().Printf("send direct message to %d failed: %v", reciever, err)
			err_str = err.Error()
		} else {
			if dm_rsp.Code != 0 {
				logger().Printf("send direct message to %d failed: [%d] %s", reciever, dm_rsp.Code, dm_rsp.Message)
				err_str = fmt.Sprintf("[%d] %s", dm_rsp.Code, dm_rsp.Message)
			} else {
				logger().Printf("send direct message to %d success", reciever)
			}
		}

		TaskLock.Lock()
		info.NumDone += 1
		result := TaskResult{
			UID:  int64(reciever),
			Done: true,
			Ts:   time.Now().Unix(),
			Err:  err_str,
		}
		info.ResultList[idx] = result
		TaskLock.Unlock()
	}
}
