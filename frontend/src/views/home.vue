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
    <el-timeline-item v-for="(item, index) in groupedArticles(filteredArticles(authors, author), currPage, pageSize)"
      :timestamp="item.date" :type="item.list[0].type" :hollow="item.list[0].hollow" placement="top">
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
              <el-button text @click="getBody(article.title, article.content, `${index},${subIndex}`)">Read</el-button>
              <el-button text @click="viewUrl(article.url)">Source Link</el-button>
            </div>
          </el-row>
        </template>
      </el-card>
    </el-timeline-item>
  </el-timeline>
  <el-pagination layout="prev, pager, next" :total="filteredArticles(authors, author).length" :page-size="pageSize"
    v-model:current-page="currPage" />
  <el-drawer v-model="viewArticleVisible" title="I am the title" :direction="'btt'" :size="'100%'">
    <template #header>
      <h3>{{ title }}</h3>
    </template>
    <editor.mavonEditor v-model="content" :subfield="false" :defaultOpen="'preview'" :toolbarsFlag="false" :boxShadow="false"
      :transition="false" />
  </el-drawer>
</template>

<script lang="ts" setup>
import editor from "mavon-editor";
import { onMounted, ref } from "vue";
import { Author, Article, Feed } from "../api/home";
import dayjs from "dayjs";
import { http, useHttp } from "../utils/http";

type ArticleEx = Article & {
  name?: string;
  type?: string;
  hollow?: boolean;
};

type GroupedArticle = {
  date: string;
  list: ArticleEx[];
};

const author = ref("");
const curr = ref("");
const authors = ref(new Array<Author>());
const pageSize = ref(16);
const currPage = ref(1);

const allArticles = (authors: Author[]) =>
  authors
    .reduce(
      (prev: ArticleEx[], cur) =>
        prev.concat(cur.article.map((item) => ({ ...item, name: cur.name, })))
      , []
    )
    .map((item) => ({
      ...item,
      time: dayjs(item.time).format("YYYY-MM-DD HH:mm:ss"),
      content: item.content,
    }))
    .sort((a, b) => (dayjs(b.time).diff(a.time)));

const filteredArticles = (authors: Author[], author?: string) =>
  allArticles(authors)
    .filter((item) => author ? item.name === author : true)
    .reduce((prev: ArticleEx[], cur) =>
      (dayjs(cur.time).format("M") != dayjs(prev[prev.length - 1]?.time).format("M"))
        ? prev.concat({ ...cur, type: "primary", })
        : prev.concat(cur)
      , [])
    .reduce((prev: ArticleEx[], cur) =>
      (dayjs(cur.time).format("YYYY") != dayjs(prev[prev.length - 1]?.time).format("YYYY"))
        ? prev.concat({ ...cur, hollow: true, })
        : prev.concat(cur)
      , []);

const groupedArticles = (items: ArticleEx[], currPage: number, pageSize: number) =>
  items
    .slice((currPage - 1) * pageSize, currPage * pageSize)
    .reduce((prev: GroupedArticle[], cur) => {
      const date = cur.time.split(" ")[0];
      const index = prev.findIndex((item) => item.date === date);
      if (index === -1) {
        prev.push({ date, list: [cur] });
      } else {
        prev[index].list.push(cur);
      }
      return prev;
    }, []);


const viewUrl = (url: string) => window.open(url);

const title = ref("");
const content = ref("");
const viewArticleVisible = ref(false);
const getBody = async (t: string, url: string, index: string) => {
  const { data, err } = await useHttp("")().get<any>(encodeURI(url), false);
  if (err.value != null || data.value == null) {
    console.error(err.value);
    return `加载 ${url} 失败：${err.value}`;
  }
  curr.value = curr.value == index ? '' : index;
  var REG_BODY = /<body[^>]*>([\s\S]*)<\/body>/;
  var result = REG_BODY.exec(data.value);
  title.value = t;
  content.value = (result && result.length === 2) ? result[1] : url;
  viewArticleVisible.value = true;
};

onMounted(async () => {
  const { data, err } = await (async () => {
    const response = await http.get<Feed>("/feed");
    if (response.err.value != null) return await useHttp("")().get<Feed>("/index.json");
    else return response;
  })()
  if (err.value != null || data.value == null) {
    console.error(err.value);
    return;
  }
  authors.value = data.value.author.map((item, index) => ({ ...item, name: item.name == "" ? `unDefined Author ${index}` : item.name }));
});
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
