import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'layout',
      component: () => import('../layouts/DefaultLayout.vue'),
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('../views/Home.vue')
        },
        {
          path: 'users',
          name: 'users',
          component: () => import('../views/user/UserList.vue')
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('../views/settings/SystemSettings.vue')
        },
        {
          path: 'roles',
          name: 'roles',
          component: () => import('../views/role/RoleList.vue')
        },
        {
          path: 'profile',
          name: 'profile',
          component: () => import('../views/user/Profile.vue')
        },
        {
          path: 'logs',
          name: 'logs',
          component: () => import('../views/log/LogList.vue')
        },
        {
          path: 'files',
          name: 'files',
          component: () => import('../views/file/FileList.vue')
        }
      ]
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/Login.vue')
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.path !== '/login' && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/')
  } else {
    next()
  }
})

export default router 