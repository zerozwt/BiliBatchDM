export default {
    rules: [
        {
            rule: /^\/$/,
            title: "index",
            component: "Index",
            key: "index",
        },
        {
            rule: /^\/progress$/,
            title: "progress",
            component: "Progress",
            key: "progress",
        },
        {
            rule: /^\/help$/,
            title: "help",
            component: "Help",
            key: "help",
        },
    ],
    default: {
        title: "404 Not Found",
        component: "404",
        key: "",
    },
}