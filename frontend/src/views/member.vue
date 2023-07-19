<template>
  <el-row style="margin-bottom: 1rem">
    <el-col :span="12">
      <el-statistic title="文章总数" :value="allArticles.length">
        <template #suffix>篇</template>
      </el-statistic>
    </el-col>
    <el-col :span="12">
      <el-statistic title="总字数" :value="wordCount">
        <template #suffix>字</template>
      </el-statistic>
    </el-col>
  </el-row>
  <el-row>
    <el-col :span="12">
      <el-statistic title="成员总数" :value="authors.length">
        <template #suffix>人</template>
      </el-statistic>
    </el-col>
    <el-col :span="12">
      <el-statistic title="历史" :value="timeCount">
        <template #suffix>天</template>
      </el-statistic>
    </el-col>
  </el-row>
  <el-divider>Members</el-divider>
  <el-card
    class="author-card"
    v-for="author in authors"
    style="box-shadow: none; margin-bottom: 1rem; padding: 0.8em"
  >
    <el-row style="margin-bottom: 16px">
      <el-text style="font-weight: 100">{{ author.name }}</el-text>
    </el-row>
    <el-row>
      <span>
        <el-text
          type="primary"
          style="margin-inline: 0.5rem; font-size: 1.25rem"
          >|</el-text
        >
        <el-text style="font-size: 1.25rem">{{ author.description }}</el-text>
      </span>
    </el-row>
    <el-row
      ><el-col :span="24" style="text-align: right">
        <el-button text @click="viewUrl(author.uri)">View Site</el-button>
      </el-col>
    </el-row>
  </el-card>
</template>

<script>
import http from "../utils/http";
import axios from "axios";
import day from "../utils/day";

export default {
  data() {
    return {
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
    };
  },
  created() {
    this.fetchFeed();
  },
  computed: {
    allArticles() {
      return this.authors
        .reduce((prev, cur) => {
          return prev.concat(
            cur.article?.map((item) => {
              return {
                ...item,
                name: cur.name,
              };
            })
          );
        }, [])
        .map((item) => {
          return {
            ...item,
            time: day(item.time).format("YYYY-MM-DD HH:mm:ss"),
          };
        })
        .sort((a, b) => {
          return new Date(b.time) - new Date(a.time);
        });
    },
    wordCount() {
      return this.allArticles.reduce((prev, cur) => {
        return prev + cur.content.length;
      }, 0);
    },
    timeCount() {
      const first = this.allArticles[this.allArticles.length - 1].time;
      const last = this.allArticles[0].time;
      return day(last).diff(day(first), "day");
    },
  },
  methods: {
    fetchFeed() {
      axios
        .get("db.json")
        .then((res) => {
          this.authors = res.data.author;
        })
        .catch((err) => {
          http.get("/feed").then((res) => {
            this.authors = res.data.author;
          });
        });
    },
    viewUrl(url) {
      window.open(url);
    },
  },
};
</script>
