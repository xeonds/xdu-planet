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
      v-for="item in filteredArticles"
      :timestamp="item?.time"
      placement="top"
    >
      <el-card class="article">
        <template #header>
          <el-row>
            <el-col :span="12">
              <el-text>{{ item?.title }}</el-text>
              <el-text type="primary" style="margin-inline: 1rem">|</el-text>
              <el-text>{{ item?.name }}</el-text>
            </el-col>
            <el-col :span="12" style="text-align: right">
              <el-button class="button" text>Read</el-button>
            </el-col>
          </el-row>
        </template>
        <el-text v-html="item?.content" style="max-height: 4rem"></el-text>
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
      authors: [
        {
          name: "张三",
          article: [
            { title: "test title", content: "test content" },
            { title: "test title", content: "test content" },
          ],
        },
        {
          name: "李四",
          article: [{ title: "test title", content: "test content" }],
        },
        {
          name: "王五",
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
        return this.allArticles.filter((item) => item?.name === this.author);
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
          console.log(this.authors);
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
.article {
  margin-bottom: 1rem;
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
</style>

<style>
p {
  text-indent: 2rem;
}
</style>
