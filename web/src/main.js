import Vue from 'vue'
import App from './App.vue'
import store from './store'
import router from './router'
import Vuex from 'vuex'

Vue.config.productionTip = false

Vue.use(Vuex);

new Vue({
  store,
  router,
  render: h => h(App)
}).$mount('#app')
