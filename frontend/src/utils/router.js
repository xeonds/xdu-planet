import { createRouter, createWebHashHistory } from "vue-router";

const routes = [
  {
    path: "/",
    component: () => import("../views/home.vue"),
  },
  {
    path: "/about",
    component: () => import("../views/about.vue"),
  },
];

export default createRouter({
  history: createWebHashHistory(),
  routes: routes,
  beforeEach: (to, _from, next) => {
    if (to.meta.title) {
      document.title = to.meta.title ? to.meta.title : "加载中";
    }
    next();
  },
});
