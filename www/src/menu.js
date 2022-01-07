import {h} from 'vue';
import Link from './components/router/Link.vue';

export default [
    {
        key: "index",
        label: () => h(Link, {href: "/"}, () => "发送私信"),
    },
    {
        key: "progress",
        label: () => h(Link, {href: "/progress"}, () => "查看进度"),
    },
    {
        key: "help",
        label: () => h(Link, {href: "/help"}, () => "使用说明"),
    },
];