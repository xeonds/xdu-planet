<template>
  <div class="header">
    <div class="avatar all">
      <el-row>
        <el-avatar @click="author = ''"> All</el-avatar>
      </el-row>
      <el-row>
        <el-text>All</el-text>
      </el-row>
    </div>
    <el-scrollbar class="scroll">
      <div class="authors">
        <div class="avatar" v-for="item in authors">
          <el-row>
            <el-avatar @click="author = item.name">{{
              item.name.substring(0, 1)
            }}</el-avatar>
          </el-row>
          <el-row>
            <el-text style="white-space: nowrap">{{ item.name }}</el-text>
          </el-row>
        </div>
      </div>
    </el-scrollbar>
  </div>
  <el-divider>Articles</el-divider>
  <el-timeline class="timeline">
    <el-timeline-item v-for="(item, index) in groupedArticles" :timestamp="item.date" :type="item.list[0].type"
      :hollow="item.list[0].hollow" placement="top">
      <el-card style="box-shadow: none; margin-bottom: 1rem" v-for="(article, subIndex) in item.list">
        <template #header>
          <el-row style="
              display: flex;
              flex-flow: row;
              justify-content: space-between;
            ">
            <span>
              <el-text class="article-title">{{ article.title }}</el-text><br />
              <el-text type="primary" style="margin-right: 0.5rem">|</el-text>
              <el-text>{{ article.name }}</el-text>
            </span>
            <div style="display: flex; flex-flow: row wrap; justify-content: right">
              <el-button text @click="
                curr =
                curr == `${index},${subIndex}` ? '' : `${index},${subIndex}`
                ">{{
    curr == `${index},${subIndex}` ? "Hide" : "Read"
  }}</el-button>
              <el-button text @click="viewUrl(article.url)">Source Link</el-button>
            </div>
          </el-row>
        </template>
        <div style="
            border-top: 1px solid var(--el-card-border-color);
          " :class="{ 'article-fold': `${index},${subIndex}` != curr }">
          <mavon-editor v-model="article.content" :subfield="false" :defaultOpen="'preview'" :toolbarsFlag="false"
            :boxShadow="false" :transition="false" />
        </div>
      </el-card>
    </el-timeline-item>
  </el-timeline>
  <el-pagination layout="prev, pager, next" :total="filteredArticles.length" :page-size="pageSize"
    v-model:current-page="currPage" />
</template>
<script>
import { mavonEditor } from "mavon-editor";
import http from "../utils/http";
import day from "../utils/day";
import axios from "axios";

export default {
  components: { mavonEditor },
  data() {
    return {
      author: "",
      curr: "",
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
      pageSize: 16,
      currPage: 1,
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
            cur.article.map((item) => {
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
            content: this.getBody(item.content),
          };
        })
        .sort((a, b) => {
          return new Date(b.time) - new Date(a.time);
        });
    },
    filteredArticles() {
      return this.allArticles
        .filter((item) => {
          return this.author ? item.name === this.author : true;
        })
        .reduce((prev, cur) => {
          if (
            day(cur.time).format("M") !=
            day(prev[prev.length - 1]?.time).format("M")
          ) {
            return prev.concat({
              ...cur,
              type: "primary",
            });
          }
          return prev.concat(cur);
        }, [])
        .reduce((prev, cur) => {
          if (
            day(cur.time).format("YYYY") !=
            day(prev[prev.length - 1]?.time).format("YYYY")
          ) {
            return prev.concat({
              ...cur,
              hollow: true,
            });
          }
          return prev.concat(cur);
        }, []);
    },
    groupedArticles() {
      return this.filteredArticles
        .slice(
          (this.currPage - 1) * this.pageSize,
          this.currPage * this.pageSize
        )
        .reduce((prev, cur) => {
          const date = cur.time.split(" ")[0];
          const index = prev.findIndex((item) => item.date === date);
          if (index === -1) {
            prev.push({
              date,
              list: [cur],
            });
          } else {
            prev[index].list.push(cur);
          }
          return prev;
        }, []);
    },
  },
  watch: {
    author() {
      this.currPage = 1;
    },
  },
  methods: {
    fetchFeed() {
      axios
        .get("db.json")
        .then((res) => {
          this.authors = res.data.author;
        })
        .catch(() => {
          http.get("/feed").then((res) => {
            this.authors = res.data.author;
          });
        });
    },
    viewUrl(url) {
      window.open(url);
    },
    getBody(content) {
      var REG_BODY = /<body[^>]*>([\s\S]*)<\/body>/;
      var result = REG_BODY.exec(content);
      if (result && result.length === 2)
        return result[1];
      return content;
    }
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

.el-avatar {
  background-color: var(--el-color-primary);
}

.el-avatar:hover {
  cursor: pointer;
  box-shadow: 0 0 10px var(--el-color-primary);
}

.all {
  border-right: #333 1px dotted;
  padding-right: 1rem;
}

.article-title {
  font-size: 1rem;
  font-weight: 400;
}
</style>

<style lang="less">
.el-card__body {
  padding: 0px;
}

.el-card__header {
  border: none;
}
</style>
