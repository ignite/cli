import Vue from "vue";
import VueRouter from "vue-router";
import Start from "@/views/Start.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    component: Start,
  },
];

const router = new VueRouter({
  mode: "hash",
  base: process.env.BASE_URL,
  routes,
});

export default router;
