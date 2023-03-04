import { createStore } from "vuex";

export default createStore({
  state() {
    return { data: [] };
  },
  mutations: {
    setData(state, payload) {
      state.data = payload;
    },
  },
});
