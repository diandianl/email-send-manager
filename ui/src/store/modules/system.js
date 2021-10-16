const state = {
  info: { sys_app_name: 'ESM' }
}

const mutations = {
  SET_INFO: (state, data) => {
    state.info = data
  }
}

const actions = {
  settingDetail({ commit }, info) {
    commit('SET_INFO', info)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
