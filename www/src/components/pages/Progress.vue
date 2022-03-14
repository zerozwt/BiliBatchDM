<template>
<div v-if="has_task">
<n-space vertical>
    <n-card title="总览">
        <n-space vertical>
        <div>进度: 已向{{num_done}}人发送 / 总共{{num_total}}人</div>
        <div>
            <n-progress v-if="all_done" :percentage="100" :indicator-placement="'inside'" status="success"/>
            <n-progress v-else type="line" :percentage="percentage" :indicator-placement="'inside'" processing/>
        </div>
        <div>群发任务开始时间: {{start_time}}</div>
        消息内容:
        <div>
            <div v-if="msg_type == '1'" class="msg-content" id="msg-content">
                <div v-for="item in content_lines" :key="item.key">{{item.text}}</div>
            </div>
            <div v-else class="msg-content" id="msg-content">
                <img :src="content" style="max-width: 100%; width: auto;"/>
            </div>
        </div>
        </n-space>
    </n-card>
    <n-card title="详情">
        <n-list id="result_list" style="margin: 0">
            <n-list-item v-for="item in result_list" :key="item.uid">
                <n-thing :title="'接收者UID: ' + item.uid.toString()">{{render_result(item)}}</n-thing>
            </n-list-item>
        </n-list>
    </n-card>
</n-space>
</div>
<div v-else>
    暂无群发消息任务，请使用“发送私信”设定私信群发任务。
</div>
</template>

<script>
import { useMessage } from "naive-ui";
import { computed, onMounted, onUnmounted, ref } from "vue";
import axios from 'axios';
const dayjs = require('dayjs');
const render_ts = ts => dayjs.unix(ts).format("YYYY-MM-DD HH:mm:ss [UTC]Z");

export default {
    props: [],
    setup(props) {
        const message = useMessage();
        let has_task = ref(false);

        let start_ts = ref(0);
        let num_total = ref(0);
        let num_done = ref(0);
        let msg_type = ref("1");
        let content = ref("");
        let result_list = ref([]);

        let percentage = computed(() => num_total.value == 0 ? 0 : Math.floor(num_done.value * 100.0 / num_total.value));
        let all_done = computed(() => num_done.value == num_total.value);
        let start_time = computed(() => render_ts(start_ts.value));

        let line_key = 0;
        let content_lines = computed(() => content.value.trim().split("\n").map(x => {return {text: x, key: ++line_key};}));

        let render_result = result => {
            if (!result.done) {
                return "尚未发送";
            }
            if (result.err.length == 0) {
                return render_ts(result.ts) + " 发送成功";
            }
            return render_ts(result.ts) + " 发送失败: " + result.err;
        };

        let query = () => {
            console.log("query", Math.random());
            axios.get("/api/progress").then(rsp => {
                if (rsp.data.code == 233) {
                    has_task.value = false;
                    return;
                }
                if (rsp.data.code != 0) {
                    message.error(`[${rsp.data.code}]请求失败: ${rsp.data.msg}`);
                    return;
                }
                has_task.value = true;

                let data = rsp.data.data;
                start_ts.value = data.start_ts;
                num_total.value = data.num_total;
                num_done.value = data.num_done;
                msg_type.value = data.msg_type;
                content.value = data.content;
                result_list.value = data.result_list;
            }).catch(err => message.error(JSON.stringify(err)));
        };

        let timeoutID = setInterval(query, 1000);
        onMounted(query);
        onUnmounted(() => clearTimeout(timeoutID));

        return {
            has_task,

            start_time,
            num_total,
            num_done,
            msg_type,
            content,
            content_lines,
            result_list,
            percentage,
            all_done,

            render_result,
        };
    },
};
</script>

<style>
.send-succ {
    color: #18a058;
}
.send-fail {
    color: #d03050;
}
.msg-content {
    background-color: rgba(128, 128, 128, 0.1);
    padding: 16px;
}
</style>