<template>
  <div class="header">
    <div class="avatar all">
      <el-row>
        <el-avatar @click="author = ''"></el-avatar>
      </el-row>
      <el-row>
        <el-text>All</el-text>
      </el-row>
    </div>
    <el-scrollbar class="scroll">
      <div class="authors">
        <div class="avatar" v-for="item in authors">
          <el-row>
            <el-avatar @click="author = item.name"></el-avatar>
          </el-row>
          <el-row>
            <el-text>{{ item.name }}</el-text>
          </el-row>
        </div>
      </div>
    </el-scrollbar>
  </div>
  <el-divider>Articles</el-divider>
  <el-timeline class="timeline">
    <el-timeline-item
      v-for="(item, i) in filteredArticles"
      :timestamp="item?.time"
      placement="top"
    >
      <el-card>
        <template #header>
          <el-row
            style="
              display: flex;
              flex-flow: row;
              justify-content: space-between;
            "
          >
            <span>
              <el-text class="article-title">{{ item?.title }}</el-text
              ><br />
              <el-text type="primary" style="margin-right: 0.5rem">|</el-text>
              <el-text>{{ item?.name }}</el-text>
            </span>
            <el-button class="button" text @click="curr = curr == i ? -1 : i"
              >Read</el-button
            >
          </el-row>
        </template>
        <div style="padding: 1rem" :class="{ 'article-fold': i != curr }">
          <el-text
            v-html="item?.content"
            style="white-space: pre-wrap; word-break: break-all"
          ></el-text>
        </div>
      </el-card>
    </el-timeline-item>
  </el-timeline>
</template>
<script>
import http from "../utils/http";

export default {
  data() {
    return {
      author: "",
      curr: -1,
      authors: [
        {
          name: "张三",
          article: [{ title: "test title", content: "test content" }],
        },
        {
          name: "李四",
          article: [{ title: "test title", content: "test content" }],
        },
      ],
      isShow: false,
    };
  },
  created() {
    this.fetchFeed();
  },
  computed: {
    allArticles() {
      return this.authors.reduce((prev, cur) => {
        return prev.concat(
          cur.article?.map((item) => {
            return {
              ...item,
              name: cur.name,
            };
          })
        );
      }, []);
    },
    filteredArticles() {
      if (this.author) {
        return this.authors
          .find((item) => item.name === this.author)
          .article.map((item) => {
            return {
              ...item,
              name: this.author,
            };
          });
      } else {
        return this.allArticles;
      }
    },
  },
  methods: {
    fetchFeed() {
      http.get("/feed").then((res) => {
        if (res.status === 200) {
          this.authors = res.data.author;
        }
      });
    },
  },
};
</script>

<style scoped>
.header {
  color: #333;
  height: 5rem;
  display: flex;
  flex-flow: row nowrap;
  align-items: center;
}
.authors {
  padding-inline: 1rem;
  display: flex;
  flex-flow: row nowrap;
  align-items: center;
}
.article-fold {
  display: none;
}
.timeline {
  padding-inline-start: 0px;
}
.avatar {
  margin: 10px 10px;
  display: flex;
  flex-flow: column nowrap;
  align-items: center;
}
.el-avatar:hover {
  cursor: pointer;
  box-shadow: 0 0 10px #ccc;
}
.all {
  border-right: #333 1px dotted;
  padding-right: 1rem;
}
.article-title {
  font-size: 1rem;
  font-weight: 400;
}
p {
  text-indent: 2rem;
}
</style>

<style lang="less">
.el-card__body {
  padding: 0px;
}
figure {
  margin: 0.5rem;
  .code {
    word-break: break-all;
  }
}
img {
  width: 100%;
}
</style>
