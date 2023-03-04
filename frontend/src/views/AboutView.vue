<template>
  <div>
    <Header></Header>
    <main class="about">
      <h3>About</h3>
      <p>
        一个简单的RSS聚合站。<br />
        共同写作，共同阅读，共同分享。这就是这个站点存在的意义。
      </p>
      <hr />
      <p>目前没怎么做UI，只是先把主要功能完成了而已。</p>
    </main>
  </div>
</template>

<style scoped>
.about {
  padding: 2rem;
  max-width: 800px;
  width: 90vw;
  min-height: calc(100vh - 465px);
}
</style>

<script>
import HeaderComponent from "../components/Header.vue";

export default {
  name: "AboutView",
  components: {
    Header: HeaderComponent,
  },
};
</script>