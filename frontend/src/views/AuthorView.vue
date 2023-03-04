<template>
  <div>
    <Header></Header>
    <main class="article-list">
      <article class="text-left" v-for="author in authors" :key="author.Name">
        <a href="{{author.Uri}}" target="_blank" class="article-title">
          <p>{{ author.Name }} {{ author.Email }}</p>
          <hr />
          <p>{{ author.Description }}</p></a
        >
      </article>
    </main>
  </div>
</template>

<style scoped>
article {
  padding: 0.8rem;
  margin: 0.5rem;
  word-break: break-all;
  box-shadow: 1px 1px 5px 0 rgba(0, 0, 0, 0.02),
    1px 1px 15px 0 rgba(0, 0, 0, 0.03);
  width: 24rem;
  max-width: 800px;
}
.article-list {
  display: flex;
  flex-flow: column;
  align-items: center;
  min-height: calc(100vh - 465px);
}
.article-title {
  color: gray;
}
.article-content {
  color: gray;
}
.article-time {
  color: gray;
}
</style>

<script>
import HeaderComponent from "@/components/Header.vue";

export default {
  name: "AuthorView",
  components: {
    Header: HeaderComponent,
  },
  computed: {
    authors: function () {
      const data = this.$store.state.data;
      return data.author;
    },
  },
};
</script>