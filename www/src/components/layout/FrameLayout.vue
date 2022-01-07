<template>
    <n-message-provider>
    <n-layout>
        <n-layout-header>{{header}}</n-layout-header>
        <n-layout>
            <n-layout-header class="menu-holder" bordered>
                <n-menu v-model:value="activeKey" mode="horizontal" :options="menu" />
            </n-layout-header>
            <n-layout-content embedded class="content-frame">
                <n-card class="content-holder"><slot></slot></n-card>
            </n-layout-content>
        </n-layout>
        <n-layout-footer bordered>{{footer}}</n-layout-footer>
    </n-layout>
    </n-message-provider>
</template>

<script>
import { toRefs, watch, ref } from "vue";

export default {
    components: {},
    props: ['header', 'footer', 'menukey', 'menu'],

    setup(props) {
        let { menukey } = toRefs(props);
        let activeKey = ref("");
        activeKey.value = menukey.value;
        watch(menukey, () => { activeKey.value = menukey.value });
        return {activeKey};
    },
}
</script>

<style>
.n-layout-header,
.n-layout-footer {
    padding: 24px;
    text-align: center;
    background: white;
}

.menu-holder {
    padding: 0 24px;
}

.n-layout-header {
    font-size: 16px;
}

.n-layout-sider {
    min-height: 80vh;
}

.n-layout-content {
    background: white;
}

.content-frame {
    padding: 24px;
    background-color: rgba(128, 128, 128, 0.1);
}

.content-holder {
    width: 60%;
    margin: 0 auto;
}
</style>