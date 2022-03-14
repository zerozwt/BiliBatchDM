<template>
<n-space vertical>
<n-card title="发送者">
    <template #header-extra> 填写方式请参照“使用说明” </template>
    <n-space vertical>
        B站UID:
        <n-input-number v-model:value="sender_uid" :show-button="false" placeholder="从B站个人主页URL获取" />
        SESS_DATA:
        <n-input v-model:value="sess" type="text" placeholder="SESS_DATA in cookie" />
        bili_jct:
        <n-input v-model:value="jct" type="text" placeholder="bili_jct in cookie" />
    </n-space>
</n-card>
<n-card title="消息">
    <n-space vertical>
        随机发送间隔:
        <n-space>
            <n-input-number v-model:value="min_time" :show-button="false" placeholder="最低间隔" ><template #suffix>秒</template></n-input-number>
            <div style="margin-top: 4px">～</div>
            <n-input-number v-model:value="max_time" :show-button="false" placeholder="最高间隔" ><template #suffix>秒</template></n-input-number>
        </n-space>
        消息类型:
        <n-radio-group v-model:value="msg_type" name="mt">
            <n-radio-button key="1" value="1">文字消息</n-radio-button>
            <n-radio-button key="2" value="2">图片消息</n-radio-button>
        </n-radio-group>
        <n-space vertical v-if="msg_type == '1'">
            消息内容:
            <n-input v-model:value="content" type="textarea" placeholder="在这里输入消息内容" :rows="5"/>
        </n-space>
        <n-space vertical v-else>
            <n-space>
                <div style="padding-top: 4px">选择图片:</div>
                <n-button type="default" @click="getPicList">刷新图片列表</n-button>
            </n-space>
            <n-radio-group v-model:value="pic_url" name="pic">
                <n-radio v-for="item in pic_list" :key="item.url" :value="item.url">
                    <n-tooltip trigger="hover">
                        <template #trigger>
                            <n-image :src="item.url" preview-disabled object-fit="scale-down"/>
                        </template>
                        {{item.url}}
                    </n-tooltip>
                </n-radio>
            </n-radio-group>
        </n-space>
    </n-space>
</n-card>
<n-card title="接收者">
    <n-space vertical>
        接收者UID:
        <n-input v-model:value="recv_list_txt" type="textarea" placeholder="在这里输入接收者的B站UID，一行一个，不能重复" :rows="10"/>
    </n-space>
</n-card>
<n-space justify="end"><n-button type="primary" style="margin-right: 24px;" @click="submit" ><div style="padding: 4px 8px;">发送</div></n-button></n-space>
</n-space>
<n-modal :mask-closable="false" display-directive="show" v-model:show="showModal">
    <n-card title="群发私信" style="max-width: 600px">
        <n-space vertical>
            是否将以下私信内容:
            <div class="msg-preview" id="msg-preview">
                <div v-if="msg_type == '1'"><div v-for="item in content_preview" :key="item.key">{{item.text}}</div></div>
                <div v-else><img :src="pic_url" style="max-width: 100%; width: auto;"/></div>
            </div>
            发送给以下B站UID:
            <ul id="reciever_list">
                <li v-for="item in recievers" :key="item">{{item}}</li>
            </ul>
            <n-space justify="end">
                <n-button type="default" :disabled="submitting" @click="showModal = false">取消</n-button>
                <n-button type="primary" :disabled="submitting" :loading="submitting" @click="doSubmit">发送</n-button>
            </n-space>
        </n-space>
    </n-card>
</n-modal>
</template>

<script>
import { useMessage } from "naive-ui";
import { computed, ref, inject, watchEffect, onMounted } from "vue";
import axios from 'axios';

const loadCache = (key, def_value) => {
    let ret = window.sessionStorage.getItem(key);
    return ret ? ret : def_value;
};

export default {
    props: [],
    setup(props) {
        const message = useMessage();
        let showModal = ref(false);
        let submitting = ref(false);
        let recievers = ref([]);

        let sender_uid = ref(parseInt(loadCache("sender_uid", "0"), 10));
        let sess = ref(loadCache("sess", ""));
        let jct = ref(loadCache("jct", ""));

        let min_time = ref(parseInt(loadCache("min_time", "5"), 10));
        let max_time = ref(parseInt(loadCache("max_time", "7"), 10));

        let msg_type = ref(loadCache("msg_type", "1"));

        let content = ref(loadCache("content", ""));
        let recv_list_txt = ref("");

        let preview_key = 0;
        let content_preview = computed(() => content.value.trim().split("\n").map(x => {return {text: x, key: ++preview_key};}));

        let pic_list = ref([]);
        let pic_url = ref("");

        let getPicList = () => {
            axios.get("/api/pic_list").then(
                rsp => {
                    if (rsp.data.code != 0) {
                        message.error(`[${rsp.data.code}]请求失败: ${rsp.data.msg}`);
                        return;
                    }
                    pic_list.value = rsp.data.data.list;
                }
            ).catch(
                err => message.error(JSON.stringify(err))
            );
        };

        onMounted(getPicList);

        watchEffect(() => {
            window.sessionStorage.setItem("sender_uid", sender_uid.value.toString());
            window.sessionStorage.setItem("sess", sess.value);
            window.sessionStorage.setItem("jct", jct.value);
            window.sessionStorage.setItem("min_time", min_time.value.toString());
            window.sessionStorage.setItem("max_time", max_time.value.toString());
            window.sessionStorage.setItem("content", content.value);
            window.sessionStorage.setItem("msg_type", msg_type.value);
        });

        let getRecvList = () => {
            if (recv_list_txt.value.trim().length == 0) {
                return {list: null, err: "接收者列表为空"};
            }
            let arr = recv_list_txt.value.trim().split("\n").map(x => x.trim());
            let err = "";
            let uniq = new Map();
            arr.forEach(item => {
                if (err.length > 0) return;
                if (!item.match(/^[0-9]+$/)) {
                    err = "接收者UID " + item + " 不是数字";
                }
                if (uniq.get(item)) {
                    err = "重复的接收者UID " + item;
                } else {
                    uniq.set(item, true);
                }
            });
            if (err.length > 0) {
                return {list: null, err: err};
            }
            return {list: arr.map(x => parseInt(x, 10)), err: ""}
        };

        let submit = () => {
            if (sender_uid.value < 1) {
                message.error("发送者UID必须大于0");
                return;
            }
            if (sess.value.trim().length == 0) {
                message.error("SESS_DATA未设置");
                return;
            }
            if (jct.value.trim().length == 0) {
                message.error("jct未设置");
                return;
            }
            if (min_time.value > max_time.value) {
                message.error("发送时间间隔设置错误: 最低间隔时间大于最高间隔时间");
                return;
            }
            if (min_time.value < 0 || max_time.value < 0) {
                message.error("发送时间间隔设置错误: 发送间隔必须大于0");
                return;
            }
            if (msg_type.value == "1" && content.value.trim().length == 0) {
                message.error("消息内容未设置");
                return;
            }
            if (msg_type.value == "2" && pic_url.value.trim().length == 0) {
                message.error("发送图片未选择");
                return;
            }
            let tmp = getRecvList();
            if (tmp.err.length > 0) {
                message.error("接收者设置错误: " + tmp.err);
                return;
            }
            submitting.value = false;
            recievers.value = tmp.list;
            showModal.value = true;
        };

        let currRoute = inject('currRoute');
        let doSubmit = () => {
            let req = {
                sender_uid: sender_uid.value,
                sess: sess.value,
                jct: jct.value,
                min_time: min_time.value,
                max_time: max_time.value,
                msg_type: msg_type.value,
                content: msg_type.value == "1" ? content.value : pic_url.value,
                recv_list: getRecvList().list,
            };
            submitting.value = true;
            axios.post("/api/batch_send", req).then(rsp => {
                submitting.value = false;
                if (rsp.data.code == 0) {
                    window.history.pushState(null, "", "/progress");
                    currRoute.value = "/progress";
                    return;
                }
                if (rsp.data.code == 233) {
                    message.error("当前有未完成的群发任务，请到“查看进度”查看任务信息");
                    return;
                }
                message.error(`[${rsp.data.code}]请求失败: ${rsp.data.msg}`);
            }).catch(err => {
                submitting.value = false;
                message.error(JSON.stringify(err));
            });
        };

        return {
            showModal,
            submitting,
            recievers,
            msg_type,
            content_preview,

            sender_uid,
            sess,
            jct,
            min_time,
            max_time,
            content,
            recv_list_txt,
            pic_list,
            pic_url,

            submit,
            doSubmit,
            getPicList,
        };
    },
};
</script>

<style>
.msg-preview {
    background-color: rgba(128, 128, 128, 0.1);
    padding: 16px;
}
.n-image--preview-disabled img {
    max-width: 300px;
    width: auto;
}
</style>