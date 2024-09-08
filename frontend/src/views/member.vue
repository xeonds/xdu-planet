<template>
  <el-row style="margin-bottom: 1rem">
    <el-col :span="8">
      <el-statistic title="文章总数" :value="allArticles(authors).length">
        <template #suffix>篇</template>
      </el-statistic>
    </el-col>
    <el-col :span="8">
      <el-statistic title="成员总数" :value="authors.length">
        <template #suffix>人</template>
      </el-statistic>
    </el-col>
    <el-col :span="8">
      <el-statistic title="历史" :value="timeCount(allArticles(authors))">
        <template #suffix>天</template>
      </el-statistic>
    </el-col>
  </el-row>
  <el-row>
  </el-row>
  <el-divider>Members</el-divider>
  <el-card class="author-card" v-for="author in authors" style="box-shadow: none; margin-bottom: 1rem; padding: 0.8em">
    <el-row style="margin-bottom: 16px">
      <el-text style="font-weight: 100">{{ author.name }}</el-text>
    </el-row>
    <el-row>
      <span>
        <el-text type="primary" style="margin-inline: 0.5rem; font-size: 1.25rem">|</el-text>
        <el-text style="font-size: 1.25rem">{{ author.description }}</el-text>
      </span>
    </el-row>
    <el-row><el-col :span="24" style="text-align: right">
        <el-button text @click="viewUrl(author.uri)">View Site</el-button>
      </el-col>
    </el-row>
  </el-card>
</template>

<script lang="ts" setup>
import { Article, Author, Feed } from "../api/home";
import { onMounted, ref } from "vue";
import dayjs from "dayjs";
import { http, useHttp } from "../utils/http";

const authors = ref(new (Array<Author>));

const allArticles = (authors: Author[]) =>
  authors
    .reduce(
      (prev: Article[], cur) =>
        prev.concat(cur.article.map((item) => ({ ...item, name: cur.name, })))
      , []
    )
    .map((item) => ({
      ...item,
      time: dayjs(item.time).format("YYYY-MM-DD HH:mm:ss"),
    }))
    .sort((a, b) => (dayjs(b.time).diff(a.time)));

const timeCount = (items: Article[]) =>
  (items.length === 0)
    ? 0
    : dayjs(items[0].time).diff(dayjs(items[items.length - 1].time), "day");

const viewUrl = (url: string) => window.open(url);

onMounted(async () => {
  const { data, err } = await (async () => {
    const response = await http.get<Feed>("/feed");
    if (response.err.value != null) return await useHttp(".")().get<Feed>("/index.json");
    return response;
  })()
  if (err.value != null || data.value == null) {
    console.error(err.value);
    return;
  }
  authors.value = data.value.author.map((item, index) => ({ ...item, name: item.name == "" ? `unDefined Author ${index}` : item.name }));
});

</script>
