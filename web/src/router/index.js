import {createRouter, createWebHashHistory} from 'vue-router'

const routes = [
    // {
    //     path: '/',
    //     redirect: '/login'
    // },
    // {
    //     path: '/init',
    //     name: 'Init',
    //     component: () => import('@/view/init/index.vue')
    // },
    {
        path: '/login',
        name: 'Login',
        component: () => import('@/view/login/index.vue')
    },
    {
        path: '/pay',
        name: 'Pay',
        component: () => import('@/view/vbox/payOrder/orderDetail.vue')
    },
    {
        path: '/pay2',
        name: 'Pay2',
        component: () => import('@/view/vbox/payOrder/orderDetailV2.vue')
    },
    {
        path: '/payTest',
        name: 'PayTest',
        component: () => import('@/view/vbox/payOrder/orderTest.vue')
    },
    {
        path: '/payDoc',
        name: 'PayDoc',
        component: () => import('@/view/doc/PayDoc.vue')
    },
// {
//   path: '/order/:ch',
//   name: 'OrderTask',
//   component: () => import('@/view/vbox/order/orderPayTask.vue')
// },
    {
        path: '/:catchAll(.*)',
        meta: {
            closeTab: true,
        },
        component: () => import('@/view/error/index.vue')
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

export default router
