import Vue from 'vue'

import 'normalize.css/normalize.css' // a modern alternative to CSS resets

import Element from 'element-ui'
import './styles/element-variables.scss'

import '@/styles/index.scss' // global css
import '@/styles/admin.scss'

import App from './App'
import store from './store'
import router from './router'

import { parseTime, resetForm } from '@/utils/costum'

import './icons' // icon

import Pagination from '@/components/Pagination'
import BasicLayout from '@/layout/BasicLayout'

// 全局方法挂载
Vue.prototype.parseTime = parseTime
Vue.prototype.resetForm = resetForm

// 全局组件挂载
Vue.component('Pagination', Pagination)
Vue.component('BasicLayout', BasicLayout)

Vue.prototype.msgSuccess = function(msg) {
  this.$message({ showClose: true, message: msg, type: 'success' })
}

Vue.prototype.msgError = function(msg) {
  this.$message({ showClose: true, message: msg, type: 'error' })
}

Vue.prototype.msgInfo = function(msg) {
  this.$message.info(msg)
}

Vue.use(Element, {
  size: 'medium' // set element-ui default size
})

Vue.config.productionTip = false

store.dispatch('permission/generateRoutes')

new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App)
})
