const BatchSendSucc = {code: 0};
const BatchSendBusy = {code: 1};
const BatchSendFail = {code: 233, msg: "Some error"};

export default {
    'post|^/api/batch_send$': opt => BatchSendSucc,
};