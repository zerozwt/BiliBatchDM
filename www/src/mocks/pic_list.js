const PicList = {
    code: 0,
    msg: "",
    data: {
        list: [
            {url: "/pics/rabbit.png"},
            {url: "/cookie.png"},
        ],
    },
};

export default {
    'get|^/api/pic_list$': opt => PicList,
};