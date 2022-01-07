var Mock = require('mockjs');

Mock.setup({timeout: "200-400"});

let confArr = [];

const files = require.context('.', true, /\.js$/);
files.keys().forEach((key) => {
    if (key == "./index.js") return;
    confArr = confArr.concat(files(key).default);
});

confArr.forEach((item) => {
    for (let [mathod_path, target] of Object.entries(item)) {
        let tmp = mathod_path.split('|');
        let method = tmp[0];
        let path = new RegExp(tmp[1]);
        console.log("MOCK", method, path);
        Mock.mock(path, method, target);
    }
});