import RouteTable from '@/routes';

export default {
    title(href) {
        return this.rule(href).title
    },

    page(href) {
        let item = this.rule(href);
        return {
            component: require(`../pages/${item.component}.vue`).default,
            menukey: item.key,
        }
    },

    rule(href) {
        let path = (new URL(window.location.protocol + "//" + window.location.hostname + href)).pathname
        for (let i = 0; i < RouteTable.rules.length; ++i) {
            if (path.match(RouteTable.rules[i].rule)) {
                return RouteTable.rules[i]
            }
        }
        return RouteTable.default
    },
}