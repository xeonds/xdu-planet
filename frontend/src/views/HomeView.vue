<template>
  <div>
    <Header></Header>
    <main class="article-list">
      <div
        class="recursive-container"
        v-for="(articles, month) in sortedArticles"
        :key="month"
      >
        <div class="mon-seperator">
          <div class="timeline"></div>
          <div>{{ month }}</div>
        </div>
        <article
          class="text-left"
          v-for="article in articles"
          :key="article.title"
        >
          <a :href="article.Url" target="_blank" class="article-title">
            <h3>{{ article.Title }}</h3>
          </a>
          <p class="article-content">{{ article.Content }}</p>
          <p class="article-time">{{ article.Time }}</p>
        </article>
      </div>
    </main>
  </div>
</template>

<script>
import HeaderComponent from "../components/Header.vue";

export default {
  name: "HomeView",
  components: {
    Header: HeaderComponent,
  },
  computed: {
    sortedArticles() {
      const sorted = {};
      const data = this.$store.state.data;
      data.article.forEach((article) => {
        const month = new Date(article.Time).toLocaleString("zh-cn", {
          month: "long",
          year: "numeric",
        });
        if (!sorted[month]) {
          sorted[month] = [];
        }
        sorted[month].push(article);
      });
      return sorted;
    },
  },
};
</script>

<style scoped>
.recursive-container {
  width: calc(100% - 1.5rem);
  max-width: 800px;
}
header {
  padding: 2rem;
  margin-top: 6rem;
  margin-bottom: 3rem;
}

article {
  padding: 0.8rem;
  margin: 0.5rem;
  word-break: break-all;
  box-shadow: 1px 1px 5px 0 rgba(0, 0, 0, 0.02),
    1px 1px 15px 0 rgba(0, 0, 0, 0.03);
  width: calc(100% - 1.5rem);
  max-width: 800px;
}

footer {
  margin-top: 6rem;
  padding: 2rem;
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

.mon-seperator {
  padding-top: 4rem;
  padding-bottom: 1.5rem;
  width: calc(100% - 1rem);
  max-width: 800px;
  text-align: left;
  font-size: 1.8rem;
  color: black;
}

.timeline {
  width: 0.75rem;
  height: 0.75rem;
  clip-path: polygon(0 0, 100% 0, 100% 100%, 0 100%);
  background-color: red;
  display: block;
  position: relative;
  top: calc(1rem + 1.5 * 0.85rem / 2 - 0.75rem / 2);
  left: -1rem;
}
</style>

