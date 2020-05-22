import Vue from 'vue'
import VueRouter from 'vue-router'
import Landing from '@/views/Landing.vue'
import Client from '@/views/Client.vue'
import Email from '@/components/Email.vue'
import EmailList from '@/components/EmailList.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Landing',
    component: Landing
  },
  {
    path: '/i',
    component: Client,
    children: [
      {
        path: '',
        name: 'EmailList',
        component: EmailList
      },
      {
        path: ':id',
        name: 'Email',
        component: Email
      }
    ]
  }
]

const router = new VueRouter({
  routes
})

export default router
