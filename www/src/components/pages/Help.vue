<template>
<n-card title="发送者信息填写方式说明">
    <n-space vertical>
        <n-card title="B站UID">
            <p>在网页上登陆B站以后，进入B站个人空间页面（<a href="https://space.bilibili.com" target="_blank">space.bilibili.com</a>）。</p>
            <p>进入个人空间后，浏览器地址栏会变成类似这样:</p>
            <img src="/uid.png" />
            <p>其中"space.bilibili.com/"后面的数字部分，就是页面上的用户（也就是你）的B站UID。</p>
        </n-card>
        <n-card title="SESS_DATA和bili_jct">
            <template #header-extra>以chrome浏览器为例</template>
            <p>在B站的任意页面，按下键盘上的F12键，打开chrome的开发者工具:</p>
            <img src="/f12menu.png" />
            <p>如上图所示，点开“>>”菜单，选择“Application”，打开“Application”界面:</p>
            <img src="/cookie.png" />
            <p>在界面左边的“Cookies”中，随意点选一个B站域名，界面右边就会显示所有的Cookie信息。</p>
            <p>从右边的列表中，找到并分别点选“SESS_DATA”和“bili_jct”，下方的“Cookie Value”就会显示对应的值，将“SESS_DATA”和“bili_jct”对应的值复制到小助手的“发送私信”界面即可。</p>
            <p><b>重要: 复制的时候不要勾选“Cookie Value”旁边的“Show URL decoded”！！！</b></p>
        </n-card>
    </n-space>
</n-card>
</template>

<script>
</script>

<style>
img {
    max-width: 100%;
    width: auto;
}
</style>