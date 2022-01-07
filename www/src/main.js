import { createApp } from 'vue';
import naive from 'naive-ui';

import './App.css';
import App from './App';

if (process.env.NODE_ENV == "development") {
    console.log("DEV_ENV");
    require('./mocks');
}

createApp(App).use(naive).mount('#app')