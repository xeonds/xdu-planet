import { createRouter, createWebHashHistory } from "vue-router";

const routes = [
  {
    path: "/",
    component: () => import("../views/home.vue"),
  },
  {
    path: "/member",
    component: () => import("../views/member.vue"),
  },
  {
    path: "/about",
    component: () => import("../views/about.vue"),
  },
];

export default createRouter({
  history: createWebHashHistory(),
  routes: routes,
});
