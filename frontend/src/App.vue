<template>
  <div class="container-fluid">
    <nav class="navbar navbar-expand-sm bg-light navbar-light fixed-top">
      <router-link class="navbar-brand" to="/">XDU-Planet</router-link>
      <button
        class="navbar-toggler"
        type="button"
        data-toggle="collapse"
        data-target="#collapsibleNavbar"
      >
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="collapsibleNavbar">
        <ul class="navbar-nav">
          <li class="nav-item">
            <router-link class="nav-link" to="/">主页</router-link>
          </li>
          <li class="nav-item">
            <router-link class="nav-link" to="/author">成员</router-link>
          </li>
          <!-- <li class="nav-item">
        <router-link class="nav-link" to="/analyze">统计</router-link>
      </li> -->
          <li class="nav-item">
            <router-link class="nav-link" to="/about">关于</router-link>
          </li>
        </ul>
      </div>
    </nav>
    <router-view />
    <Footer></Footer>
  </div>
</template>

<script>
import FooterComponent from "./components/Footer.vue";

export default {
  name: "App",
  created: function () {
    this.fetchData();
  },
  components: {
    Footer: FooterComponent,
  },
  methods: {
    fetchData: async function () {
      const response = await this.axios.get("db.json");
      this.$store.commit("setData", response.data);
    },
  },
};
</script>

<style>
.container-fluid {
  display: flex;
  flex-flow: column;
  align-items: center;
  padding: 0;
}
</style>