import Vue from "vue";
import VueRouter from "vue-router";
import Start from "@/views/Start.vue";
import Blocks from "@/views/Blocks.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    component: Start,
  },
  {
    path: "/blocks",
    component: Blocks,
  },
];

const router = new VueRouter({
  mode: "hash",
  base: process.env.BASE_URL,
  routes,
});

export default router;
