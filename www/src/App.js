import { computed, h, provide, reactive, ref } from 'vue';
import FrameLayout from './components/layout/FrameLayout.vue';
import router from './components/router/router';
import AppMenu from './menu';

export default {
    name: 'App',
    components: {},

    setup() {
        let currRoute = ref(window.location.pathname);
        window.addEventListener('popstate', () => {
            currRoute.value = window.location.pathname
        })

        let FrameParams = reactive({
            header: "B站私信群发小助手",
            footer: "Bilibili Batch Direct Message Assistant",
            menukey: "",
            menu: AppMenu,
        });

        let page = computed(() => {
            let routePage = router.page(currRoute.value);
            FrameParams.menukey = routePage.menukey;
            return h(FrameLayout, FrameParams, () => h(routePage.component));
        });

        provide('router', router);
        provide('currRoute', currRoute);

        return {page};
    },

    render() {
        return h(this.page);
    }
}
