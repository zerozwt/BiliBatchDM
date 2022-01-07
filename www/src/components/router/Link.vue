<template>
<a :href="href" @click.prevent="link"><slot></slot></a>
</template>

<script>
import { inject } from "vue";

export default {
    props: ['href'],

    setup(props) {
        let router = inject('router');
        let currRoute = inject('currRoute');
        let link = () => {
            window.history.pushState(null, router.title(props.href), props.href);
            currRoute.value = props.href;
        };
        return { link, };
    },
}
</script>

<style>
</style>