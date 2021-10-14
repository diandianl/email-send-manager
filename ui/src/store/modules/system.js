import storage from '@/utils/storage'
const state = {
  info: storage.get('app_info')
}

const mutations = {
  SET_INFO: (state, data) => {
    state.info = data
    storage.set('app_info', data)
  }
}

const actions = {
  settingDetail({ commit }) {
    commit('SET_INFO', { sys_app_name: 'Email Send Manager' })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
