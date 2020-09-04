import Vue from "vue";
import VueRouter from "vue-router";
// import Start from "@/views/Start.vue";
import Welcome from "@/views/Welcome.vue";
import Blocks from "@/views/Blocks.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    component: Welcome,
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
