const ProgressRsp = {
    code: 0,
    msg: "",
    data: {
        start_ts: 1641528000,
        num_total: 3,
        num_done: 2,
        content: "测试一下\n多行的消息",
        result_list: [
            {uid: 1001, done: true, err: "", ts: 1641528000},
            {uid: 1002, done: true, err: "Some error", ts: 1641528010},
            {uid: 1003, done: false, err: "", ts: 1641528020},
        ],
    },
};

export default {
    'get|^/api/progress$': opt => ProgressRsp,
};